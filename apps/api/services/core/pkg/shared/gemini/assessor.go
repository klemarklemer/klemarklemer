// Package gemini holds the reasoning layer behind the Assessment Agent.
//
// The agent grounds every recommendation in the claim record it is given: the
// policy's coverage ceiling, its deductible, its in-force window, the estimated
// loss, and the documents on file. Nothing about the outcome is pre-decided.
//
// The package degrades on purpose. When no Gemini credential is configured the
// constructor returns a deterministic assessor that applies the same underwriting
// rules locally, so the service still boots, still reasons over real claim data,
// and still reports which engine produced the answer. A demo or a CI run without
// a key therefore behaves honestly instead of pretending to be a model.
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

// Recommendation outcomes. Mirrored from the shared claim domain so this package
// stays free of a dependency on it and can be unit tested in isolation.
const (
	OutcomeApprove      = "APPROVE"
	OutcomeReject       = "REJECT"
	OutcomeManualReview = "MANUAL_REVIEW"
)

// Engine identifiers reported on every result so an auditor can tell which
// reasoning path produced a recommendation.
const (
	SourceDeterministic = "deterministic-rules"
)

const (
	// defaultModel is the Gemini model used when GEMINI_MODEL is unset.
	defaultModel = "gemini-3.5-flash"
	// requestTimeout bounds a single assessment call so one slow response cannot
	// hold a claim transaction open.
	requestTimeout = 25 * time.Second
	// assessmentTemperature keeps recommendations reproducible for the same claim.
	assessmentTemperature = 0.2
)

// AssessmentInput is the evidence an assessment reasons over. It is a flat value
// deliberately: the assessor must not reach back into the database, so whatever
// is not on this struct cannot influence the outcome.
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

// AssessmentResult is a recommendation, never a decision. A human still binds it
// through the decision gate.
type AssessmentResult struct {
	Outcome    string
	Confidence float64
	Reasons    string
	Source     string
}

// Assessor produces a grounded recommendation for one claim.
type Assessor interface {
	Assess(ctx context.Context, in AssessmentInput) (AssessmentResult, error)
}

// New returns a Gemini-backed assessor when a credential is configured, and the
// deterministic assessor otherwise. It never returns a nil Assessor, so callers
// do not need a nil branch: an unconfigured deployment degrades rather than fails.
func New(ctx context.Context) (Assessor, string) {
	apiKey := firstNonEmpty(os.Getenv("GEMINI_API_KEY"), os.Getenv("GOOGLE_API_KEY"))
	if apiKey == "" {
		return deterministicAssessor{}, SourceDeterministic
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return deterministicAssessor{}, SourceDeterministic
	}

	model := firstNonEmpty(os.Getenv("GEMINI_MODEL"), defaultModel)
	return &geminiAssessor{client: client, model: model}, model
}

type geminiAssessor struct {
	client *genai.Client
	model  string
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
		Source:     a.model,
	}, nil
}

// responseSchema constrains the model to the exact shape the claim record needs,
// so a malformed answer is rejected by the API rather than parsed defensively here.
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

// deterministicAssessor applies the same underwriting rules without a model. It
// exists so an unconfigured environment still reasons over the real claim rather
// than emitting a canned answer.
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

	case !hasRequiredDocuments(in.DocumentTypes):
		return result(OutcomeManualReview, 0.60,
			fmt.Sprintf("Claim %s is within the %.2f limit on policy %s and would settle at %.2f, but the evidence on file (%s) lacks a police report or damage photo. Substantiation is too thin to approve.",
				in.ClaimNumber, in.MaxCoverageAmount, in.PolicyNumber, payable, documentList(in.DocumentTypes))), nil

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

// settlement is what would actually be paid: the loss net of the deductible,
// floored at zero.
func settlement(in AssessmentInput) float64 {
	return math.Max(0, in.EstimatedLoss-in.DeductibleAmount)
}

func hasRequiredDocuments(docs []string) bool {
	var police, photo bool
	for _, d := range docs {
		switch strings.ToUpper(strings.TrimSpace(d)) {
		case "POLICE_REPORT":
			police = true
		case "DAMAGE_PHOTO":
			photo = true
		}
	}
	return police && photo
}

func documentList(docs []string) string {
	if len(docs) == 0 {
		return "no supporting documents"
	}
	return strings.Join(docs, ", ")
}

// normalizeOutcome rejects anything the claim record cannot store, so a model that
// invents an outcome fails loudly instead of writing an unknown value to the row.
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
