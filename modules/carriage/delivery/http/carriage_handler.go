package http

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	_reqContext "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/middleware/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/response"
	carriage "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/carriage"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/carriage/dto"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ResponseError struct {
	Message string `json:"message"`
}

type CarriageHandler struct {
	CarriageUseCase carriage.Usecase
	validator       *validator.Validate
	mwPageRequest   _reqContext.IMiddlewarePageRequest
}

func NewCarriageHandler(e *echo.Echo, us carriage.Usecase, mwP _reqContext.IMiddlewarePageRequest) {
	handler := &CarriageHandler{
		CarriageUseCase: us,
		validator:       validator.New(),
		mwPageRequest:   mwP,
	}

	r := e.Group("v1/carriage")

	r.POST("", handler.CreateCarriage)
	r.GET("", handler.GetIndexCarriage, handler.mwPageRequest.PageRequestCtx)
	r.GET("/all", handler.GetAllCarriage)
	r.GET("/:id", handler.GetCarriageByID)
	r.PUT("/:id", handler.UpdateCarriage)
	r.DELETE("/:id", handler.DeleteCarriage)
}

func (handler *CarriageHandler) CreateCarriage(c echo.Context) error {

	authId := "authId" // get from middleware

	req := new(dto.ReqCreateCarriage)
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

	res, err := handler.CarriageUseCase.CreateCarriage(c, req, authId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespCarriage(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *CarriageHandler) GetIndexCarriage(c echo.Context) error {
	pageRequest := c.Get("page_request").(*request.PageRequest)

	res, total, err := handler.CarriageUseCase.GetIndexCarriage(*pageRequest)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	respCarriage := []dto.RespCarriage{}

	for _, v := range res {
		respCarriage = append(respCarriage, dto.ToRespCarriage(v))
	}

	respPag := response.PaginationResponse{}
	respPag, err = respPag.SetResponse(respCarriage, total, pageRequest.PerPage, pageRequest.Page)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, respPag)
}

func (handler *CarriageHandler) GetAllCarriage(c echo.Context) error {

	res, err := handler.CarriageUseCase.GetAllCarriage()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	respCarriage := []dto.RespCarriage{}

	for _, v := range res {
		respCarriage = append(respCarriage, dto.ToRespCarriage(v))
	}

	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(respCarriage)

	return c.JSON(http.StatusOK, resp)
}

func (handler *CarriageHandler) GetCarriageByID(c echo.Context) error {

	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid ID format"})
	}

	res, err := handler.CarriageUseCase.GetCarriageByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespCarriage(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *CarriageHandler) UpdateCarriage(c echo.Context) error {

	authId := "authId" // get from middleware
	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid ID format"})
	}

	req := new(dto.ReqUpdateCarriage)
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

	res, err := handler.CarriageUseCase.UpdateCarriage(id, req, authId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespCarriage(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *CarriageHandler) DeleteCarriage(c echo.Context) error {

	authId := "authId" // get from middleware
	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid ID format"})
	}

	res, err := handler.CarriageUseCase.SoftDeleteCarriage(id, authId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespCarriage(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}
