package unittest

import (
	"database/sql"
	"fmt"
	"testing"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/contractor/dto"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/contractor/repository"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/contractor/usecase"
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

	contractorRepo := repository.NewContractorRepository(connDB)

	contractorUseCase := usecase.NewContractorUsecase(contractorRepo, 10)

	insert := dto.ReqCreateContractor{
		Name: "Contractor 2",
	}

	c := echo.New().NewContext(nil, nil)

	res, err := contractorUseCase.CreateContractor(c, &insert, "authId")

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

	contractorRepo := repository.NewContractorRepository(connDB)

	contractorUseCase := usecase.NewContractorUsecase(contractorRepo, 10)

	resAll, err := contractorUseCase.GetAllContractor()

	t.Log(err)
	t.Log(len(resAll))

	res, err := contractorUseCase.GetContractorByID(resAll[len(resAll)-1].ID.String())

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

	contractorRepo := repository.NewContractorRepository(connDB)

	contractorUseCase := usecase.NewContractorUsecase(contractorRepo, 10)

	resAll, err := contractorUseCase.GetAllContractor()

	t.Log(err)

	contractorUp := dto.ReqUpdateContractor{
		Name: "Contractor 1 Updated",
	}

	res, err := contractorUseCase.UpdateContractor(resAll[0].ID.String(), &contractorUp, uuid.New().String())

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

	contractorRepo := repository.NewContractorRepository(connDB)

	contractorUseCase := usecase.NewContractorUsecase(contractorRepo, 10)

	t.Log(err)

	res, err := contractorUseCase.SoftDeleteContractor("a4b98e1b-39e1-11ef-856e-00ff7af1e5ed", uuid.New().String())

	t.Log(err)

	t.Log(res)
}
