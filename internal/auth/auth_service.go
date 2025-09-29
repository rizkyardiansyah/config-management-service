package auth

import (
	"errors"

	"sass.com/configsvc/internal/models"
)

type AuthService interface {
	Login(username, password string) (string, string, error)
}

type AuthServiceImpl struct {
	userRepo UserRepository
}

type UserRepository interface {
	FindByUsername(username string) (*models.User, error)
	VerifyPassword(u *models.User, password string) bool
}

func NewAuthService(userRepo UserRepository) AuthService {
	return &AuthServiceImpl{userRepo: userRepo}
}

// Login verifies user and returns access + refresh tokens
func (s *AuthServiceImpl) Login(username, password string) (string, string, error) {
	u, err := s.userRepo.FindByUsername(username)
	if err != nil || u == nil {
		return "", "", errors.New("invalid credentials")
	}
	if !s.userRepo.VerifyPassword(u, password) {
		return "", "", errors.New("invalid credentials")
	}

	access, err := createAccessToken(*u)
	if err != nil {
		return "", "", err
	}
	refresh, err := createRefreshToken()
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}
