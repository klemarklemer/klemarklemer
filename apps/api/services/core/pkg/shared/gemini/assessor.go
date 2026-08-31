// Package gemini holds the reasoning layer behind the Assessment Agent.
//
// Without a credential it degrades to deterministic underwriting rules rather
// than failing, so a keyless CI run or fresh clone still reasons over real claim
// data and reports which engine answered instead of imitating a model.
package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"google.golang.org/genai"
)

// Mirrored from the shared claim domain to keep this package independently testable.
const (
	OutcomeApprove      = "APPROVE"
	OutcomeReject       = "REJECT"
	OutcomeManualReview = "MANUAL_REVIEW"
)

const (
	SourceDeterministic = "deterministic-rules"
	BackendVertexAI     = "vertex-ai"
	BackendGeminiAPI    = "gemini-api"
)

const (
	defaultModel = "gemini-2.5-flash"
	// Bounds one call so a slow response cannot hold a claim transaction open.
	requestTimeout = 25 * time.Second
	// Low, so the same claim yields a reproducible recommendation.
	assessmentTemperature = 0.2
)

// AssessmentInput is flat by design: the assessor never reaches back into the
// database, so anything absent here cannot influence the outcome.
type AssessmentInput struct {
	ClaimNumber         string
	ClaimType           string
	Severity            string
	IncidentDescription string
	EstimatedLoss       float64

	PolicyNumber      string
	CoverageType      string
	PolicyStatus      string
	MaxCoverageAmount float64
	DeductibleAmount  float64
	PolicyExpiry      time.Time

	DocumentTypes []string
}

// AssessmentResult is a recommendation, never a decision.
type AssessmentResult struct {
	Outcome    string
	Confidence float64
	Reasons    string
	Source     string
}

type Assessor interface {
	Assess(ctx context.Context, in AssessmentInput) (AssessmentResult, error)
}

// New selects a backend from the environment alone, so the same binary runs on
// Vertex AI (GOOGLE_GENAI_USE_VERTEXAI + GOOGLE_CLOUD_PROJECT, via ADC), on the
// Gemini API (GEMINI_API_KEY), or on deterministic rules. It never returns nil.
// The second value names the engine selected.
func New(ctx context.Context) (Assessor, string) {
	useVertex := vertexConfigured()
	apiKey := firstNonEmpty(os.Getenv("GEMINI_API_KEY"), os.Getenv("GOOGLE_API_KEY"))

	if !useVertex && apiKey == "" {
		return deterministicAssessor{}, SourceDeterministic
	}

	// Leaving Backend unset lets the SDK resolve Vertex from the environment;
	// setting it here would pin the Gemini API.
	cfg := &genai.ClientConfig{}
	if !useVertex {
		cfg.APIKey = apiKey
		cfg.Backend = genai.BackendGeminiAPI
	}

	client, err := genai.NewClient(ctx, cfg)
	if err != nil {
		return deterministicAssessor{}, SourceDeterministic
	}

	model := firstNonEmpty(os.Getenv("GEMINI_MODEL"), defaultModel)
	return &geminiAssessor{client: client, model: model, backend: backendLabel(useVertex)}, model + " via " + backendLabel(useVertex)
}

// Mirrors the SDK's truthiness rule so the two cannot disagree.
func vertexConfigured() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("GOOGLE_GENAI_USE_VERTEXAI"))) {
	case "1", "true":
		return os.Getenv("GOOGLE_CLOUD_PROJECT") != ""
	default:
		return false
	}
}

func backendLabel(useVertex bool) string {
	if useVertex {
		return BackendVertexAI
	}
	return BackendGeminiAPI
}

type geminiAssessor struct {
	client  *genai.Client
	model   string
	backend string
}

const systemInstruction = `You are the Assessment Agent in an insurance claims operations platform.

Given one motor claim and the policy behind it, recommend an outcome. You are advising a
human claims officer who makes the binding decision, so your job is to be correct and
auditable, not agreeable.

Apply these underwriting rules:
- A policy that is not in force (expired, or status other than ACTIVE) can never be APPROVE.
- An estimated loss above the policy's maximum coverage is never a clean APPROVE.
- An estimated loss at or below the deductible leaves nothing payable, so it is REJECT.
- Missing a POLICE_REPORT or a DAMAGE_PHOTO weakens the evidence; prefer MANUAL_REVIEW
  over APPROVE when substantiation is thin.
- Anything genuinely ambiguous is MANUAL_REVIEW. That is a valid, useful answer.

Cite the actual figures you were given - policy number, coverage ceiling, deductible,
estimated loss, and the settlement that would follow. Never invent a document, a reference
number, or an amount that is not in the input. Confidence is your genuine certainty from
0.0 to 1.0; reserve values above 0.9 for claims where every rule points the same way.`

