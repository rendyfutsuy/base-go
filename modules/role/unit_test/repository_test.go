package unittest

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq" // PostgreSQL driver

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/role/dto"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/role/repository"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
)

func TestCreateRepository(t *testing.T) {

	role1 := dto.ToDBCreateRole{
		Name:        "Role 2",
		Code:        uuid.New().String(),
		CreatedByID: uuid.New(),
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

	roleRepo := repository.NewRoleRepository(connDB)

	// create table
	err = roleRepo.CreateTable("D:/ngoding/roketin/tugure-roketin/v2/blips-v2-backend/database/migrations/create_role_table.up.sql")
	t.Log(err)

	res, err := roleRepo.CreateRole(role1)

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

	roleRepo := repository.NewRoleRepository(connDB)

	resAll, total, err := roleRepo.GetIndexRole(request.PageRequest{
		Page:    1,
		PerPage: 2,
	})

	// resAll, err := roleRepo.GetAllRole()

	t.Log(err)
	t.Log(total)
	t.Log(resAll)

	// resOne, err := roleRepo.GetRoleByID(resAll[0].ID)

	// t.Log(err)
	// t.Log(resOne)
}
