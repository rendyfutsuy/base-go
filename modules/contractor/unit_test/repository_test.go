package unittest

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq" // PostgreSQL driver

	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/helper/request"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/contractor/dto"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/contractor/repository"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/utils"
)

func TestCreateRepository(t *testing.T) {

	insert := dto.ToDBCreateContractor{
		Name:        "data 1",
		Code:        uuid.New().String(),
		Address:     "address 1",
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

	repo := repository.NewContractorRepository(connDB)

	// create table

	res, err := repo.CreateContractor(insert)

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

	contractorRepo := repository.NewContractorRepository(connDB)

	resAll, total, err := contractorRepo.GetIndexContractor(request.PageRequest{
		Page:    1,
		PerPage: 2,
	})

	// resAll, err := contractorRepo.GetAllContractor()

	t.Log(err)
	t.Log(total)
	t.Log(resAll)

	resOne, err := contractorRepo.GetContractorByID(resAll[0].ID)

	t.Log(err)
	t.Log(resOne)
}