func (a *geminiAssessor) Assess(ctx context.Context, in AssessmentInput) (AssessmentResult, error) {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	resp, err := a.client.Models.GenerateContent(ctx, a.model, genai.Text(buildPrompt(in)), &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(systemInstruction, genai.RoleUser),
		Temperature:       genai.Ptr[float32](assessmentTemperature),
		ResponseMIMEType:  "application/json",
		ResponseSchema:    responseSchema(),
	})
	if err != nil {
		return AssessmentResult{}, fmt.Errorf("gemini assessment call: %w", err)
	}

	var parsed struct {
		Outcome    string  `json:"outcome"`
		Confidence float64 `json:"confidence"`
		Reasons    string  `json:"reasons"`
	}
	if err := json.Unmarshal([]byte(resp.Text()), &parsed); err != nil {
		return AssessmentResult{}, fmt.Errorf("gemini assessment response was not valid JSON: %w", err)
	}

	outcome, err := normalizeOutcome(parsed.Outcome)
	if err != nil {
		return AssessmentResult{}, err
	}
	if strings.TrimSpace(parsed.Reasons) == "" {
		return AssessmentResult{}, fmt.Errorf("gemini assessment returned no reasoning")
	}

	return AssessmentResult{
		Outcome:    outcome,
		Confidence: clampConfidence(parsed.Confidence),
		Reasons:    strings.TrimSpace(parsed.Reasons),
		Source:     a.model + " via " + a.backend,
	}, nil
}

// Constrains the model so a malformed answer is rejected by the API rather than
// parsed defensively here.
func responseSchema() *genai.Schema {
	return &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"outcome": {
				Type:        genai.TypeString,
				Description: "One of APPROVE, REJECT, or MANUAL_REVIEW.",
			},
			"confidence": {
				Type:        genai.TypeNumber,
				Description: "Genuine certainty in the outcome, from 0.0 to 1.0.",
			},
			"reasons": {
				Type:        genai.TypeString,
				Description: "Auditable justification citing the actual policy and loss figures supplied.",
			},
		},
		Required: []string{"outcome", "confidence", "reasons"},
	}
}

