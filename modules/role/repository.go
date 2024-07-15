package role

import (
	// "database/sql"

	models "git.roketin.com/tugure/blips/backend/v2/blips-v2-backend/models"
	"github.com/google/uuid"
)

// Repository represent the role's repository contract
type Repository interface {
	CreateRole(role models.Role) (id uuid.UUID, err error)
}
