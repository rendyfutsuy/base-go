package usecase

import (

	models "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/role/dto"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

)

func (u *roleUsecase) CreateRole(c echo.Context, req *dto.ReqCreateRole) (id uuid.UUID, err error) {

	// userLoggedIn := c.Get("user").(models.User)
	// createdBy := userLoggedIn.UserCode
	// createdAt := time.Now()

	role := models.Role{
		Name:      req.Name,
		Deletable: true,
	}

	id, err = u.roleRepo.CreateRole(role)
	if err != nil {
		return id, err
	}


	return id, err
}
