package repository

import (
	"database/sql"

	models "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/modules/role"
	"github.com/google/uuid"
)

type roleRepository struct {
	Conn *sql.DB
}

func NewRoleRepository(Conn *sql.DB) role.Repository {
	return &roleRepository{Conn}
}


func (repo *roleRepository) CreateRole(role models.Role) (id uuid.UUID, err error) {

	err = repo.Conn.QueryRow(
		`INSERT INTO roles 
			(name, deletable, created_at, created_by) 
		VALUES 
			($1, $2, $3, $4) 
		RETURNING id`,
		role.Name,
		role.Deletable,
	).Scan(&id)

	if err != nil {
		return id, err
	}

	return id, err
}