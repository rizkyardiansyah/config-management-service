package auth

import (
	"errors"

	"sass.com/configsvc/internal/config"
	"sass.com/configsvc/internal/secrets"
)

type AuthService interface {
	Login(username, password string) (string, string, error)
}

type AuthServiceImpl struct {
	userRepo UserRepository
	cfg      *config.Config
	secrets  *secrets.Secrets
}

func NewAuthService(userRepo UserRepository, cfg *config.Config, secrets *secrets.Secrets) AuthService {
	return &AuthServiceImpl{
		userRepo: userRepo,
		cfg:      cfg,
		secrets:  secrets,
	}
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

	access, err := createAccessToken(*u, s.cfg, s.secrets)
	if err != nil {
		return "", "", err
	}

	// To simplify, we don't store refresh token and its expiration time to Sqlite
	refresh, err := createRefreshToken()
	if err != nil {
		return "", "", err
	}
	return access, refresh, nil
}
