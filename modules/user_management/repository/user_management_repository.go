package repository

import (
	"database/sql"
	"os"

	user "github.com/rendyfutsuy/base-go/modules/user_management"
	"github.com/rendyfutsuy/base-go/utils"
)

type userRepository struct {
	Conn *sql.DB
}

func NewUserManagementRepository(Conn *sql.DB) user.Repository {
	return &userRepository{Conn}
}

func (repo *userRepository) CreateTable(sqlFilePath string) (err error) {

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
