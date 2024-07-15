package unittest

import (
	"database/sql"
	"fmt"
	"testing"

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/category/dto"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/category/repository"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/category/usecase"
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

	categoryRepo := repository.NewCategoryRepository(connDB)

	categoryUseCase := usecase.NewCategoryUsecase(categoryRepo, 10)

	category := dto.ReqCreateCategory{
		Name: "Category 2",
	}

	c := echo.New().NewContext(nil, nil)

	res, err := categoryUseCase.CreateCategory(c, &category, uuid.New().String())

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

	categoryRepo := repository.NewCategoryRepository(connDB)

	categoryUseCase := usecase.NewCategoryUsecase(categoryRepo, 10)

	resAll, err := categoryUseCase.GetAllCategory()

	t.Log(err)
	t.Log(len(resAll))

	res, err := categoryUseCase.GetCategoryByID(resAll[len(resAll)-1].ID.String())

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

	categoryRepo := repository.NewCategoryRepository(connDB)

	categoryUseCase := usecase.NewCategoryUsecase(categoryRepo, 10)

	resAll, err := categoryUseCase.GetAllCategory()

	t.Log(err)

	categoryUp := dto.ReqUpdateCategory{
		Name: "Category 1 Updated",
	}

	res, err := categoryUseCase.UpdateCategory(resAll[0].ID.String(), &categoryUp, uuid.New().String())

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

	categoryRepo := repository.NewCategoryRepository(connDB)

	categoryUseCase := usecase.NewCategoryUsecase(categoryRepo, 10)

	t.Log(err)

	res, err := categoryUseCase.SoftDeleteCategory("a4b98e1b-39e1-11ef-856e-00ff7af1e5ed", uuid.New().String())

	t.Log(err)

	t.Log(res)
}
