package handler

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"user_api_gateway/config"
	"user_api_gateway/genproto/user_service"
	"user_api_gateway/pkg/etc"
	"user_api_gateway/pkg/helpers"
	"user_api_gateway/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @Router /auth/login [post]
// @Summary Login
// @Description Login
// @Tags auth
// @Accept  json
// @Produce  json
// @Param body body user_service.LoginRequest true "User"
// @Success 200 {object} user_service.SuccessResponse
// @Failure 400 {object} user_service.ErrorResponse
func (h *handler) Login(ctx *gin.Context) {
	var (
		body    user_service.LoginRequest
		session *user_service.Session
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	user, err := h.grpcClient.UserService().GetSingle(ctx, &user_service.UserSingleRequest{
		Username: body.Username,
		Email:    body.Email,
	})

	if err != nil {
		h.ReturnError(ctx, config.ErrorNotFound, "User not found", 404)
		return
	}

	if user.UserType == "user" && body.Platform == "admin" {
		h.ReturnError(ctx, config.ErrorForbidden, "User can't login to admin web", 400)
		return
	} else if user.UserType == "admin" && body.Platform != "admin" {
		h.ReturnError(ctx, config.ErrorForbidden, "Admin can only login to admin web", 400)
		return
	}

	if !etc.CheckPasswordHash(body.Password, user.Password) {
		h.ReturnError(ctx, config.ErrorInvalidPass, "Incorrect password", http.StatusBadRequest)
		return
	}

	// create session
	newSession := &user_service.Session{
		UserId:       user.Id,
		UserAgent:    ctx.Request.UserAgent(),
		IsActive:     true,
		IpAddress:    ctx.ClientIP(),
		ExpiresAt:    time.Now().Add(time.Hour * 999999).Format(time.RFC3339),
		LastActiveAt: time.Now().Format(time.RFC3339),
		Platform:     body.Platform,
	}

	session, err = h.grpcClient.SessionService().Create(ctx, newSession)
	if h.HandleDbError(ctx, err, "Error while creating new session") {
		return
	}

	// generate jwt token
	jwtFields := map[string]interface{}{
		"sub":        user.Id,
		"user_role":  user.UserRole,
		"user_type":  user.UserType,
		"platform":   body.Platform,
		"session_id": session.Id,
	}

	user.AccessToken, err = jwt.GenerateJWT(jwtFields, h.cfg.JWT)
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, "Oops, something went wrong!!!", http.StatusInternalServerError)
		return
	}

	ctx.JSON(200, gin.H{
		"user":    user,
		"session": session,
	})
}

// Logout godoc
// @Router      /auth/logout [post]
// @Summary     Logout
// @Description Logout
// @Security    BearerAuth
// @Tags        auth
// @Accept      json
// @Produce     json
// @Success     200 {object} user_service.SuccessResponse
// @Failure     400 {object} user_service.ErrorResponse
func (h *handler) Logout(ctx *gin.Context) {

	sessionID := ctx.GetHeader("session_id")
	if sessionID == "" {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid session ID", 400)
		return
	}

	_, err := h.grpcClient.SessionService().Delete(ctx, &user_service.SessionSingleRequest{
		Id: sessionID,
	})
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, "Error deleting session", 500)
		return
	}

	ctx.JSON(200, user_service.SuccessResponse{
		Message: "Successfully logged out",
	})
}

