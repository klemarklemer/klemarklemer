package domain

import (
	"time"

	shareddomain "monorepo/services/core/pkg/shared/domain"

	"github.com/golangid/candi/candishared"
)

// ResponsePolicyList model
type ResponsePolicyList struct {
	Meta candishared.Meta `json:"meta"`
	Data []ResponsePolicy `json:"data"`
}

// ResponsePolicy model
type ResponsePolicy struct {
	ID                int     `json:"id"`
	PolicyNumber      string  `json:"policy_number"`
	PolicyHolderName  string  `json:"policy_holder_name"`
	VehiclePlate      string  `json:"vehicle_plate"`
	VehicleModel      string  `json:"vehicle_model"`
	CoverageType      string  `json:"coverage_type"`
	MaxCoverageAmount float64 `json:"max_coverage_amount"`
	DeductibleAmount  float64 `json:"deductible_amount"`
	Status            string  `json:"status"`
	EffectiveDate     string  `json:"effective_date"`
	ExpiryDate        string  `json:"expiry_date"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

// Serialize from db model
func (r *ResponsePolicy) Serialize(source *shareddomain.Policy) {
	r.ID = source.ID
	r.PolicyNumber = source.PolicyNumber
	r.PolicyHolderName = source.PolicyHolderName
	r.VehiclePlate = source.VehiclePlate
	r.VehicleModel = source.VehicleModel
	r.CoverageType = source.CoverageType
	r.MaxCoverageAmount = source.MaxCoverageAmount
	r.DeductibleAmount = source.DeductibleAmount
	r.Status = source.Status
	r.EffectiveDate = source.EffectiveDate.Format(time.RFC3339)
	r.ExpiryDate = source.ExpiryDate.Format(time.RFC3339)
	r.CreatedAt = source.CreatedAt.Format(time.RFC3339)
	r.UpdatedAt = source.UpdatedAt.Format(time.RFC3339)
}
