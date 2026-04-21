package models
 
import (
	"time"

	"github.com/google/uuid"
)
	
type NewUserDto struct {
	GitHubID    int64  `json:"github_id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	AvatarURL   string `json:"avatar_url"`
	AccessToken string `json:"access_token"`
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

func NewUser(dto NewUserDto, encryptionKey []byte) User {
	return User{
		ID:          uuid.New().String(),
		GitHubID:    dto.GitHubID,
		Username:    dto.Username,
		Email:       dto.Email,
		AvatarURL:   dto.AvatarURL,
		AccessToken: crypt.Encrypt(dto.AccessToken, encryptionKey),
		CreatedAt:   time.Now(),
		LastLoginAt: time.Now(),
	}
}
