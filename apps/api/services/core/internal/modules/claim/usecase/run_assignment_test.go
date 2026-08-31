package usecase

import (
	"testing"

	shareddomain "monorepo/services/core/pkg/shared/domain"

	"github.com/stretchr/testify/assert"
)

func officer(name, role string, workload int, skill float64, available bool) shareddomain.Officer {
	return shareddomain.Officer{
		Name: name, Role: role, CurrentWorkload: workload,
		MotorSkillRating: skill, IsAvailable: available,
	}
}

func surveyor(name, specialty string, workload int, available bool) shareddomain.Officer {
	o := officer(name, shareddomain.RoleSurveyor, workload, 4.5, available)
	o.Specialty = &specialty
	return o
}

// The seeded roster: every surveyor outscores every claims officer, so an
// unfiltered scorer hands a motor claim to a property inspector.
func seededRoster() []shareddomain.Officer {
	return []shareddomain.Officer{
		officer("Alex Rivera", "Claims Officer", 4, 4.80, true),
		officer("David Chen", "Senior Claims Officer", 8, 4.90, true),
		officer("Elena Rostova", "Claims Officer", 2, 4.20, true),
		officer("Marcus Webb", shareddomain.RoleSurveyor, 1, 4.70, true),
		officer("Priya Sharma", shareddomain.RoleSurveyor, 0, 4.50, true),
		officer("James Okafor", shareddomain.RoleSurveyor, 2, 4.60, true),
	}
}

func TestSelectClaimOwner_GivenSeededRoster_NeverPicksASurveyor(t *testing.T) {
	best, _, _, _ := selectClaimOwner(seededRoster())

	assert.NotNil(t, best)
	assert.False(t, best.IsSurveyor(), "Loop 2 assigned a surveyor to own a claim")
	assert.Equal(t, "Elena Rostova", best.Name, "best claims officer is the one with the lowest workload")
}

func TestSelectClaimOwner_GivenSurveyorScoresHighest_StillPicksTheOfficer(t *testing.T) {
	roster := []shareddomain.Officer{
		officer("Priya Sharma", shareddomain.RoleSurveyor, 0, 5.00, true),
		officer("Alex Rivera", "Claims Officer", 9, 4.00, true),
	}

	best, _, _, _ := selectClaimOwner(roster)

	assert.NotNil(t, best)
	assert.Equal(t, "Alex Rivera", best.Name)
}

func TestSelectClaimOwner_GivenOnlyUnavailableOfficers_ReturnsNil(t *testing.T) {
	roster := []shareddomain.Officer{
		officer("Alex Rivera", "Claims Officer", 1, 4.80, false),
		officer("Priya Sharma", shareddomain.RoleSurveyor, 0, 4.50, true),
	}

	best, workload, skill, total := selectClaimOwner(roster)

	assert.Nil(t, best, "an unavailable officer must not be assigned as a fallback")
	assert.Zero(t, workload)
	assert.Zero(t, skill)
	assert.Zero(t, total)
}

func TestSelectClaimOwner_GivenNoOfficers_ReturnsNil(t *testing.T) {
	best, _, _, _ := selectClaimOwner(nil)

	assert.Nil(t, best)
}

func TestSelectClaimOwner_GivenSeveralOfficers_PicksTheHighestCombinedScore(t *testing.T) {
	roster := []shareddomain.Officer{
		officer("Heavy Load", "Claims Officer", 9, 5.00, true),
		officer("Balanced", "Claims Officer", 2, 4.20, true),
	}

	best, workload, skill, total := selectClaimOwner(roster)

	assert.Equal(t, "Balanced", best.Name)
	assert.InDelta(t, 4.0, workload, 0.001)
	assert.InDelta(t, 4.2, skill, 0.001)
	assert.InDelta(t, 8.2, total, 0.001)
}

func TestSelectSurveyor_GivenAMotorClaim_PrefersAVehicleInspector(t *testing.T) {
	roster := []shareddomain.Officer{
		officer("Elena Rostova", "Claims Officer", 0, 4.20, true),
		surveyor("Priya Sharma", "Property Inspection", 0, true),
		surveyor("Marcus Webb", "Vehicle Inspection", 1, true),
		surveyor("James Okafor", "Vehicle Inspection", 2, true),
	}

	best := selectSurveyor(roster, "MOTOR")

	assert.NotNil(t, best)
	assert.True(t, best.IsSurveyor())
	assert.Equal(t, "Marcus Webb", best.Name,
		"a property inspector must not be sent to a car crash, even carrying less work")
}

func TestSelectSurveyor_GivenNoFittingSpecialty_FallsBackToLeastLoaded(t *testing.T) {
	roster := []shareddomain.Officer{
		surveyor("Priya Sharma", "Property Inspection", 3, true),
		surveyor("Ana Costa", "Marine Inspection", 1, true),
	}

	best := selectSurveyor(roster, "MOTOR")

	assert.NotNil(t, best)
	assert.Equal(t, "Ana Costa", best.Name, "no vehicle inspector free, so the least loaded surveyor takes it")
}

func TestSelectSurveyor_GivenOnlyClaimsOfficers_ReturnsNil(t *testing.T) {
	roster := []shareddomain.Officer{
		officer("Alex Rivera", "Claims Officer", 0, 5.00, true),
	}

	assert.Nil(t, selectSurveyor(roster, "MOTOR"), "a claims officer must never be sent to inspect")
}

func TestSelectSurveyor_GivenAllSurveyorsBusy_ReturnsNil(t *testing.T) {
	roster := []shareddomain.Officer{
		officer("Marcus Webb", shareddomain.RoleSurveyor, 1, 4.70, false),
		officer("Priya Sharma", shareddomain.RoleSurveyor, 0, 4.50, false),
	}

	assert.Nil(t, selectSurveyor(roster, "MOTOR"))
}
