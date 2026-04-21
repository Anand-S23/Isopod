package models
 
import (
	"time"

	"github.com/google/uuid"
)
	
type GHUserDto struct {
	ID          int64  `json:"id"`
	Login       string `json:"login"`
	Email       string `json:"email"`
	AvatarURL   string `json:"avatar_url"`
}

type User struct {
    ID          string    `json:"id"`
    GitHubID    int64     `json:"github_id"`
    Username    string    `json:"username"`
    Email       string    `json:"email"`
    AvatarURL   string    `json:"avatar_url"`
    AccessToken string    `json:"-"`
    CreatedAt   time.Time `json:"created_at"`
    LastLoginAt time.Time `json:"last_login"`
}

func NewUser(dto GHUserDto, encryptionKey []byte) User {
	return User{
		ID:          uuid.New().String(),
		GitHubID:    dto.ID,
		Username:    dto.Login,
		Email:       dto.Email,
		AvatarURL:   dto.AvatarURL,
		AccessToken: crypt.Encrypt(dto.AccessToken, encryptionKey),
		CreatedAt:   time.Now(),
		LastLoginAt: time.Now(),
	}
}
