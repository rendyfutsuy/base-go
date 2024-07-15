package unittest

import (
	"database/sql"
	"fmt"
	"testing"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/role/dto"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/role/repository"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/role/usecase"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func TestCreateUsecase(t *testing.T) {

	utils.InitConfig("config.json")
	dbString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		utils.ConfigVars.String("database.blips.host"),
		utils.ConfigVars.Int("database.blips.port"),
		utils.ConfigVars.String("database.blips.user"),
		utils.ConfigVars.String("database.blips.password"),
		utils.ConfigVars.String("database.blips.db_name"),
		utils.ConfigVars.String("database.blips.sslmode"),
	)

	connDB, err := sql.Open("postgres", dbString)
	t.Log(err)

	roleRepo := repository.NewRoleRepository(connDB)

	roleUseCase := usecase.NewRoleUsecase(roleRepo, 10)

	role := dto.ReqCreateRole{
		Name: "Role 2",
	}

	c := echo.New().NewContext(nil, nil)

	res, err := roleUseCase.CreateRole(c, &role, uuid.New())

	t.Log(err)
	t.Log(res)
}

func TestGetUsecase(t *testing.T) {

	utils.InitConfig("config.json")
	dbString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		utils.ConfigVars.String("database.blips.host"),
		utils.ConfigVars.Int("database.blips.port"),
		utils.ConfigVars.String("database.blips.user"),
		utils.ConfigVars.String("database.blips.password"),
		utils.ConfigVars.String("database.blips.db_name"),
		utils.ConfigVars.String("database.blips.sslmode"),
	)

	connDB, err := sql.Open("postgres", dbString)
	t.Log(err)

	roleRepo := repository.NewRoleRepository(connDB)

	roleUseCase := usecase.NewRoleUsecase(roleRepo, 10)

	resAll, err := roleUseCase.GetAllRole()

	t.Log(err)
	t.Log(len(resAll))

	res, err := roleUseCase.GetRoleByID(resAll[len(resAll)-1].ID.String())

	t.Log(err)

	t.Log(res)
}

func TestUpdateUsecase(t *testing.T) {

	utils.InitConfig("config.json")
	dbString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		utils.ConfigVars.String("database.blips.host"),
		utils.ConfigVars.Int("database.blips.port"),
		utils.ConfigVars.String("database.blips.user"),
		utils.ConfigVars.String("database.blips.password"),
		utils.ConfigVars.String("database.blips.db_name"),
		utils.ConfigVars.String("database.blips.sslmode"),
	)

	connDB, err := sql.Open("postgres", dbString)
	t.Log(err)

	roleRepo := repository.NewRoleRepository(connDB)

	roleUseCase := usecase.NewRoleUsecase(roleRepo, 10)

	resAll, err := roleUseCase.GetAllRole()

	t.Log(err)

	roleUp := dto.ReqUpdateRole{
		Name: "Role 1 Updated",
	}

	res, err := roleUseCase.UpdateRole(resAll[0].ID.String(), &roleUp, uuid.New())

	t.Log(err)

	t.Log(res)
}

func TestDeleteUsecase(t *testing.T) {

	utils.InitConfig("config.json")
	dbString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		utils.ConfigVars.String("database.blips.host"),
		utils.ConfigVars.Int("database.blips.port"),
		utils.ConfigVars.String("database.blips.user"),
		utils.ConfigVars.String("database.blips.password"),
		utils.ConfigVars.String("database.blips.db_name"),
		utils.ConfigVars.String("database.blips.sslmode"),
	)

	connDB, err := sql.Open("postgres", dbString)
	t.Log(err)

	roleRepo := repository.NewRoleRepository(connDB)

	roleUseCase := usecase.NewRoleUsecase(roleRepo, 10)

	t.Log(err)

	res, err := roleUseCase.SoftDeleteRole("a4b98e1b-39e1-11ef-856e-00ff7af1e5ed", uuid.New())

	t.Log(err)

	t.Log(res)
}
