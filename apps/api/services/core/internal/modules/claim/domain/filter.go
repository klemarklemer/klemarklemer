package domain

import "github.com/golangid/candi/candishared"

// FilterClaim model
type FilterClaim struct {
	candishared.Filter
	ID          *int     `json:"id"`
	ClaimNumber string   `json:"claimNumber"`
	Stage       string   `json:"stage"`
	Status      string   `json:"status"`
	StartDate   string   `json:"startDate"`
	EndDate     string   `json:"endDate"`
	Preloads    []string `json:"-"`
}
