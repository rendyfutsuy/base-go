package http

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	_reqContext "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/middleware/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/response"
	contractor "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/contractor"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/contractor/dto"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ResponseError struct {
	Message string `json:"message"`
}

type ContractorHandler struct {
	ContractorUseCase contractor.Usecase
	validator         *validator.Validate
	mwPageRequest     _reqContext.IMiddlewarePageRequest
}

func NewContractorHandler(e *echo.Echo, us contractor.Usecase, mwP _reqContext.IMiddlewarePageRequest) {
	handler := &ContractorHandler{
		ContractorUseCase: us,
		validator:         validator.New(),
		mwPageRequest:     mwP,
	}

	r := e.Group("v1/contractor")

	r.POST("", handler.CreateContractor)
	r.GET("", handler.GetIndexContractor, handler.mwPageRequest.PageRequestCtx)
	r.GET("/all", handler.GetAllContractor)
	r.GET("/:id", handler.GetContractorByID)
	r.PUT("/:id", handler.UpdateContractor)
	r.DELETE("/:id", handler.DeleteContractor)
}

func (handler *ContractorHandler) CreateContractor(c echo.Context) error {

	authId := "authId" // get from middleware

	req := new(dto.ReqCreateContractor)
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

	res, err := handler.ContractorUseCase.CreateContractor(c, req, authId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespContractor(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *ContractorHandler) GetIndexContractor(c echo.Context) error {
	pageRequest := c.Get("page_request").(*request.PageRequest)

	res, total, err := handler.ContractorUseCase.GetIndexContractor(*pageRequest)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	respContractor := []dto.RespContractor{}

	for _, v := range res {
		respContractor = append(respContractor, dto.ToRespContractor(v))
	}

	respPag := response.PaginationResponse{}
	respPag, err = respPag.SetResponse(respContractor, total, pageRequest.PerPage, pageRequest.Page)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, respPag)
}

func (handler *ContractorHandler) GetAllContractor(c echo.Context) error {

	res, err := handler.ContractorUseCase.GetAllContractor()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	respContractor := []dto.RespContractor{}

	for _, v := range res {
		respContractor = append(respContractor, dto.ToRespContractor(v))
	}

	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(respContractor)

	return c.JSON(http.StatusOK, resp)
}

func (handler *ContractorHandler) GetContractorByID(c echo.Context) error {

	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid ID format"})
	}

	res, err := handler.ContractorUseCase.GetContractorByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespContractor(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *ContractorHandler) UpdateContractor(c echo.Context) error {

	authId := "authId" // get from middleware
	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid ID format"})
	}

	req := new(dto.ReqUpdateContractor)
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

	res, err := handler.ContractorUseCase.UpdateContractor(id, req, authId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespContractor(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *ContractorHandler) DeleteContractor(c echo.Context) error {

	authId := "authId" // get from middleware
	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid ID format"})
	}

	res, err := handler.ContractorUseCase.SoftDeleteContractor(id, authId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespContractor(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}
