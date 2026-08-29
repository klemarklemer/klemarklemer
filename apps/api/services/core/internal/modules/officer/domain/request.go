package domain

import (
	shareddomain "monorepo/services/core/pkg/shared/domain"
)

// RequestOfficer model
type RequestOfficer struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Email            string  `json:"email"`
	Role             string  `json:"role"`
	CurrentWorkload  int     `json:"current_workload"`
	MotorSkillRating float64 `json:"motor_skill_rating"`
	IsAvailable      bool    `json:"is_available"`
}

// Deserialize to db model
func (r *RequestOfficer) Deserialize() (res shareddomain.Officer) {
	res.ID = r.ID
	res.Name = r.Name
	res.Email = r.Email
	res.Role = r.Role
	res.CurrentWorkload = r.CurrentWorkload
	res.MotorSkillRating = r.MotorSkillRating
	res.IsAvailable = r.IsAvailable
	return
}
