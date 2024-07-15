package http

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	_reqContext "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/middleware/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/response"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helpers/middleware"
	role "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/role"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/role/dto"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ResponseError struct {
	Message        string `json:"message"`
	middlewareAuth middleware.IMiddlewareAuth
}

type RoleHandler struct {
	RoleUseCase    role.Usecase
	validator      *validator.Validate
	mwPageRequest  _reqContext.IMiddlewarePageRequest
	middlewareAuth middleware.IMiddlewareAuth
}

func NewRoleHandler(e *echo.Echo, us role.Usecase, mwP _reqContext.IMiddlewarePageRequest, middlewareAuth middleware.IMiddlewareAuth) {
	handler := &RoleHandler{
		RoleUseCase:    us,
		validator:      validator.New(),
		mwPageRequest:  mwP,
		middlewareAuth: middlewareAuth,
	}

	r := e.Group("v1/role")

	r.Use(handler.middlewareAuth.AuthorizationCheck)
	r.GET("", handler.GetIndexRole, handler.mwPageRequest.PageRequestCtx)
	r.POST("", handler.CreateRole)
	r.GET("/all", handler.GetAllRole)
	r.GET("/:id", handler.GetRoleByID)
	r.PUT("/:id", handler.UpdateRole)
	r.DELETE("/:id", handler.DeleteRole)
}

func (handler *RoleHandler) CreateRole(c echo.Context) error {
	authId := c.Get("authId") // get from middleware

	req := new(dto.ReqCreateRole)
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

	res, err := handler.RoleUseCase.CreateRole(c, req, authId.(uuid.UUID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespRole(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *RoleHandler) GetIndexRole(c echo.Context) error {
	pageRequest := c.Get("page_request").(*request.PageRequest)

	res, total, err := handler.RoleUseCase.GetIndexRole(*pageRequest)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	respRole := []dto.RespRole{}

	for _, v := range res {
		respRole = append(respRole, dto.ToRespRole(v))
	}

	respPag := response.PaginationResponse{}
	respPag, err = respPag.SetResponse(respRole, total, pageRequest.PerPage, pageRequest.Page)

	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, respPag)
}

func (handler *RoleHandler) GetAllRole(c echo.Context) error {

	res, err := handler.RoleUseCase.GetAllRole()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	respRole := []dto.RespRole{}

	for _, v := range res {
		respRole = append(respRole, dto.ToRespRole(v))
	}

	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(respRole)

	return c.JSON(http.StatusOK, resp)
}

func (handler *RoleHandler) GetRoleByID(c echo.Context) error {

	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid ID format"})
	}

	res, err := handler.RoleUseCase.GetRoleByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespRole(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *RoleHandler) UpdateRole(c echo.Context) error {

	authId := c.Get("authId") // get from middleware
	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid ID format"})
	}

	req := new(dto.ReqUpdateRole)
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

	res, err := handler.RoleUseCase.UpdateRole(id, req, authId.(uuid.UUID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespRole(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}

func (handler *RoleHandler) DeleteRole(c echo.Context) error {

	authId := c.Get("authId") // get from middleware
	id := c.Param("id")

	// validate id
	err := uuid.Validate(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid ID format"})
	}

	res, err := handler.RoleUseCase.SoftDeleteRole(id, authId.(uuid.UUID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	resResp := dto.ToRespRole(*res)
	resp := response.NonPaginationResponse{}
	resp, _ = resp.SetResponse(resResp)

	return c.JSON(http.StatusOK, resp)
}
