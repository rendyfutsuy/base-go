package unittest

import (
	"database/sql"
	"fmt"
	"testing"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/carriage/dto"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/carriage/repository"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/carriage/usecase"
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

	carriageRepo := repository.NewCarriageRepository(connDB)

	carriageUseCase := usecase.NewCarriageUsecase(carriageRepo, 10)

	insert := dto.ReqCreateCarriage{
		Name: "Carriage 2",
	}

	c := echo.New().NewContext(nil, nil)

	res, err := carriageUseCase.CreateCarriage(c, &insert, "authId")

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

	carriageRepo := repository.NewCarriageRepository(connDB)

	carriageUseCase := usecase.NewCarriageUsecase(carriageRepo, 10)

	resAll, err := carriageUseCase.GetAllCarriage()

	t.Log(err)
	t.Log(len(resAll))

	res, err := carriageUseCase.GetCarriageByID(resAll[len(resAll)-1].ID.String())

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

	carriageRepo := repository.NewCarriageRepository(connDB)

	carriageUseCase := usecase.NewCarriageUsecase(carriageRepo, 10)

	resAll, err := carriageUseCase.GetAllCarriage()

	t.Log(err)

	carriageUp := dto.ReqUpdateCarriage{
		Name: "Carriage 1 Updated",
	}

	res, err := carriageUseCase.UpdateCarriage(resAll[0].ID.String(), &carriageUp, uuid.New().String())

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

	carriageRepo := repository.NewCarriageRepository(connDB)

	carriageUseCase := usecase.NewCarriageUsecase(carriageRepo, 10)

	t.Log(err)

	res, err := carriageUseCase.SoftDeleteCarriage("a4b98e1b-39e1-11ef-856e-00ff7af1e5ed", uuid.New().String())

	t.Log(err)

	t.Log(res)
}
