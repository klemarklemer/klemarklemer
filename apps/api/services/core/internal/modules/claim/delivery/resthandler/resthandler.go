package resthandler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"monorepo/services/core/internal/modules/claim/domain"
	"monorepo/services/core/pkg/shared/usecase"

	"github.com/golangid/candi/candihelper"
	restserver "github.com/golangid/candi/codebase/app/rest_server"
	"github.com/golangid/candi/codebase/factory/dependency"
	"github.com/golangid/candi/codebase/interfaces"
	"github.com/golangid/candi/tracer"
	"github.com/golangid/candi/wrapper"
)

// RestHandler handler
type RestHandler struct {
	mw        interfaces.Middleware
	uc        usecase.Usecase
	validator interfaces.Validator
}

// NewRestHandler create new rest handler
func NewRestHandler(uc usecase.Usecase, deps dependency.Dependency) *RestHandler {
	return &RestHandler{
		uc: uc, mw: deps.GetMiddleware(), validator: deps.GetValidator(),
	}
}

// Mount handler with root "/"
func (h *RestHandler) Mount(root interfaces.RESTRouter) {
	v1Claim := root.Group(candihelper.V1 + "/claim")

	v1Claim.GET("/", h.getAllClaim)
	v1Claim.GET("/:id", h.getDetailClaimByID)
	v1Claim.POST("/", h.createClaim)
	v1Claim.POST("/:id/documents", h.uploadDocument)
	v1Claim.POST("/:id/intake", h.evaluateIntake)
	v1Claim.POST("/:id/assignment", h.runAssignment)
	v1Claim.POST("/:id/assessment", h.runAssessment)
	v1Claim.POST("/:id/approval", h.submitHumanApproval)

	// Demo utility route
	root.Group(candihelper.V1 + "/demo").POST("/reset", h.resetDemo)
}

func (h *RestHandler) getAllClaim(rw http.ResponseWriter, req *http.Request) {
	trace, ctx := tracer.StartTraceWithContext(req.Context(), "ClaimDeliveryREST:GetAllClaim")
	defer trace.Finish()

	var filter domain.FilterClaim
	if err := candihelper.ParseFromQueryParam(req.URL.Query(), &filter); err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, "Failed parse filter", err).JSON(rw)
		return
	}

	result, err := h.uc.Claim().GetAllClaim(ctx, &filter)
	if err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	response := wrapper.NewHTTPResponse(http.StatusOK, "Success", result.Data)
	response.Meta = result.Meta
	response.JSON(rw)
}

func (h *RestHandler) getDetailClaimByID(rw http.ResponseWriter, req *http.Request) {
	trace, ctx := tracer.StartTraceWithContext(req.Context(), "ClaimDeliveryREST:GetDetailClaimByID")
	defer trace.Finish()

	id, _ := strconv.Atoi(restserver.URLParam(req, "id"))
	data, err := h.uc.Claim().GetDetailClaim(ctx, id)
	if err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	wrapper.NewHTTPResponse(http.StatusOK, "Success", data).JSON(rw)
}

func (h *RestHandler) createClaim(rw http.ResponseWriter, req *http.Request) {
	trace, ctx := tracer.StartTraceWithContext(req.Context(), "ClaimDeliveryREST:CreateClaim")
	defer trace.Finish()

	body, _ := io.ReadAll(req.Body)
	var payload domain.RequestCreateClaim
	if err := json.Unmarshal(body, &payload); err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	res, err := h.uc.Claim().CreateClaim(ctx, &payload)
	if err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	wrapper.NewHTTPResponse(http.StatusCreated, "Success", res).JSON(rw)
}

func (h *RestHandler) uploadDocument(rw http.ResponseWriter, req *http.Request) {
	trace, ctx := tracer.StartTraceWithContext(req.Context(), "ClaimDeliveryREST:UploadDocument")
	defer trace.Finish()

	id, _ := strconv.Atoi(restserver.URLParam(req, "id"))
	body, _ := io.ReadAll(req.Body)

	var payload domain.RequestUploadDocument
	if len(body) > 0 {
		_ = json.Unmarshal(body, &payload)
	}

	res, err := h.uc.Claim().UploadDocument(ctx, id, &payload)
	if err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	wrapper.NewHTTPResponse(http.StatusOK, "Document uploaded and intake evaluated", res).JSON(rw)
}

func (h *RestHandler) evaluateIntake(rw http.ResponseWriter, req *http.Request) {
	trace, ctx := tracer.StartTraceWithContext(req.Context(), "ClaimDeliveryREST:EvaluateIntake")
	defer trace.Finish()

	id, _ := strconv.Atoi(restserver.URLParam(req, "id"))
	res, err := h.uc.Claim().EvaluateIntake(ctx, id)
	if err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	wrapper.NewHTTPResponse(http.StatusOK, "Intake evaluated successfully", res).JSON(rw)
}

func (h *RestHandler) runAssignment(rw http.ResponseWriter, req *http.Request) {
	trace, ctx := tracer.StartTraceWithContext(req.Context(), "ClaimDeliveryREST:RunAssignment")
	defer trace.Finish()

	id, _ := strconv.Atoi(restserver.URLParam(req, "id"))
	res, err := h.uc.Claim().RunAssignment(ctx, id)
	if err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	wrapper.NewHTTPResponse(http.StatusOK, "Assignment completed", res).JSON(rw)
}

func (h *RestHandler) runAssessment(rw http.ResponseWriter, req *http.Request) {
	trace, ctx := tracer.StartTraceWithContext(req.Context(), "ClaimDeliveryREST:RunAssessment")
	defer trace.Finish()

	id, _ := strconv.Atoi(restserver.URLParam(req, "id"))
	res, err := h.uc.Claim().RunAssessment(ctx, id)
	if err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	wrapper.NewHTTPResponse(http.StatusOK, "Assessment recommendation generated", res).JSON(rw)
}

func (h *RestHandler) submitHumanApproval(rw http.ResponseWriter, req *http.Request) {
	trace, ctx := tracer.StartTraceWithContext(req.Context(), "ClaimDeliveryREST:SubmitHumanApproval")
	defer trace.Finish()

	id, _ := strconv.Atoi(restserver.URLParam(req, "id"))
	body, _ := io.ReadAll(req.Body)

	var payload domain.RequestHumanApproval
	if err := json.Unmarshal(body, &payload); err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	res, err := h.uc.Claim().SubmitHumanApproval(ctx, id, &payload)
	if err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	wrapper.NewHTTPResponse(http.StatusOK, "Human approval processed successfully", res).JSON(rw)
}

func (h *RestHandler) resetDemo(rw http.ResponseWriter, req *http.Request) {
	trace, ctx := tracer.StartTraceWithContext(req.Context(), "ClaimDeliveryREST:ResetDemo")
	defer trace.Finish()

	res, err := h.uc.Claim().ResetDemo(ctx)
	if err != nil {
		wrapper.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(rw)
		return
	}

	wrapper.NewHTTPResponse(http.StatusOK, "Demo state reset successfully", res).JSON(rw)
}
