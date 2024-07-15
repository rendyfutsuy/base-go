package unittest

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq" // PostgreSQL driver

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/conveyance/dto"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/conveyance/repository"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
)

func TestCreateRepository(t *testing.T) {

	insert := dto.ToDBCreateConveyance{
		Name:        "data 1",
		Code:        uuid.New().String(),
		Type:        "type 1",
		CreatedByID: "authId",
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

	repo := repository.NewConveyanceRepository(connDB)

	// create table

	res, err := repo.CreateConveyance(insert)

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

	conveyanceRepo := repository.NewConveyanceRepository(connDB)

	resAll, total, err := conveyanceRepo.GetIndexConveyance(request.PageRequest{
		Page:    1,
		PerPage: 2,
	})

	// resAll, err := conveyanceRepo.GetAllConveyance()

	t.Log(err)
	t.Log(total)
	t.Log(resAll)

	resOne, err := conveyanceRepo.GetConveyanceByID(resAll[0].ID)

	t.Log(err)
	t.Log(resOne)
}
