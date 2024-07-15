package http

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	_reqContext "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/middleware/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/response"
	category "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/category"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/category/dto"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ResponseError struct {
	Message string `json:"message"`
}

type CategoryHandler struct {
	CategoryUseCase category.Usecase
	validator       *validator.Validate
	mwPageRequest   _reqContext.IMiddlewarePageRequest
}

func NewCategoryHandler(e *echo.Echo, us category.Usecase, mwP _reqContext.IMiddlewarePageRequest) {
	handler := &CategoryHandler{
		CategoryUseCase: us,
		validator:       validator.New(),
		mwPageRequest:   mwP,
	}

	r := e.Group("v1/category")

	r.POST("", handler.CreateCategory)
	r.GET("", handler.GetIndexCategory, handler.mwPageRequest.PageRequestCtx)
	r.GET("/all", handler.GetAllCategory)
	r.GET("/:id", handler.GetCategoryByID)
	r.PUT("/:id", handler.UpdateCategory)
	r.DELETE("/:id", handler.DeleteCategory)
}

func (handler *CategoryHandler) CreateCategory(c echo.Context) error {

	authId := "authId" // get from middleware

	req := new(dto.ReqCreateCategory)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	if err := handler.validator.Struct(req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			var errorMessages []string
			for _, fe := range ve {
				// Construct a human-friendly error message
				fieldName := strings.ToLower(fe.Field())
				errorMessages = append(errorMessages, fmt.Sprintf("%s is required", fieldName))
			}
			// Join all error messages and return as a single string
			errorMessage := strings.Join(errorMessages, ", ")
			return c.JSON(http.StatusBadRequest, ResponseError{Message: errorMessage})
		}
		// Fallback for any other errors
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	res, err := handler.CategoryUseCase.CreateCategory(c, req, authId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespCategory(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *CategoryHandler) GetIndexCategory(c echo.Context) error {
	pageRequest := c.Get("page_request").(*request.PageRequest)

	res, total, err := handler.CategoryUseCase.GetIndexCategory(*pageRequest)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	respCategory := []dto.RespCategory{}

	for _, v := range res {
		respCategory = append(respCategory, dto.ToRespCategory(v))
	}

	respPag := response.PaginationResponse{}
	respPag, err = respPag.SetResponse(respCategory, total, pageRequest.PerPage, pageRequest.Page)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, respPag)
}

func (handler *CategoryHandler) GetAllCategory(c echo.Context) error {

	res, err := handler.CategoryUseCase.GetAllCategory()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	respCategory := []dto.RespCategory{}

	for _, v := range res {
		respCategory = append(respCategory, dto.ToRespCategory(v))
	}

	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(respCategory)

	return c.JSON(http.StatusOK, resp)
}

func (handler *CategoryHandler) GetCategoryByID(c echo.Context) error {

	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid ID format"})
	}

	res, err := handler.CategoryUseCase.GetCategoryByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespCategory(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *CategoryHandler) UpdateCategory(c echo.Context) error {

	authId := "authId" // get from middleware
	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid ID format"})
	}

	req := new(dto.ReqUpdateCategory)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	if err := handler.validator.Struct(req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			var errorMessages []string
			for _, fe := range ve {
				// Construct a human-friendly error message
				fieldName := strings.ToLower(fe.Field())
				errorMessages = append(errorMessages, fmt.Sprintf("%s is required", fieldName))
			}
			// Join all error messages and return as a single string
			errorMessage := strings.Join(errorMessages, ", ")
			return c.JSON(http.StatusBadRequest, ResponseError{Message: errorMessage})
		}
		// Fallback for any other errors
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	res, err := handler.CategoryUseCase.UpdateCategory(id, req, authId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespCategory(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *CategoryHandler) DeleteCategory(c echo.Context) error {

	authId := "authId" // get from middleware
	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid ID format"})
	}

	res, err := handler.CategoryUseCase.SoftDeleteCategory(id, authId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespCategory(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}
