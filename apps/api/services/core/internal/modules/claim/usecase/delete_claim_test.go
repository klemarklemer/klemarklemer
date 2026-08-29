package usecase

import "testing"

func TestSkip_MissingMocks(t *testing.T) {
	t.Skip("mocks not yet generated; integration tests require running Postgres")
}
