package domain

import (
	"time"
)

// ClaimsOfficer domain entity
type ClaimsOfficer struct {
	ID               int       `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	Name             string    `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Email            string    `gorm:"column:email;type:varchar(255);uniqueIndex;not null" json:"email"`
	Role             string    `gorm:"column:role;type:varchar(64);not null;default:'Claims Officer'" json:"role"`
	Specialty        *string   `gorm:"column:specialty;type:varchar(64)" json:"specialty,omitempty"`
	Region           *string   `gorm:"column:region;type:varchar(64)" json:"region,omitempty"`
	CurrentWorkload  int       `gorm:"column:current_workload;not null;default:0" json:"current_workload"`
	MotorSkillRating float64   `gorm:"column:motor_skill_rating;type:decimal(3,2);not null;default:4.0" json:"motor_skill_rating"`
	IsAvailable      bool      `gorm:"column:is_available;not null;default:true" json:"is_available"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName return table name of ClaimsOfficer model
func (ClaimsOfficer) TableName() string {
	return "claims_officers"
}

// Officer alias for ClaimsOfficer
type Officer = ClaimsOfficer

// RoleSurveyor is the one officer role that inspects damage rather than owning a
// claim. Loop 2 scores claim ownership, so it must exclude these.
const RoleSurveyor = "Surveyor"

// IsSurveyor reports whether this officer inspects damage rather than owning claims.
func (o ClaimsOfficer) IsSurveyor() bool {
	return o.Role == RoleSurveyor
}
