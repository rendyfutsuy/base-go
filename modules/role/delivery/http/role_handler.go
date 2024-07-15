package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	// "github.com/google/uuid"
	"github.com/labstack/echo/v4"

	// "github.com/sirupsen/logrus"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/role"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/role/dto"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// RoleHandler  represent the http handler for role
type RoleHandler struct {
	RoleUseCase    role.Usecase
	validator      *validator.Validate
}

// NewRoleHandler will initialize the role/ resources endpoint
func NewRoleHandler(e *echo.Echo, us role.Usecase) {
	handler := &RoleHandler{
		RoleUseCase:    us,
	}

	r := e.Group("role")
	// r.Use(middlewareAuth.AuthorizationCheck)

	r.POST("",
		handler.CreateRole,
	)
	
}

func (handler *RoleHandler) CreateRole(c echo.Context) error {


	req := new(dto.ReqCreateRole)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	if err := handler.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	_, err := handler.RoleUseCase.CreateRole(c, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "success")
}
