package role

import (
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/role/dto"
	"github.com/labstack/echo/v4"
	"github.com/google/uuid"
)

// Usecase represent the role's usecases
type Usecase interface {
	CreateRole(c echo.Context, req *dto.ReqCreateRole) (id uuid.UUID, err error)
}