// Register godoc
// @Router /auth/register [post]
// @Summary Register
// @Description Register
// @Tags auth
// @Accept  json
// @Produce  json
// @Param body body user_service.RegisterRequest true "User"
// @Success 200 {object} user_service.User
// @Failure 400 {object} user_service.ErrorResponse
func (h *handler) Register(ctx *gin.Context) {
	var (
		body *user_service.RegisterRequest
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	_, err = h.grpcClient.UserService().GetSingle(ctx, &user_service.UserSingleRequest{
		Username: body.Username,
		Email:    body.Email,
	})

	if err == nil {
		h.ReturnError(ctx, config.ErrorConflict, "User already exists", 409)
		return
	}

	if err := helpers.ValidateEmailAddress(body.Email); err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid email address", 400)
		return
	}

	if err := helpers.ValidatePassword(body.Password); err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid password", 400)
		return
	}

	hashedPassword, err := etc.GeneratePasswordHash(body.Password)
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, "Oops, something went wrong!!!", http.StatusInternalServerError)
		return
	}
	body.Password = string(hashedPassword)

	user, err := h.grpcClient.UserService().Create(ctx, &user_service.User{
		FullName: body.Fullname,
		UserType: "user",
		UserRole: "user",
		UserName: body.Username,
		Email:    body.Email,
		Status:   "inverify",
		Password: body.Password,
		Gender:   body.Gender,
	})
	if h.HandleDbError(ctx, err, "Error creating user") {
		return
	}

	// send verification code to user
	otp := etc.GenerateOTP(6)
	err = h.redis.Set(ctx, fmt.Sprintf("otp-%s", user.Email), otp, 5*60)
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, "Error setting OTP", 500)
		return
	}

	// send otp code to user's email
	emailBody, err := etc.GenerateOtpEmailBody(otp)
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, "Error sending OTP", 500)
		return
	}

	err = etc.SendEmail(os.Getenv("GMAIL_HOST"), os.Getenv("GMAIL_PORT"), os.Getenv("GMAIL_USER"), os.Getenv("GMAIL_PASSWORD"), user.Email, emailBody)
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, "Error ssending OTP", 500)
		return
	}

	ctx.JSON(201, user_service.SuccessResponse{
		Message: "User registered successfully, please verify your email address",
	})
}

// Register godoc
// @Router /auth/verify-email [post]
// @Summary Register
// @Description Register
// @Tags auth
// @Accept  json
// @Produce  json
// @Param body  body user_service.VerifyEmailRequest true "User"
// @Success 200 {object} user_service.VerifyEmailResponse
// @Failure 400 {object} user_service.ErrorResponse
func (h *handler) VerifyEmail(ctx *gin.Context) {
	var (
		body *user_service.VerifyEmailRequest
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	key := fmt.Sprintf("otp-%s", body.Email)

	otp, err := h.redis.Get(ctx, key)
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, "Ooops, something went wrong", http.StatusInternalServerError)
		return
	}

	if otp != body.Otp {
		h.ReturnError(ctx, config.ErrorBadRequest, "Incorrect otp", http.StatusBadRequest)
		return
	}

	user, err := h.grpcClient.UserService().GetSingle(ctx, &user_service.UserSingleRequest{
		Email: body.Email,
	})
	if h.HandleDbError(ctx, err, "get single user") {
		return
	}

	user.Status = "active"

	_, err = h.grpcClient.UserService().Update(ctx, user)
	if h.HandleDbError(ctx, err, "update user") {
		return
	}

	// create session
	newSession := user_service.Session{
		UserId:       user.Id,
		IpAddress:    ctx.ClientIP(),
		ExpiresAt:    time.Now().Add(time.Hour * 999999).Format(time.RFC3339),
		UserAgent:    ctx.Request.UserAgent(),
		IsActive:     true,
		LastActiveAt: time.Now().Format(time.RFC3339),
		Platform:     body.Platform,
	}

	session, err := h.grpcClient.SessionService().Create(ctx, &newSession)
	if h.HandleDbError(ctx, err, "Error while creating new session") {
		return
	}

	// generate jwt token
	jwtFields := map[string]interface{}{
		"sub":        user.Id,
		"user_role":  user.UserRole,
		"user_type":  user.UserType,
		"platform":   body.Platform,
		"session_id": session.Id,
	}

	user.AccessToken, err = jwt.GenerateJWT(jwtFields, h.cfg.JWT)
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, "Oops, something went wrong!!!", http.StatusInternalServerError)
		return
	}

	ctx.JSON(200, gin.H{
		"user":    user,
		"session": session,
	})
}
