package unittest

import (
	"database/sql"
	"fmt"
	"testing"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/account/dto"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/account/repository"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/account/usecase"
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

	accountRepo := repository.NewAccountRepository(connDB)

	accountUseCase := usecase.NewAccountUsecase(accountRepo, 10)

	account := dto.ReqCreateAccount{
		Name: "Account 2",
	}

	c := echo.New().NewContext(nil, nil)

	res, err := accountUseCase.CreateAccount(c, &account, uuid.New())

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

	accountRepo := repository.NewAccountRepository(connDB)

	accountUseCase := usecase.NewAccountUsecase(accountRepo, 10)

	resAll, err := accountUseCase.GetAllAccount()

	t.Log(err)
	t.Log(len(resAll))

	res, err := accountUseCase.GetAccountByID(resAll[len(resAll)-1].ID.String())

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

	accountRepo := repository.NewAccountRepository(connDB)

	accountUseCase := usecase.NewAccountUsecase(accountRepo, 10)

	resAll, err := accountUseCase.GetAllAccount()

	t.Log(err)

	accountUp := dto.ReqUpdateAccount{
		Name: "Account 1 Updated",
	}

	res, err := accountUseCase.UpdateAccount(resAll[0].ID.String(), &accountUp, uuid.New())

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

	accountRepo := repository.NewAccountRepository(connDB)

	accountUseCase := usecase.NewAccountUsecase(accountRepo, 10)

	t.Log(err)

	res, err := accountUseCase.SoftDeleteAccount("a4b98e1b-39e1-11ef-856e-00ff7af1e5ed", uuid.New())

	t.Log(err)

	t.Log(res)
}
