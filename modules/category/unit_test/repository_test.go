package unittest

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq" // PostgreSQL driver

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/category/dto"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/category/repository"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
)

func TestCreateRepository(t *testing.T) {

	category1 := dto.ToDBCreateCategory{
		Name:        "Category 2",
		Code:        uuid.New().String(),
		CreatedByID: uuid.New().String(),
	}

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

	// create table
	err = categoryRepo.CreateTable("D:/ngoding/roketin/tugure-roketin/v2/blips-v2-backend/database/migrations/create_category_table.up.sql")
	t.Log(err)

	res, err := categoryRepo.CreateCategory(nil, category1)

	t.Log(err)
	t.Log(res)
}

func TestGetRepository(t *testing.T) {

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

	resAll, total, err := categoryRepo.GetIndexCategory(request.PageRequest{
		Page:    1,
		PerPage: 2,
	})

	// resAll, err := categoryRepo.GetAllCategory()

	t.Log(err)
	t.Log(total)
	t.Log(resAll)

	// resOne, err := categoryRepo.GetCategoryByID(resAll[0].ID)

	// t.Log(err)
	// t.Log(resOne)
}
