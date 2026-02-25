package token_storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rendyfutsuy/base-go/constants"
	"github.com/rendyfutsuy/base-go/models"
	"gorm.io/gorm"
)

type LocalStorage struct {
	DB *gorm.DB
}

func NewLocalStorage(db *gorm.DB) *LocalStorage {
	return &LocalStorage{
		DB: db,
	}
}

func (s *LocalStorage) SaveSession(ctx context.Context, user models.User, accessToken, refreshToken, accessJTI, refreshJTI string, refreshTTL time.Duration) error {
	now := time.Now().UTC()
	token := models.JWTToken{
		UserId:           user.ID,
		AccessToken:      accessToken,
		AccessJTI:        accessJTI,
		RefreshToken:     refreshToken,
		RefreshJTI:       refreshJTI,
		RefreshExpiresAt: now.Add(refreshTTL),
		IsUsed:           false,
		CreatedAt:        now,
		UpdatedAt:        &now,
	}

	return s.DB.WithContext(ctx).Create(&token).Error
}

func (s *LocalStorage) GetRefreshTokenMetadata(ctx context.Context, refreshJTI string) (RefreshTokenMeta, error) {
	var token models.JWTToken
	err := s.DB.WithContext(ctx).
		Where("refresh_jti = ?", refreshJTI).
		First(&token).Error

	if err != nil {
		return RefreshTokenMeta{}, err
	}

	return RefreshTokenMeta{
		UserID:    token.UserId,
		ExpiresAt: token.RefreshExpiresAt,
		Used:      token.IsUsed,
		AccessJTI: token.AccessJTI,
	}, nil
}

func (s *LocalStorage) MarkRefreshTokenUsed(ctx context.Context, refreshJTI string) error {
	return s.DB.WithContext(ctx).
		Model(&models.JWTToken{}).
		Where("refresh_jti = ?", refreshJTI).
		Update("is_used", true).Error
}

func (s *LocalStorage) DestroySession(ctx context.Context, accessToken string) error {
	// Delete by access token
	return s.DB.WithContext(ctx).
		Where("access_token = ?", accessToken).
		Delete(&models.JWTToken{}).Error
}

func (s *LocalStorage) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	return s.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.JWTToken{}).Error
}

func (s *LocalStorage) ValidateAccessToken(ctx context.Context, accessToken string) (models.User, error) {
	var token models.JWTToken
	// Check if session exists for this access token
	err := s.DB.WithContext(ctx).
		Where("access_token = ?", accessToken).
		First(&token).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("invalid session")
		}
		return models.User{}, err
	}

	// Retrieve User with Role
	var user models.User
	err = s.DB.WithContext(ctx).
		Table("users usr").
		Select(`usr.id, usr.full_name, usr.email, usr.username, usr.is_active, usr.gender, usr.role_id, usr.is_first_time_login,
			(SELECT f.file_path FROM files_to_module ftm
				JOIN files f ON f.id = ftm.file_id AND f.deleted_at IS NULL
				WHERE ftm.module_type = ? AND ftm.module_id = usr.id AND ftm.type = ?
				ORDER BY ftm.created_at DESC
				LIMIT 1
			) AS avatar_url,
			roles.name as role_name, usr.verified_at`,
			constants.ModuleTypeUser, constants.FileTypeAvatar,
		).
		Joins("LEFT JOIN roles ON roles.id = usr.role_id AND roles.deleted_at IS NULL").
		Where("usr.id = ? AND usr.deleted_at IS NULL", token.UserId).
		Scan(&user).Error

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
