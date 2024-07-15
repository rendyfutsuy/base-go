package http

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	_reqContext "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/middleware/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/response"
	conveyance "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/conveyance"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/conveyance/dto"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ResponseError struct {
	Message string `json:"message"`
}

type ConveyanceHandler struct {
	ConveyanceUseCase conveyance.Usecase
	validator         *validator.Validate
	mwPageRequest     _reqContext.IMiddlewarePageRequest
}

func NewConveyanceHandler(e *echo.Echo, us conveyance.Usecase, mwP _reqContext.IMiddlewarePageRequest) {
	handler := &ConveyanceHandler{
		ConveyanceUseCase: us,
		validator:         validator.New(),
		mwPageRequest:     mwP,
	}

	r := e.Group("v1/conveyance")

	r.POST("", handler.CreateConveyance)
	r.GET("", handler.GetIndexConveyance, handler.mwPageRequest.PageRequestCtx)
	r.GET("/all", handler.GetAllConveyance)
	r.GET("/:id", handler.GetConveyanceByID)
	r.PUT("/:id", handler.UpdateConveyance)
	r.DELETE("/:id", handler.DeleteConveyance)
}

func (handler *ConveyanceHandler) CreateConveyance(c echo.Context) error {

	authId := "authId" // get from middleware

	req := new(dto.ReqCreateConveyance)
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

	res, err := handler.ConveyanceUseCase.CreateConveyance(c, req, authId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespConveyance(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *ConveyanceHandler) GetIndexConveyance(c echo.Context) error {
	pageRequest := c.Get("page_request").(*request.PageRequest)

	res, total, err := handler.ConveyanceUseCase.GetIndexConveyance(*pageRequest)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	respConveyance := []dto.RespConveyance{}

	for _, v := range res {
		respConveyance = append(respConveyance, dto.ToRespConveyance(v))
	}

	respPag := response.PaginationResponse{}
	respPag, err = respPag.SetResponse(respConveyance, total, pageRequest.PerPage, pageRequest.Page)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, respPag)
}

func (handler *ConveyanceHandler) GetAllConveyance(c echo.Context) error {

	res, err := handler.ConveyanceUseCase.GetAllConveyance()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	respConveyance := []dto.RespConveyance{}

	for _, v := range res {
		respConveyance = append(respConveyance, dto.ToRespConveyance(v))
	}

	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(respConveyance)

	return c.JSON(http.StatusOK, resp)
}

func (handler *ConveyanceHandler) GetConveyanceByID(c echo.Context) error {

	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid ID format"})
	}

	res, err := handler.ConveyanceUseCase.GetConveyanceByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespConveyance(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *ConveyanceHandler) UpdateConveyance(c echo.Context) error {

	authId := "authId" // get from middleware
	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid ID format"})
	}

	req := new(dto.ReqUpdateConveyance)
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

	res, err := handler.ConveyanceUseCase.UpdateConveyance(id, req, authId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespConveyance(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *ConveyanceHandler) DeleteConveyance(c echo.Context) error {

	authId := "authId" // get from middleware
	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid ID format"})
	}

	res, err := handler.ConveyanceUseCase.SoftDeleteConveyance(id, authId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespConveyance(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}
