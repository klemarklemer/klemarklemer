package core

import (
	"net/http"
	"time"

	"github.com/golangid/candi/candiutils"
)

type coreRESTImpl struct {
	host    string
	authKey string
	httpReq candiutils.HTTPRequest
}

// NewCoreServiceREST constructor
func NewCoreServiceREST(host string, authKey string) Core {

	return &coreRESTImpl{
		host:    host,
		authKey: authKey,
		httpReq: candiutils.NewHTTPRequest(
			candiutils.HTTPRequestSetRetries(5),
			candiutils.HTTPRequestSetSleepBetweenRetry(500*time.Millisecond),
			candiutils.HTTPRequestSetHTTPErrorCodeThreshold(http.StatusBadRequest),
			candiutils.HTTPRequestSetBreakerName("core"),
		),
	}
}
