package gemini

import (
	"context"
	"strings"
	"testing"
	"time"
)

// baseInput is a clean, fully substantiated claim that should approve. Each test
// below mutates exactly one field, so a differing outcome is attributable to that
// field and nothing else.
func baseInput() AssessmentInput {
	return AssessmentInput{
		ClaimNumber:         "CLM-2026-0001",
		ClaimType:           "MOTOR",
		Severity:            "MEDIUM",
		IncidentDescription: "Single-vehicle collision with a road barrier.",
		EstimatedLoss:       4200,
		PolicyNumber:        "POL-MOTOR-2026-0001",
		CoverageType:        "Comprehensive",
		PolicyStatus:        "ACTIVE",
		MaxCoverageAmount:   45000,
		DeductibleAmount:    500,
		PolicyExpiry:        time.Now().Add(365 * 24 * time.Hour),
		DocumentTypes:       []string{"POLICE_REPORT", "DAMAGE_PHOTO"},
	}
}

func TestDeterministicAssessor_GivenVaryingClaims_ReturnsOutcomeGroundedInInput(t *testing.T) {
	cases := []struct {
		name          string
		mutate        func(*AssessmentInput)
		wantOutcome   string
		wantInReasons string
	}{
		{
			name:          "in-force policy, loss within cover, documents complete",
			mutate:        func(*AssessmentInput) {},
			wantOutcome:   OutcomeApprove,
			wantInReasons: "3700.00", // 4200 loss - 500 deductible
		},
		{
			name:          "loss above the coverage ceiling cannot auto-approve",
			mutate:        func(in *AssessmentInput) { in.EstimatedLoss = 90000 },
			wantOutcome:   OutcomeManualReview,
			wantInReasons: "exceeds",
		},
		{
			name:          "loss at or under the deductible pays nothing",
			mutate:        func(in *AssessmentInput) { in.EstimatedLoss = 400 },
			wantOutcome:   OutcomeReject,
			wantInReasons: "deductible",
		},
		{
			name:          "lapsed policy is not covered",
			mutate:        func(in *AssessmentInput) { in.PolicyStatus = "LAPSED" },
			wantOutcome:   OutcomeReject,
			wantInReasons: "LAPSED",
		},
		{
			name:          "policy expired before assessment is not covered",
			mutate:        func(in *AssessmentInput) { in.PolicyExpiry = time.Now().Add(-24 * time.Hour) },
			wantOutcome:   OutcomeReject,
			wantInReasons: "expired",
		},
		{
			name:          "thin evidence goes to a human rather than approving",
			mutate:        func(in *AssessmentInput) { in.DocumentTypes = []string{"DAMAGE_PHOTO"} },
			wantOutcome:   OutcomeManualReview,
			wantInReasons: "police report",
		},
		{
			name:          "no policy attached means coverage is unestablished",
			mutate:        func(in *AssessmentInput) { in.PolicyNumber = "" },
			wantOutcome:   OutcomeManualReview,
			wantInReasons: "no policy record",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			in := baseInput()
			tc.mutate(&in)

			got, err := deterministicAssessor{}.Assess(context.Background(), in)
			if err != nil {
				t.Fatalf("Assess returned error: %v", err)
			}
			if got.Outcome != tc.wantOutcome {
				t.Errorf("outcome = %q, want %q\nreasons: %s", got.Outcome, tc.wantOutcome, got.Reasons)
			}
			if !strings.Contains(strings.ToLower(got.Reasons), strings.ToLower(tc.wantInReasons)) {
				t.Errorf("reasons did not cite %q\ngot: %s", tc.wantInReasons, got.Reasons)
			}
			if got.Confidence < 0 || got.Confidence > 1 {
				t.Errorf("confidence %v out of range", got.Confidence)
			}
			if got.Source != SourceDeterministic {
				t.Errorf("source = %q, want %q", got.Source, SourceDeterministic)
			}
		})
	}
}

// The defect this replaces returned an identical APPROVE/0.94 for every claim.
// This pins that two materially different claims can no longer share an answer.
func TestDeterministicAssessor_GivenDifferentClaims_DoesNotReturnIdenticalRecommendation(t *testing.T) {
	covered := baseInput()

	uncovered := baseInput()
	uncovered.EstimatedLoss = 250 // below the deductible

	a, err := deterministicAssessor{}.Assess(context.Background(), covered)
	if err != nil {
		t.Fatalf("Assess(covered): %v", err)
	}
	b, err := deterministicAssessor{}.Assess(context.Background(), uncovered)
	if err != nil {
		t.Fatalf("Assess(uncovered): %v", err)
	}

	if a.Outcome == b.Outcome && a.Reasons == b.Reasons {
		t.Fatalf("a payable and an unpayable claim produced the same recommendation: %+v", a)
	}
}

func TestSettlement_GivenLossBelowDeductible_FloorsAtZero(t *testing.T) {
	in := baseInput()
	in.EstimatedLoss = 100
	in.DeductibleAmount = 500

	if got := settlement(in); got != 0 {
		t.Errorf("settlement = %v, want 0", got)
	}
}

func TestNormalizeOutcome_GivenUnknownValue_ReturnsError(t *testing.T) {
	if _, err := normalizeOutcome("MAYBE"); err == nil {
		t.Error("expected an error for an outcome the claim record cannot store")
	}
	for _, valid := range []string{"approve", " REJECT ", "Manual_Review"} {
		if _, err := normalizeOutcome(valid); err != nil {
			t.Errorf("normalizeOutcome(%q) errored: %v", valid, err)
		}
	}
}

func TestClampConfidence_GivenOutOfRangeValues_ConstrainsToUnitInterval(t *testing.T) {
	for _, tc := range []struct{ in, want float64 }{{-3, 0}, {0.42, 0.42}, {7, 1}} {
		if got := clampConfidence(tc.in); got != tc.want {
			t.Errorf("clampConfidence(%v) = %v, want %v", tc.in, got, tc.want)
		}
	}
}

func TestNew_GivenNoCredential_FallsBackToDeterministicEngine(t *testing.T) {
	t.Setenv("GEMINI_API_KEY", "")
	t.Setenv("GOOGLE_API_KEY", "")

	assessor, engine := New(context.Background())
	if assessor == nil {
		t.Fatal("New returned a nil Assessor; an unconfigured deployment must still boot")
	}
	if engine != SourceDeterministic {
		t.Errorf("engine = %q, want %q", engine, SourceDeterministic)
	}
}
