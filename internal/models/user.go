package models

import (
	"time"

	"github.com/Anand-S23/isopod/pkg/crypt"
	"github.com/google/uuid"
)

type GHUserDto struct {
	ID          int64  `json:"id"`
	Login       string `json:"login"`
	Email       string `json:"email"`
	AvatarURL   string `json:"avatar_url"`
	AccessToken string `json:"access_token"`
}

type User struct {
	ID          string    `json:"id"         db:"id"`
	GitHubID    int64     `json:"github_id"  db:"github_id"`
	Username    string    `json:"username"   db:"username"`
	Email       string    `json:"email"      db:"email"`
	AvatarURL   string    `json:"avatar_url" db:"avatar_url"`
	AccessToken string    `json:"-"          db:"access_token"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	LastLoginAt time.Time `json:"last_login" db:"last_login_at"`
}

func NewUser(dto GHUserDto, encryptionKey []byte) (User, error) {
	encryptedToken, err := crypt.Encrypt(dto.AccessToken, encryptionKey)
	if err != nil {
		return User{}, err
	}

	return User{
		ID:          uuid.New().String(),
		GitHubID:    dto.ID,
		Username:    dto.Login,
		Email:       dto.Email,
		AvatarURL:   dto.AvatarURL,
		AccessToken: encryptedToken,
		CreatedAt:   time.Now(),
		LastLoginAt: time.Now(),
	}, nil
}
