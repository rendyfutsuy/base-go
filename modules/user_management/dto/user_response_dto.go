package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/models"
)

type ReqConfirmationUserPassword struct {
	Password string `json:"password" validate:"required"`
}

type RespUser struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"name"`
}

type RespUserIndex struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	RoleName  string    `json:"role_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Deletable bool      `json:"deletable"`
}

type ReqUserIndexFilter struct {
	RoleIds   []uuid.UUID `query:"role_ids" json:"role_ids"`
	RoleName  string      `query:"role_name" json:"role_name"`
	SortBy    string      `query:"sort_by" json:"sort_by"`
	SortOrder string      `query:"sort_order" json:"sort_order"`
}

type RespPermissionGroupUserDetail struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"name"`
	RoleName string    `json:"role_name"`
	RoleId   uuid.UUID `json:"role_id"`
}

type RespUserDetail struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Nik       string    `json:"nik"`
	RoleId    uuid.UUID `json:"role_id"`
	RoleName  string    `json:"role_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Deletable bool      `json:"deletable"`
}

// to get role info for compact use
func ToRespUser(userDb models.User) RespUser {

	return RespUser{
		ID:       userDb.ID,
		FullName: userDb.FullName,
	}

}

// for get role for index use
func ToRespUserIndex(userDb models.User) RespUserIndex {

	return RespUserIndex{
		ID:        userDb.ID,
		FullName:  userDb.FullName,
		Username:  userDb.Username,
		Email:     userDb.Email,
		RoleName:  userDb.RoleName,
		Deletable: userDb.Deletable,
		CreatedAt: userDb.CreatedAt,
		UpdatedAt: userDb.UpdatedAt,
	}

}

// to get role info with references
func ToRespUserDetail(userDb models.User) RespUserDetail {

	return RespUserDetail{
		ID:        userDb.ID,
		FullName:  userDb.FullName,
		Username:  userDb.Username,
		Email:     userDb.Email,
		Nik:       userDb.Nik,
		RoleId:    userDb.RoleId,
		RoleName:  userDb.RoleName,
		Deletable: userDb.Deletable,
		CreatedAt: userDb.CreatedAt,
		UpdatedAt: userDb.UpdatedAt,
	}
}
