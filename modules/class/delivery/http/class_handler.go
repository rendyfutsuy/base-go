package http

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	_reqContext "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/middleware/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/response"
	class "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/class"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/class/dto"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ResponseError struct {
	Message string `json:"message"`
}

type ClassHandler struct {
	ClassUseCase  class.Usecase
	validator     *validator.Validate
	mwPageRequest _reqContext.IMiddlewarePageRequest
}

func NewClassHandler(e *echo.Echo, us class.Usecase, mwP _reqContext.IMiddlewarePageRequest) {
	handler := &ClassHandler{
		ClassUseCase:  us,
		validator:     validator.New(),
		mwPageRequest: mwP,
	}

	r := e.Group("v1/class")

	r.POST("", handler.CreateClass)
	r.GET("", handler.GetIndexClass, handler.mwPageRequest.PageRequestCtx)
	r.GET("/all", handler.GetAllClass)
	r.GET("/:id", handler.GetClassByID)
	r.PUT("/:id", handler.UpdateClass)
	r.DELETE("/:id", handler.DeleteClass)
}

func (handler *ClassHandler) CreateClass(c echo.Context) error {

	authId := "authId" // get from middleware

	req := new(dto.ReqCreateClass)
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

	res, err := handler.ClassUseCase.CreateClass(c, req, authId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespClass(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *ClassHandler) GetIndexClass(c echo.Context) error {
	pageRequest := c.Get("page_request").(*request.PageRequest)

	res, total, err := handler.ClassUseCase.GetIndexClass(*pageRequest)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	respClass := []dto.RespClass{}

	for _, v := range res {
		respClass = append(respClass, dto.ToRespClass(v))
	}

	respPag := response.PaginationResponse{}
	respPag, err = respPag.SetResponse(respClass, total, pageRequest.PerPage, pageRequest.Page)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, respPag)
}

func (handler *ClassHandler) GetAllClass(c echo.Context) error {

	res, err := handler.ClassUseCase.GetAllClass()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	respClass := []dto.RespClass{}

	for _, v := range res {
		respClass = append(respClass, dto.ToRespClass(v))
	}

	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(respClass)

	return c.JSON(http.StatusOK, resp)
}

func (handler *ClassHandler) GetClassByID(c echo.Context) error {

	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid ID format"})
	}

	res, err := handler.ClassUseCase.GetClassByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespClass(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *ClassHandler) UpdateClass(c echo.Context) error {

	authId := "authId" // get from middleware
	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid ID format"})
	}

	req := new(dto.ReqUpdateClass)
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

	res, err := handler.ClassUseCase.UpdateClass(id, req, authId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespClass(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *ClassHandler) DeleteClass(c echo.Context) error {

	authId := "authId" // get from middleware
	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid ID format"})
	}

	res, err := handler.ClassUseCase.SoftDeleteClass(id, authId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespClass(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}
