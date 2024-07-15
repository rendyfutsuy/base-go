package unittest

import (
	"database/sql"
	"fmt"
	"testing"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/conveyance/dto"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/conveyance/repository"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/conveyance/usecase"
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

	conveyanceRepo := repository.NewConveyanceRepository(connDB)

	conveyanceUseCase := usecase.NewConveyanceUsecase(conveyanceRepo, 10)

	insert := dto.ReqCreateConveyance{
		Name: "Conveyance 2",
	}

	c := echo.New().NewContext(nil, nil)

	res, err := conveyanceUseCase.CreateConveyance(c, &insert, "authId")

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

	conveyanceRepo := repository.NewConveyanceRepository(connDB)

	conveyanceUseCase := usecase.NewConveyanceUsecase(conveyanceRepo, 10)

	resAll, err := conveyanceUseCase.GetAllConveyance()

	t.Log(err)
	t.Log(len(resAll))

	res, err := conveyanceUseCase.GetConveyanceByID(resAll[len(resAll)-1].ID.String())

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

	conveyanceRepo := repository.NewConveyanceRepository(connDB)

	conveyanceUseCase := usecase.NewConveyanceUsecase(conveyanceRepo, 10)

	resAll, err := conveyanceUseCase.GetAllConveyance()

	t.Log(err)

	conveyanceUp := dto.ReqUpdateConveyance{
		Name: "Conveyance 1 Updated",
	}

	res, err := conveyanceUseCase.UpdateConveyance(resAll[0].ID.String(), &conveyanceUp, uuid.New().String())

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

	conveyanceRepo := repository.NewConveyanceRepository(connDB)

	conveyanceUseCase := usecase.NewConveyanceUsecase(conveyanceRepo, 10)

	t.Log(err)

	res, err := conveyanceUseCase.SoftDeleteConveyance("a4b98e1b-39e1-11ef-856e-00ff7af1e5ed", uuid.New().String())

	t.Log(err)

	t.Log(res)
}
