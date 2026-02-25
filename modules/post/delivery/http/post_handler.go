package http

import (
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rendyfutsuy/base-go/helpers/middleware"
	_reqContext "github.com/rendyfutsuy/base-go/helpers/middleware/request"
	"github.com/rendyfutsuy/base-go/helpers/request"
	"github.com/rendyfutsuy/base-go/helpers/response"
	"github.com/rendyfutsuy/base-go/models"
	"github.com/rendyfutsuy/base-go/modules/post"
	"github.com/rendyfutsuy/base-go/modules/post/dto"
)

type PostHandler struct {
	Usecase              post.Usecase
	validator            *validator.Validate
	mwPageRequest        _reqContext.IMiddlewarePageRequest
	middlewareAuth       middleware.IMiddlewareAuth
	middlewarePermission middleware.IMiddlewarePermission
}

func NewPostHandler(e *echo.Echo, uc post.Usecase, mwP _reqContext.IMiddlewarePageRequest, auth middleware.IMiddlewareAuth, middlewarePermission middleware.IMiddlewarePermission) {
	h := &PostHandler{Usecase: uc, validator: validator.New(), mwPageRequest: mwP, middlewareAuth: auth, middlewarePermission: middlewarePermission}

	// Public routes
	e.GET("/v1/post", h.GetIndex, h.mwPageRequest.PageRequestCtx)
	e.GET("/v1/post/:id", h.GetByID)

	// Protected routes
	r := e.Group("/v1/post")
	r.Use(h.middlewareAuth.AuthorizationCheck)

	permissionToCreate := []string{"post.create"}
	permissionToUpdate := []string{"post.update"}
	permissionToDelete := []string{"post.delete"}

	r.POST("", h.Create, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToCreate))
	r.PUT("/:id", h.Update, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToUpdate))
	r.DELETE("/:id", h.Delete, middleware.RequireActivatedUser, h.middlewarePermission.PermissionValidation(permissionToDelete))
}

// Create Post
// @Summary      Create post
// @Description  Create a post with optional thumbnail upload
// @Tags         Post
// @Accept       json
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        payload   body     dto.ReqCreatePost  true  "Post payload"
// @Param        thumbnail formData file                 false "Thumbnail file"
// @Success      200       {object} response.NonPaginationResponse{data=dto.RespPost}
// @Failure      400       {object} response.NonPaginationResponse
// @Router       /v1/post [post]
func (h *PostHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()

	req := new(dto.ReqCreatePost)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	thumbnailFile, _ := c.FormFile("thumbnail")
	var thumbnailData []byte
	var thumbnailName string
	if thumbnailFile != nil {
		src, err := thumbnailFile.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
		}
		defer src.Close()
		thumbnailData, err = io.ReadAll(src)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
		}
		thumbnailName = thumbnailFile.Filename
	}

	var userID string
	if user := c.Get("user"); user != nil {
		if u, ok := user.(models.User); ok {
			userID = u.ID.String()
		}
	}
	res, err := h.Usecase.Create(ctx, req, userID, thumbnailData, thumbnailName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(dto.ToRespPost(*res))
	return c.JSON(http.StatusOK, resp)
}

// Update Post
// @Summary      Update post
// @Description  Update a post with optional thumbnail upload
// @Tags         Post
// @Accept       json
// @Accept       multipart/form-data
// @Produce      json
// @Security     BearerAuth
// @Param        id        path     string               true  "Post ID"
// @Param        payload   body     dto.ReqUpdatePost  true  "Post payload"
// @Param        thumbnail formData file                 false "Thumbnail file"
// @Success      200       {object} response.NonPaginationResponse{data=dto.RespPost}
// @Failure      400       {object} response.NonPaginationResponse
// @Router       /v1/post/{id} [put]
func (h *PostHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")

	req := new(dto.ReqUpdatePost)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	thumbnailFile, _ := c.FormFile("thumbnail")
	var thumbnailData []byte
	var thumbnailName string
	if thumbnailFile != nil {
		src, err := thumbnailFile.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
		}
		defer src.Close()
		thumbnailData, err = io.ReadAll(src)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
		}
		thumbnailName = thumbnailFile.Filename
	}

	var userID string
	if user := c.Get("user"); user != nil {
		if u, ok := user.(models.User); ok {
			userID = u.ID.String()
		}
	}
	res, err := h.Usecase.Update(ctx, id, req, userID, thumbnailData, thumbnailName)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(dto.ToRespPost(*res))
	return c.JSON(http.StatusOK, resp)
}

// Delete Post
// @Summary      Delete post
// @Description  Delete a post by ID
// @Tags         Post
// @Produce      json
// @Security     BearerAuth
// @Param        id   path string true "Post ID"
// @Success      200  {object} response.NonPaginationResponse
// @Failure      400  {object} response.NonPaginationResponse
// @Router       /v1/post/{id} [delete]
func (h *PostHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	if err := h.Usecase.Delete(ctx, id, ""); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(struct {
		Message string `json:"message"`
	}{Message: "Successfully deleted Post"})
	return c.JSON(http.StatusOK, resp)
}

// Get Posts
// @Summary      Get paginated list of posts
// @Description  Retrieve a paginated list of posts with optional filtering
// @Tags         Post
// @Accept       json
// @Produce      json
// @Param        page        query   int                       false "Page number"     default(1)
// @Param        per_page    query   int                       false "Items per page"  default(10)
// @Param        sort_by     query   string                    false "Sort column"
// @Param        sort_order  query   string                    false "Sort order (asc/desc)"
// @Param        search      query   string                    false "Search query"
// @Param        filter      query   dto.ReqPostIndexFilter  false "Filter options"
// @Success      200         {object} response.PaginationResponse{data=[]dto.RespPostIndex}
// @Failure      400         {object} response.NonPaginationResponse
// @Router       /v1/post [get]
func (h *PostHandler) GetIndex(c echo.Context) error {
	ctx := c.Request().Context()
	pageRequest := c.Get("page_request").(*request.PageRequest)

	filter := new(dto.ReqPostIndexFilter)
	if err := c.Bind(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	if err := c.Validate(filter); err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}

	res, total, err := h.Usecase.GetIndex(ctx, *pageRequest, *filter)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	respPosts := make([]dto.RespPostIndex, 0, len(res))
	for _, v := range res {
		respPosts = append(respPosts, dto.ToRespPostIndex(v))
	}
	respPag := response.PaginationResponse{}
	respPag, err = respPag.SetResponse(respPosts, total, pageRequest.PerPage, pageRequest.Page)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	return c.JSON(http.StatusOK, respPag)
}

// Get Post By ID
// @Summary      Get post by ID
// @Description  Retrieve post detail and its parameters
// @Tags         Post
// @Produce      json
// @Param        id   path string true "Post ID"
// @Success      200  {object} response.NonPaginationResponse{data=dto.RespPost}
// @Failure      400  {object} response.NonPaginationResponse
// @Router       /v1/post/{id} [get]
func (h *PostHandler) GetByID(c echo.Context) error {
	ctx := c.Request().Context()
	id := c.Param("id")
	res, err := h.Usecase.GetByID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	level, lang, topics, err := h.Usecase.GetParameterReferences(ctx, id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.SetErrorResponse(http.StatusBadRequest, err.Error()))
	}
	out := dto.ToRespPost(*res)
	out.Level = level
	out.Lang = lang
	out.Topics = topics
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(out)
	return c.JSON(http.StatusOK, resp)
}
