package domain

import "github.com/golangid/candi/candishared"

// FilterOfficer model
type FilterOfficer struct {
	candishared.Filter
	ID        *int `json:"id"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Preloads  []string `json:"-"`
}
