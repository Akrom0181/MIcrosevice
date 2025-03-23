package handler

import (
	"net/http"
	"strconv"

	"user_api_gateway/config"
	"user_api_gateway/genproto/user_service"
	"user_api_gateway/pkg/etc"
	"user_api_gateway/pkg/helpers"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateUser godoc
// @Router        /user [post]
// @Summary       Create user
// @Description   API for creating user
// @Security      BearerAuth
// @Tags          user
// @Accept        json
// @Produce       json
// @Param         user body user_service.User true "User"
// @Success       200 {object} user_service.User
// @Failure       404 {object} user_service.ErrorResponse
// @Failure       500 {object} user_service.ErrorResponse
func (h *handler) CreateUser(ctx *gin.Context) {
	var (
		body user_service.User
		resp *user_service.User
	)

	if ctx.GetHeader("user_type") != "admin" {
		h.ReturnError(ctx, config.ErrorForbidden, "Invalid request body", 403)
		return
	}

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	if err := helpers.ValidateEmailAddress(body.Email); err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "invalid email address"+body.Email, 400)
		return
	}

	hashedPassword, err := etc.GeneratePasswordHash(body.Password)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Error hashing password", 400)
		return
	}

	body.Password = string(hashedPassword)

	resp, err = h.grpcClient.UserService().Create(ctx.Request.Context(), &body)
	if h.HandleDbError(ctx, err, "Error creating user") {
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// GetListUser godoc
// @Router         /user/list [GET]
// @Summary        Get list of users
// @Description    API for getting list of users
// @Security       BearerAuth
// @Tags           user
// @Accept         json
// @Produce        json
// @Param          page query int true "page"
// @Param          limit query int true "limit"
// @Param          search query string false "search"
// @Success        200 {object} user_service.GetListUserResponse
// @Failure        404 {object} user_service.ErrorResponse
// @Failure        500 {object} user_service.ErrorResponse
func (h *handler) GetUsers(ctx *gin.Context) {
	var (
		req  user_service.GetListUserRequest
		resp *user_service.GetListUserResponse
		err  error
	)

	if ctx.GetHeader("user_type") != "admin" {
		h.ReturnError(ctx, config.ErrorForbidden, "Invalid request body", 403)
		return
	}

	req.Search = ctx.Query("search")

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

	resp, err = h.grpcClient.UserService().GetList(ctx.Request.Context(), &req)
	if h.HandleDbError(ctx, err, "Error getting users") {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetUser godoc
// @Router         /user/{id} [GET]
// @Summary        Get a single user by ID
// @Description    API for getting a single user by ID
// @Security       BearerAuth
// @Tags           user
// @Accept         json
// @Produce        json
// @Success        200 {object} user_service.User
// @Failure        404 {object} user_service.ErrorResponse
// @Failure        500 {object} user_service.ErrorResponse
func (h *handler) GetUser(ctx *gin.Context) {
	var (
		id   = ctx.GetHeader("sub")
		resp *user_service.User
	)

	req := &user_service.UserSingleRequest{
		Id: id,
	}

	resp, err := h.grpcClient.UserService().GetSingle(ctx.Request.Context(), req)
	if h.HandleDbError(ctx, err, "Error getting user") {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// UpdateUser godoc
// @Router          /user [PUT]
// @Summary         Update a user by ID
// @Description     API for updating a user by ID
// @Security        BearerAuth
// @Tags            user
// @Accept          json
// @Produce         json
// @Param           user body user_service.User true "User"
// @Success         200 {object} user_service.User
// @Failure         404 {object} user_service.ErrorResponse
// @Failure         500 {object} user_service.ErrorResponse
func (h *handler) UpdateUser(ctx *gin.Context) {
	var (
		body user_service.User
		resp *user_service.User
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	if ctx.GetHeader("user_type") == "user" {
		body.Id = ctx.GetHeader("sub")
	}

	if body.Password != "" {
		hashedPassword, hashErr := etc.GeneratePasswordHash(body.Password)
		if hashErr != nil {
			h.ReturnError(ctx, config.ErrorBadRequest, "Error hashing password", 400)
			return
		}
		body.Password = string(hashedPassword)
	} else {
		h.ReturnError(ctx, config.ErrorBadRequest, "Password cannot be empty", 400)
		return
	}

	resp, err = h.grpcClient.UserService().Update(ctx.Request.Context(), &body)
	if h.HandleDbError(ctx, err, "Error updating user") {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// DeleteUser godoc
// @Router        /user/{id} [DELETE]
// @Summary       Delete a teacher by ID
// @Description   API for deleting a user by ID
// @Security      BearerAuth
// @Tags          user
// @Accept        json
// @Produce       json
// @Param         id path string true "User ID"
// @Success       200 {object} user_service.SuccessResponse
// @Failure       404 {object} user_service.ErrorResponse
// @Failure       500 {object} user_service.ErrorResponse
func (h *handler) DeleteUser(ctx *gin.Context) {
	var (
		id   = ctx.Param("id")
		err  error
		resp = &emptypb.Empty{}
	)

	if ctx.GetHeader("user_type") != "admin" {
		h.ReturnError(ctx, config.ErrorForbidden, "Invalid request body", 403)
		return
	}

	req := &user_service.UserPrimaryKey{
		Id: id,
	}

	resp, err = h.grpcClient.UserService().Delete(ctx.Request.Context(), req)
	if h.HandleDbError(ctx, err, "Error deleting user") {
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
