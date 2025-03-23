package handler

import (
	"net/http"
	"strconv"
	"user_api_gateway/config"
	"user_api_gateway/genproto/user_service"

	"github.com/gin-gonic/gin"
)

// GetSession godoc
// @Router /session/{id} [get]
// @Summary Get a session by ID
// @Description Get a session by ID
// @Security BearerAuth
// @Tags session
// @Accept  json
// @Produce  json
// @Param id path string true "Session ID"
// @Success 200 {object} user_service.Session
// @Failure 400 {object} user_service.ErrorResponse
func (h *handler) GetSession(ctx *gin.Context) {
	var (
		req = &user_service.SessionSingleRequest{}
	)

	req.Id = ctx.Param("id")

	session, err := h.grpcClient.SessionService().GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting session") {
		return
	}

	ctx.JSON(200, session)
}

// GetSessions godoc
// @Router /session/list [get]
// @Summary Get a list of users
// @Description Get a list of users
// @Security BearerAuth
// @Tags session
// @Accept  json
// @Produce  json
// @Param page query number true "page"
// @Param limit query number true "limit"
// @Param user_id query string false "user_id"
// @Success 200 {object} user_service.GetListSessionResponse
// @Failure 400 {object} user_service.ErrorResponse
func (h *handler) GetSessions(ctx *gin.Context) {
	var (
		req  user_service.GetListSessionRequest
		resp *user_service.GetListSessionResponse
		err  error
	)

	req.Search = ctx.Query("user_id")

	if ctx.GetHeader("user_type") == "user" {
		req.Search = ctx.GetHeader("sub")
	}

	page, err := strconv.ParseUint(ctx.DefaultQuery("page", "1"), 10, 64)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid page", 400)
		return
	}

	limit, err := strconv.ParseUint(ctx.DefaultQuery("limit", "10"), 10, 64)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid limit", 400)
		return
	}

	req.Page = page
	req.Limit = limit

	resp, err = h.grpcClient.SessionService().GetList(ctx.Request.Context(), &req)
	if h.HandleDbError(ctx, err, "Error getting session") {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateSession godoc
// @Router /session [put]
// @Summary Update a session
// @Description Update a session
// @Security BearerAuth
// @Tags session
// @Accept  json
// @Produce  json
// @Param session body user_service.Session true "Session object"
// @Success 200 {object} user_service.Session
// @Failure 400 {object} user_service.ErrorResponse
func (h *handler) UpdateSession(ctx *gin.Context) {
	var (
		body *user_service.Session
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	session, err := h.grpcClient.SessionService().Update(ctx, body)
	if h.HandleDbError(ctx, err, "Error updating session") {
		return
	}

	ctx.JSON(200, session)
}

// DeleteSession godoc
// @Router /session/{id} [delete]
// @Summary Delete a session
// @Description Delete a session
// @Security BearerAuth
// @Tags session
// @Accept  json
// @Produce  json
// @Param id path string true "Session ID"
// @Success 200 {object} user_service.SuccessResponse
// @Failure 400 {object} user_service.ErrorResponse
func (h *handler) DeleteSession(ctx *gin.Context) {
	var (
		req = &user_service.SessionSingleRequest{}
	)

	req.Id = ctx.Param("id")

	_, err := h.grpcClient.SessionService().Delete(ctx, req)
	if h.HandleDbError(ctx, err, "Error deleting session") {
		return
	}

	ctx.JSON(200, user_service.SuccessResponse{
		Message: "Session deleted successfully",
	})
}
