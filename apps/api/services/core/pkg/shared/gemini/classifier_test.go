package gemini

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func classify(t *testing.T, in ClassificationInput) ClassificationResult {
	t.Helper()
	got, err := deterministicClassifier{}.Classify(context.Background(), in)
	assert.NoError(t, err)
	return got
}

func TestDeterministicClassify_GivenLossAboveHalfTheCeiling_IsHighAndNeedsASurvey(t *testing.T) {
	got := classify(t, ClassificationInput{
		ClaimNumber: "CLM-1", EstimatedLoss: 30000, MaxCoverageAmount: 45000,
		DocumentTypes: []string{"POLICE_REPORT", "DAMAGE_PHOTO"},
	})

	assert.Equal(t, SeverityHigh, got.Severity)
	assert.True(t, got.SurveyRequired, "a loss consuming most of the cover warrants an inspection")
	assert.Contains(t, got.Reasons, "surveyor")
}

func TestDeterministicClassify_GivenModestLoss_IsLowAndNeedsNoSurvey(t *testing.T) {
	got := classify(t, ClassificationInput{
		ClaimNumber: "CLM-2", EstimatedLoss: 1000, MaxCoverageAmount: 45000,
		DocumentTypes: []string{"POLICE_REPORT", "DAMAGE_PHOTO"},
	})

	assert.Equal(t, SeverityLow, got.Severity)
	assert.False(t, got.SurveyRequired)
}

func TestDeterministicClassify_GivenMidRangeLoss_IsMedium(t *testing.T) {
	got := classify(t, ClassificationInput{
		ClaimNumber: "CLM-3", EstimatedLoss: 9000, MaxCoverageAmount: 45000,
		DocumentTypes: []string{"POLICE_REPORT", "DAMAGE_PHOTO"},
	})

	assert.Equal(t, SeverityMedium, got.Severity)
	assert.False(t, got.SurveyRequired)
}

func TestDeterministicClassify_GivenNoPolicyCeiling_FallsBackToMedium(t *testing.T) {
	got := classify(t, ClassificationInput{ClaimNumber: "CLM-4", EstimatedLoss: 9000})

	assert.Equal(t, SeverityMedium, got.Severity)
	assert.Contains(t, got.Reasons, "no policy ceiling")
}

// The old intake wrote missing_documents as a constant ["POLICE_REPORT"], so a
// claim missing only the photo was reported as missing the wrong document.
func TestDeterministicClassify_NamesTheDocumentsActuallyAbsent(t *testing.T) {
	onlyPhotoMissing := classify(t, ClassificationInput{
		ClaimNumber: "CLM-5", EstimatedLoss: 1000, MaxCoverageAmount: 45000,
		DocumentTypes: []string{"POLICE_REPORT"},
	})
	assert.Equal(t, []string{"DAMAGE_PHOTO"}, onlyPhotoMissing.MissingDocuments)

	nothingOnFile := classify(t, ClassificationInput{
		ClaimNumber: "CLM-6", EstimatedLoss: 1000, MaxCoverageAmount: 45000,
	})
	assert.Equal(t, []string{"POLICE_REPORT", "DAMAGE_PHOTO"}, nothingOnFile.MissingDocuments)

	complete := classify(t, ClassificationInput{
		ClaimNumber: "CLM-7", EstimatedLoss: 1000, MaxCoverageAmount: 45000,
		DocumentTypes: []string{"DAMAGE_PHOTO", "POLICE_REPORT"},
	})
	assert.Empty(t, complete.MissingDocuments)
}

func TestDeterministicClassify_AlwaysReportsItsEngine(t *testing.T) {
	got := classify(t, ClassificationInput{ClaimNumber: "CLM-8", EstimatedLoss: 500, MaxCoverageAmount: 45000})

	assert.Equal(t, SourceDeterministic, got.Source)
	assert.NotEmpty(t, got.Reasons)
}

func TestNormalizeSeverity(t *testing.T) {
	for _, in := range []string{"high", " HIGH ", "High"} {
		got, err := normalizeSeverity(in)
		assert.NoError(t, err)
		assert.Equal(t, SeverityHigh, got)
	}

	_, err := normalizeSeverity("CATASTROPHIC")
	assert.Error(t, err, "an unknown severity must be rejected, not silently defaulted")
}
