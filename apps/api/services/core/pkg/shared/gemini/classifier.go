package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

const (
	SeverityLow    = "LOW"
	SeverityMedium = "MEDIUM"
	SeverityHigh   = "HIGH"
)

// Thresholds as a share of the policy ceiling. A loss is judged relative to the
// cover behind it, so the same figure is severe on a small policy and routine on
// a large one.
const (
	highSeverityShare   = 0.50
	mediumSeverityShare = 0.15
)

// ClassificationInput is what Loop 1 knows before a claim is triaged.
type ClassificationInput struct {
	ClaimNumber         string
	ClaimType           string
	IncidentDescription string
	EstimatedLoss       float64
	MaxCoverageAmount   float64
	DocumentTypes       []string
}

// ClassificationResult is Loop 1's triage of one claim.
type ClassificationResult struct {
	Severity         string
	SurveyRequired   bool
	MissingDocuments []string
	Reasons          string
	Source           string
}

type Classifier interface {
	Classify(ctx context.Context, in ClassificationInput) (ClassificationResult, error)
}

// NewClassifier mirrors New: same credential resolution, same deterministic
// fallback, so an unconfigured environment still triages real claim data.
func NewClassifier(ctx context.Context) (Classifier, string) {
	client, model, backend, ok := newModelClient(ctx)
	if !ok {
		return deterministicClassifier{}, SourceDeterministic
	}
	return &geminiClassifier{client: client, model: model, backend: backend}, model + " via " + backend
}

type geminiClassifier struct {
	client  *genai.Client
	model   string
	backend string
}

const classificationInstruction = `You are the Intake Agent in an insurance claims operations platform.

You triage one motor claim before it is assigned to anyone. Two judgements are yours:

SEVERITY - LOW, MEDIUM or HIGH. Weigh the estimated loss against the policy ceiling and
what the incident description actually says. A minor scrape is LOW however expensive the
policy. Structural damage, injury, multiple vehicles, suspected total loss, or a loss
consuming a large share of the cover is HIGH.

SURVEY_REQUIRED - whether a surveyor must physically inspect before assessment. Require
one when the described damage cannot be verified from photographs alone, when the loss is
large relative to cover, or when the account is internally inconsistent. Do not require one
for small, clearly described damage with photographs on file; an unnecessary inspection
costs days.

Justify both from the figures and the wording you were given. Never invent a document, an
amount, or a detail the description does not contain. If the description is empty, say so
and lean on the figures alone.`

func (c *geminiClassifier) Classify(ctx context.Context, in ClassificationInput) (ClassificationResult, error) {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	resp, err := c.client.Models.GenerateContent(ctx, c.model, genai.Text(buildClassificationPrompt(in)), &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(classificationInstruction, genai.RoleUser),
		Temperature:       genai.Ptr[float32](assessmentTemperature),
		ResponseMIMEType:  "application/json",
		ResponseSchema:    classificationSchema(),
	})
	if err != nil {
		return ClassificationResult{}, fmt.Errorf("gemini classification call: %w", err)
	}

	var parsed struct {
		Severity       string `json:"severity"`
		SurveyRequired bool   `json:"survey_required"`
		Reasons        string `json:"reasons"`
	}
	if err := json.Unmarshal([]byte(resp.Text()), &parsed); err != nil {
		return ClassificationResult{}, fmt.Errorf("gemini classification response was not valid JSON: %w", err)
	}

	severity, err := normalizeSeverity(parsed.Severity)
	if err != nil {
		return ClassificationResult{}, err
	}
	if strings.TrimSpace(parsed.Reasons) == "" {
		return ClassificationResult{}, fmt.Errorf("gemini classification returned no reasoning")
	}

	return ClassificationResult{
		Severity:       severity,
		SurveyRequired: parsed.SurveyRequired,
		// Which required documents are absent is set arithmetic, not judgement.
		// Computing it here keeps the model from naming a document that is on file.
		MissingDocuments: missingDocuments(in.DocumentTypes),
		Reasons:          strings.TrimSpace(parsed.Reasons),
		Source:           c.model + " via " + c.backend,
	}, nil
}

func classificationSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"severity": {
				Type:        genai.TypeString,
				Description: "One of LOW, MEDIUM, or HIGH.",
			},
			"survey_required": {
				Type:        genai.TypeBoolean,
				Description: "True when a surveyor must inspect before the claim can be assessed.",
			},
			"reasons": {
				Type:        genai.TypeString,
				Description: "Auditable justification for the severity and the survey decision, citing the figures and wording supplied.",
			},
		},
		Required: []string{"severity", "survey_required", "reasons"},
	}
}

func normalizeSeverity(raw string) (string, error) {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case SeverityLow:
		return SeverityLow, nil
	case SeverityMedium:
		return SeverityMedium, nil
	case SeverityHigh:
		return SeverityHigh, nil
	default:
		return "", fmt.Errorf("gemini classification returned an unknown severity %q", raw)
	}
}

func buildClassificationPrompt(in ClassificationInput) string {
	var b strings.Builder

	b.WriteString("Triage this claim.\n\nCLAIM\n")
	fmt.Fprintf(&b, "  Number:     %s\n", in.ClaimNumber)
	fmt.Fprintf(&b, "  Type:       %s\n", in.ClaimType)
	fmt.Fprintf(&b, "  Est. loss:  %.2f\n", in.EstimatedLoss)
	if in.MaxCoverageAmount > 0 {
		fmt.Fprintf(&b, "  Policy ceiling: %.2f (this loss is %.0f%% of cover)\n",
			in.MaxCoverageAmount, 100*in.EstimatedLoss/in.MaxCoverageAmount)
	} else {
		b.WriteString("  Policy ceiling: not established for this claim\n")
	}
	fmt.Fprintf(&b, "  Incident:   %s\n", valueOrUnknown(in.IncidentDescription))

	b.WriteString("\nDOCUMENTS ON FILE\n")
	if len(in.DocumentTypes) == 0 {
		b.WriteString("  none\n")
	} else {
		for _, d := range in.DocumentTypes {
			fmt.Fprintf(&b, "  - %s\n", d)
		}
	}

	return b.String()
}

// Applies the same triage without a model, so an unconfigured environment still
// classifies from the real claim rather than writing constants.
type deterministicClassifier struct{}

func (deterministicClassifier) Classify(_ context.Context, in ClassificationInput) (ClassificationResult, error) {
	severity := severityFromLoss(in)
	missing := missingDocuments(in.DocumentTypes)

	// Without a model to read the incident description, a large loss relative to
	// cover is the only signal available for sending an inspector.
	surveyRequired := severity == SeverityHigh

	var reasons string
	switch {
	case in.MaxCoverageAmount <= 0:
		reasons = fmt.Sprintf("Claim %s has no policy ceiling established, so severity falls back to the estimated loss of %.2f alone.",
			in.ClaimNumber, in.EstimatedLoss)
	default:
		reasons = fmt.Sprintf("An estimated loss of %.2f is %.0f%% of the %.2f ceiling on this claim, which is %s severity.",
			in.EstimatedLoss, 100*in.EstimatedLoss/in.MaxCoverageAmount, in.MaxCoverageAmount, severity)
	}
	if surveyRequired {
		reasons += " A loss of that share warrants a surveyor inspecting before assessment."
	}
	if len(missing) > 0 {
		reasons += fmt.Sprintf(" %s still missing from the evidence on file.", strings.Join(missing, " and "))
	}

	return ClassificationResult{
		Severity:         severity,
		SurveyRequired:   surveyRequired,
		MissingDocuments: missing,
		Reasons:          reasons,
		Source:           SourceDeterministic,
	}, nil
}

func severityFromLoss(in ClassificationInput) string {
	if in.MaxCoverageAmount <= 0 {
		return SeverityMedium
	}
	switch share := in.EstimatedLoss / in.MaxCoverageAmount; {
	case share >= highSeverityShare:
		return SeverityHigh
	case share >= mediumSeverityShare:
		return SeverityMedium
	default:
		return SeverityLow
	}
}
