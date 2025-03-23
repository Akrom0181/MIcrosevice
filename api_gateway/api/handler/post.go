package handler

import (
	"net/http"
	"strconv"
	"user_api_gateway/config"
	"user_api_gateway/genproto/post_service"
	"user_api_gateway/genproto/user_service"
	"user_api_gateway/pkg/etc"
	"user_api_gateway/pkg/gemini"

	"github.com/gin-gonic/gin"
)

// CreatePost godoc
// @Router /post [post]
// @Summary Create a new post
// @Description Create a new post
// @Security BearerAuth
// @Tags post
// @Accept  json
// @Produce  json
// @Param post body post_service.Post true "Post object"
// @Success 201 {object} post_service.Post
// @Failure 400 {object} user_service.ErrorResponse
func (h *handler) CreatePost(ctx *gin.Context) {
	var (
		body *post_service.Post
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	body.OwnerId = ctx.GetHeader("sub")

	if body.OwnerId == "" {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid owner id", 400)
		return
	}

	defaultTags, err := h.grpcClient.PostAttachmentService().GetDefaultTags(ctx, &post_service.GetDefaultTagsRequest{})
	if h.HandleDbError(ctx, err, "Error getting default tags") {
		return
	}

	tags, err := gemini.AskFromGemini(body.Content, defaultTags.Tags)
	if h.HandleDbError(ctx, err, "Error getting tags from gemini") {
		return
	}

	contentHashtags, err := etc.GetTags(ctx, body.Content)
	if h.HandleDbError(ctx, err, "Error getting tags") {
		return
	}

	body.Tags = map[string]*post_service.StringList{
		"gemini_tags":      {Values: tags},
		"owner_id":         {Values: []string{body.OwnerId}},
		"content_hashtags": {Values: contentHashtags},
	}

	post, err := h.grpcClient.PostService().Create(ctx, body)
	if h.HandleDbError(ctx, err, "Error creating post") {
		return
	}

	attachmentList, err := h.grpcClient.PostAttachmentService().MultipleUpsert(ctx, &post_service.AttachmentMultipleInsertRequest{
		PostId:      post.Id,
		Attachments: body.Attachments,
	})
	if err != nil {
		h.HandleDbError(ctx, err, "Error creating post attachments")
		return
	}
	post.Attachments = attachmentList.Items
	if h.HandleDbError(ctx, err, "Error creating post") {
		return
	}

	ctx.JSON(201, post)
}

// GetPost godoc
// @Router /post/{id} [get]
// @Summary Get a post by ID
// @Description Get a post by ID
// @Security BearerAuth
// @Tags post
// @Accept  json
// @Produce  json
// @Param id path string true "Post ID"
// @Success 200 {object} post_service.Post
// @Failure 400 {object} user_service.ErrorResponse
func (h *handler) GetPost(ctx *gin.Context) {
	var (
		req = &post_service.PostSingleRequest{}
	)

	req.Id = ctx.Param("id")

	post, err := h.grpcClient.PostService().GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting post") {
		return
	}

	postAttachments, err := h.grpcClient.PostAttachmentService().GetList(ctx, &post_service.GetListAttachmentRequest{
		Page:   1,
		Limit:  100,
		Search: post.Id})
	if h.HandleDbError(ctx, err, "Error getting post attachments") {
		return
	}

	post.Attachments = postAttachments.Items

	user, err := h.grpcClient.UserService().GetSingle(ctx, &user_service.UserSingleRequest{Id: post.OwnerId})
	if h.HandleDbError(ctx, err, "Error getting post owner") {
		return
	}
	post.OwnerId = user.Id

	ctx.JSON(200, gin.H{
		"post": post,
		"user": user,
	})
}

// GetPosts godoc
// @Router /post/list [get]
// @Summary Get a list of posts
// @Description Get a list of posts
// @Security BearerAuth
// @Tags post
// @Accept  json
// @Produce  json
// @Param page query number true "page"
// @Param limit query number true "limit"
// @Param owner_id query string false "owner_id"
// @Success 200 {object} post_service.PostList
// @Failure 400 {object} user_service.ErrorResponse
func (h *handler) GetPosts(ctx *gin.Context) {
	var (
		req = &post_service.GetListPostRequest{}
	)

	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	req.Search = ctx.Query("owner_id")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid page", 400)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid limit", 400)
		return
	}

	req.Page = uint64(page)
	req.Limit = uint64(limit)

	posts, err := h.grpcClient.PostService().GetList(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting posts") {
		return
	}

	ctx.JSON(200, posts)
}

// UpdatePost godoc
// @Router /post [put]
// @Summary Update a post
// @Description Update a post
// @Security BearerAuth
// @Tags post
// @Accept  json
// @Produce  json
// @Param post body post_service.Post true "Post object"
// @Success 200 {object} post_service.Post
// @Failure 400 {object} user_service.ErrorResponse
func (h *handler) UpdatePost(ctx *gin.Context) {
	var (
		body *post_service.Post
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	if body.OwnerId != ctx.GetHeader("sub") {
		h.ReturnError(ctx, config.ErrorForbidden, "You have no access to the post", http.StatusForbidden)
		return
	}

	defaultTags, err := h.grpcClient.PostAttachmentService().GetDefaultTags(ctx, &post_service.GetDefaultTagsRequest{})
	if h.HandleDbError(ctx, err, "Error getting default tags") {
		return
	}

	tags, err := gemini.AskFromGemini(body.Content, defaultTags.Tags)
	if h.HandleDbError(ctx, err, "Error getting tags from gemini") {
		return
	}

	contentHashtags, err := etc.GetTags(ctx, body.Content)
	if h.HandleDbError(ctx, err, "Error getting tags") {
		return
	}

	body.Tags = map[string]*post_service.StringList{
		"ownerId":    {Values: []string{body.OwnerId}},
		"geminiTags": {Values: tags},
		"getTagBody": {Values: contentHashtags},
	}

	post, err := h.grpcClient.PostService().Update(ctx, body)
	if h.HandleDbError(ctx, err, "Error updating post") {
		return
	}

	if h.HandleDbError(ctx, err, "Error upserting post attachments") {
		return
	}

	ctx.JSON(200, post)
}

// DeletePost godoc
// @Router /post/{id} [delete]
// @Summary Delete a post
// @Description Delete a post
// @Security BearerAuth
// @Tags post
// @Accept  json
// @Produce  json
// @Param id path string true "Post ID"
// @Success 200 {object} user_service.SuccessResponse
// @Failure 400 {object} user_service.ErrorResponse
func (h *handler) DeletePost(ctx *gin.Context) {
	var (
		req = &post_service.PostSingleRequest{}
	)

	req.Id = ctx.Param("id")

	post, err := h.grpcClient.PostService().GetSingle(ctx, req)
	if h.HandleDbError(ctx, err, "Error getting post") {
		return
	}

	if post.OwnerId != ctx.GetHeader("sub") {
		h.ReturnError(ctx, config.ErrorForbidden, "You have no access to the post", http.StatusForbidden)
		return
	}

	_, err = h.grpcClient.PostService().Delete(ctx, req)
	if h.HandleDbError(ctx, err, "Error deleting post") {
		return
	}

	_, err = h.grpcClient.PostAttachmentService().Delete(ctx, &post_service.AttachmentSingleRequest{Id: req.Id})
	if h.HandleDbError(ctx, err, "Error deleting post attachments") {
		return
	}

	ctx.JSON(200, user_service.SuccessResponse{
		Message: "Post deleted successfully",
	})
}
