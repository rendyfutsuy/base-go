package repository

import (
	"database/sql"
	"os"

	role "github.com/rendyfutsuy/base-go/modules/role_management"
	"github.com/rendyfutsuy/base-go/utils"
)

type roleRepository struct {
	Conn *sql.DB
}

func NewRoleManagementRepository(Conn *sql.DB) role.Repository {
	return &roleRepository{Conn}
}

func (repo *roleRepository) CreateTable(sqlFilePath string) (err error) {

	sqlCommands, err := os.ReadFile(sqlFilePath)
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	_, err = repo.Conn.Exec(string(sqlCommands))
	if err != nil {
		utils.Logger.Error(err.Error())
		return err
	}

	return err
}
