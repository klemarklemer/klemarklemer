package domain

import (
	shareddomain "monorepo/services/core/pkg/shared/domain"
)

// RequestPolicy model
type RequestPolicy struct {
	ID                int     `json:"id"`
	PolicyNumber      string  `json:"policy_number"`
	PolicyHolderName  string  `json:"policy_holder_name"`
	VehiclePlate      string  `json:"vehicle_plate"`
	VehicleModel      string  `json:"vehicle_model"`
	CoverageType      string  `json:"coverage_type"`
	MaxCoverageAmount float64 `json:"max_coverage_amount"`
	DeductibleAmount  float64 `json:"deductible_amount"`
	Status            string  `json:"status"`
}

// Deserialize to db model
func (r *RequestPolicy) Deserialize() (res shareddomain.Policy) {
	res.ID = r.ID
	res.PolicyNumber = r.PolicyNumber
	res.PolicyHolderName = r.PolicyHolderName
	res.VehiclePlate = r.VehiclePlate
	res.VehicleModel = r.VehicleModel
	res.CoverageType = r.CoverageType
	res.MaxCoverageAmount = r.MaxCoverageAmount
	res.DeductibleAmount = r.DeductibleAmount
	res.Status = r.Status
	return
}
