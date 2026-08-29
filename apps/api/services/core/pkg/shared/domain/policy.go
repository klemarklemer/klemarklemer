package domain

import (
	"time"
)

// Policy domain entity
type Policy struct {
	ID                int       `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	PolicyNumber      string    `gorm:"column:policy_number;type:varchar(64);uniqueIndex;not null" json:"policy_number"`
	PolicyHolderName  string    `gorm:"column:policy_holder_name;type:varchar(255);not null" json:"policy_holder_name"`
	VehiclePlate      string    `gorm:"column:vehicle_plate;type:varchar(32);not null" json:"vehicle_plate"`
	VehicleModel      string    `gorm:"column:vehicle_model;type:varchar(128);not null" json:"vehicle_model"`
	CoverageType      string    `gorm:"column:coverage_type;type:varchar(64);not null" json:"coverage_type"`
	MaxCoverageAmount float64   `gorm:"column:max_coverage_amount;type:decimal(15,2);not null" json:"max_coverage_amount"`
	DeductibleAmount  float64   `gorm:"column:deductible_amount;type:decimal(15,2);not null" json:"deductible_amount"`
	EffectiveDate     time.Time `gorm:"column:effective_date;not null" json:"effective_date"`
	ExpiryDate        time.Time `gorm:"column:expiry_date;not null" json:"expiry_date"`
	Status            string    `gorm:"column:status;type:varchar(32);not null;default:'ACTIVE'" json:"status"`
	CreatedAt         time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName return table name of Policy model
func (Policy) TableName() string {
	return "policies"
}
