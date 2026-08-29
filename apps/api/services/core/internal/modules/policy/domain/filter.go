package domain

import "github.com/golangid/candi/candishared"

// FilterPolicy model
type FilterPolicy struct {
	candishared.Filter
	ID        *int `json:"id"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Preloads  []string `json:"-"`
}
