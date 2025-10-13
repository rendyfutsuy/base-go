package http

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	_reqContext "github.com/rendyfutsuy/base-go.git/helper/middleware/request"
	"github.com/rendyfutsuy/base-go.git/helper/request"
	"github.com/rendyfutsuy/base-go.git/helper/response"
	"github.com/rendyfutsuy/base-go.git/helpers/middleware"
	account "github.com/rendyfutsuy/base-go.git/modules/account"
	"github.com/rendyfutsuy/base-go.git/modules/account/dto"
)

type ResponseError struct {
	Message        string `json:"message"`
	middlewareAuth middleware.IMiddlewareAuth
}

type AccountHandler struct {
	AccountUseCase account.Usecase
	validator      *validator.Validate
	mwPageRequest  _reqContext.IMiddlewarePageRequest
	middlewareAuth middleware.IMiddlewareAuth
}

func NewAccountHandler(e *echo.Echo, us account.Usecase, mwP _reqContext.IMiddlewarePageRequest, middlewareAuth middleware.IMiddlewareAuth) {
	handler := &AccountHandler{
		AccountUseCase: us,
		validator:      validator.New(),
		mwPageRequest:  mwP,
		middlewareAuth: middlewareAuth,
	}

	r := e.Group("v1/account")

	r.Use(handler.middlewareAuth.AuthorizationCheck)
	r.GET("", handler.GetIndexAccount, handler.mwPageRequest.PageRequestCtx)
	r.POST("", handler.CreateAccount)
	r.GET("/all", handler.GetAllAccount)
	r.GET("/:id", handler.GetAccountByID)
	r.PUT("/:id", handler.UpdateAccount)
	r.DELETE("/:id", handler.DeleteAccount)
}

func (handler *AccountHandler) CreateAccount(c echo.Context) error {
	authId := c.Get("authId") // get from middleware

	req := new(dto.ReqCreateAccount)
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

	res, err := handler.AccountUseCase.CreateAccount(c, req, authId.(uuid.UUID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespAccount(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *AccountHandler) GetIndexAccount(c echo.Context) error {
	pageRequest := c.Get("page_request").(*request.PageRequest)

	res, total, err := handler.AccountUseCase.GetIndexAccount(*pageRequest)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	respAccount := []dto.RespAccount{}

	for _, v := range res {
		respAccount = append(respAccount, dto.ToRespAccount(v))
	}

	respPag := response.PaginationResponse{}
	respPag, err = respPag.SetResponse(respAccount, total, pageRequest.PerPage, pageRequest.Page)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, respPag)
}

func (handler *AccountHandler) GetAllAccount(c echo.Context) error {

	res, err := handler.AccountUseCase.GetAllAccount()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	respAccount := []dto.RespAccount{}

	for _, v := range res {
		respAccount = append(respAccount, dto.ToRespAccount(v))
	}

	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(respAccount)

	return c.JSON(http.StatusOK, resp)
}

func (handler *AccountHandler) GetAccountByID(c echo.Context) error {

	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid ID format"})
	}

	res, err := handler.AccountUseCase.GetAccountByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespAccount(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *AccountHandler) UpdateAccount(c echo.Context) error {

	authId := c.Get("authId") // get from middleware
	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid ID format"})
	}

	req := new(dto.ReqUpdateAccount)
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

	res, err := handler.AccountUseCase.UpdateAccount(id, req, authId.(uuid.UUID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespAccount(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *AccountHandler) DeleteAccount(c echo.Context) error {

	authId := c.Get("authId") // get from middleware
	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid ID format"})
	}

	res, err := handler.AccountUseCase.SoftDeleteAccount(id, authId.(uuid.UUID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespAccount(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}
