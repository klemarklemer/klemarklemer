package domain

import (
	"time"

	shareddomain "monorepo/services/core/pkg/shared/domain"

	"github.com/golangid/candi/candishared"
)

// ResponseOfficerList model
type ResponseOfficerList struct {
	Meta candishared.Meta  `json:"meta"`
	Data []ResponseOfficer `json:"data"`
}

// ResponseOfficer model
type ResponseOfficer struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Email            string  `json:"email"`
	Role             string  `json:"role"`
	CurrentWorkload  int     `json:"current_workload"`
	MotorSkillRating float64 `json:"motor_skill_rating"`
	IsAvailable      bool    `json:"is_available"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

// Serialize from db model
func (r *ResponseOfficer) Serialize(source *shareddomain.Officer) {
	r.ID = source.ID
	r.Name = source.Name
	r.Email = source.Email
	r.Role = source.Role
	r.CurrentWorkload = source.CurrentWorkload
	r.MotorSkillRating = source.MotorSkillRating
	r.IsAvailable = source.IsAvailable
	r.CreatedAt = source.CreatedAt.Format(time.RFC3339)
	r.UpdatedAt = source.UpdatedAt.Format(time.RFC3339)
}