func buildPrompt(in AssessmentInput) string {
	var b strings.Builder

	b.WriteString("Assess this claim.\n\nCLAIM\n")
	fmt.Fprintf(&b, "  Number:      %s\n", in.ClaimNumber)
	fmt.Fprintf(&b, "  Type:        %s\n", in.ClaimType)
	fmt.Fprintf(&b, "  Severity:    %s\n", in.Severity)
	fmt.Fprintf(&b, "  Est. loss:   %.2f\n", in.EstimatedLoss)
	fmt.Fprintf(&b, "  Incident:    %s\n", valueOrUnknown(in.IncidentDescription))

	b.WriteString("\nPOLICY\n")
	if in.PolicyNumber == "" {
		b.WriteString("  No policy record is attached to this claim.\n")
	} else {
		fmt.Fprintf(&b, "  Number:      %s\n", in.PolicyNumber)
		fmt.Fprintf(&b, "  Coverage:    %s\n", in.CoverageType)
		fmt.Fprintf(&b, "  Status:      %s\n", in.PolicyStatus)
		fmt.Fprintf(&b, "  Max cover:   %.2f\n", in.MaxCoverageAmount)
		fmt.Fprintf(&b, "  Deductible:  %.2f\n", in.DeductibleAmount)
		fmt.Fprintf(&b, "  Expires:     %s\n", in.PolicyExpiry.Format(time.DateOnly))
		fmt.Fprintf(&b, "  Settlement if approved: %.2f\n", settlement(in))
	}

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

// Applies the same underwriting rules without a model, so an unconfigured
// environment still reasons over the real claim.
type deterministicAssessor struct{}

func (deterministicAssessor) Assess(_ context.Context, in AssessmentInput) (AssessmentResult, error) {
	payable := settlement(in)

	switch {
	case in.PolicyNumber == "":
		return result(OutcomeManualReview, 0.50,
			fmt.Sprintf("Claim %s has no policy record attached, so coverage cannot be established. Routing to manual review.",
				in.ClaimNumber)), nil

	case !strings.EqualFold(in.PolicyStatus, "ACTIVE"):
		return result(OutcomeReject, 0.90,
			fmt.Sprintf("Policy %s is in status %s rather than ACTIVE at the time of assessment, so the claim is not covered.",
				in.PolicyNumber, in.PolicyStatus)), nil

	case !in.PolicyExpiry.IsZero() && in.PolicyExpiry.Before(time.Now()):
		return result(OutcomeReject, 0.90,
			fmt.Sprintf("Policy %s expired on %s, before this assessment, so the claim falls outside the in-force period.",
				in.PolicyNumber, in.PolicyExpiry.Format(time.DateOnly))), nil

	case in.EstimatedLoss > in.MaxCoverageAmount:
		return result(OutcomeManualReview, 0.75,
			fmt.Sprintf("Estimated loss of %.2f exceeds the %.2f ceiling on policy %s. The excess needs an underwriter, so this cannot be auto-approved.",
				in.EstimatedLoss, in.MaxCoverageAmount, in.PolicyNumber)), nil

	case payable <= 0:
		return result(OutcomeReject, 0.85,
			fmt.Sprintf("Estimated loss of %.2f does not exceed the %.2f deductible on policy %s, so no settlement is payable.",
				in.EstimatedLoss, in.DeductibleAmount, in.PolicyNumber)), nil

	case len(missingDocuments(in.DocumentTypes)) > 0:
		return result(OutcomeManualReview, 0.60,
			fmt.Sprintf("Claim %s is within the %.2f limit on policy %s and would settle at %.2f, but %s is missing from the evidence on file (%s). Substantiation is too thin to approve.",
				in.ClaimNumber, in.MaxCoverageAmount, in.PolicyNumber, payable,
				strings.Join(missingDocuments(in.DocumentTypes), " and "), documentList(in.DocumentTypes))), nil

	default:
		return result(OutcomeApprove, 0.88,
			fmt.Sprintf("Policy %s is in force with %s cover to %.2f. The estimated loss of %.2f sits within that limit and is substantiated by %s, settling at %.2f after the %.2f deductible.",
				in.PolicyNumber, in.CoverageType, in.MaxCoverageAmount, in.EstimatedLoss,
				documentList(in.DocumentTypes), payable, in.DeductibleAmount)), nil
	}
}

func result(outcome string, confidence float64, reasons string) AssessmentResult {
	return AssessmentResult{
		Outcome:    outcome,
		Confidence: confidence,
		Reasons:    reasons,
		Source:     SourceDeterministic,
	}
}

func settlement(in AssessmentInput) float64 {
	return math.Max(0, in.EstimatedLoss-in.DeductibleAmount)
}

// Returns the required document types actually absent, so the recommendation can
// name them instead of listing every requirement whether or not it is met.
func missingDocuments(docs []string) []string {
	present := make(map[string]bool, len(docs))
	for _, d := range docs {
		present[strings.ToUpper(strings.TrimSpace(d))] = true
	}

	var missing []string
	for _, required := range []string{"POLICE_REPORT", "DAMAGE_PHOTO"} {
		if !present[required] {
			missing = append(missing, required)
		}
	}
	return missing
}

func documentList(docs []string) string {
	if len(docs) == 0 {
		return "no supporting documents"
	}
	return strings.Join(docs, ", ")
}

// Rejects anything the claim record cannot store, so an invented outcome is
// refused rather than written to the row.
func normalizeOutcome(raw string) (string, error) {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case OutcomeApprove:
		return OutcomeApprove, nil
	case OutcomeReject:
		return OutcomeReject, nil
	case OutcomeManualReview:
		return OutcomeManualReview, nil
	default:
		return "", fmt.Errorf("gemini assessment returned unrecognised outcome %q", raw)
	}
}

func clampConfidence(c float64) float64 {
	if math.IsNaN(c) {
		return 0
	}
	return math.Min(1, math.Max(0, c))
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func valueOrUnknown(s string) string {
	if strings.TrimSpace(s) == "" {
		return "(not provided)"
	}
	return s
}
