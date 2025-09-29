package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"sass.com/configsvc/internal/models"
	"sass.com/configsvc/internal/secrets"
)

func setupTestRouterWithMW(secret []byte) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	secs := &secrets.Secrets{JWTsecret: secret}
	r.Use(AuthMiddleware(secs))

	// Dummy handler
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	return r
}

func generateTestToken(secret []byte, expired bool) string {
	claims := jwt.MapClaims{
		"sub":  uuid.New(),
		"role": string(models.RoleUser),
		"iat":  time.Now().Unix(),
	}
	if expired {
		claims["exp"] = time.Now().Add(-time.Minute).Unix() // already expired
	} else {
		claims["exp"] = time.Now().Add(time.Minute).Unix()
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, _ := t.SignedString(secret)
	return token
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	secret := []byte("testsecret")
	r := setupTestRouterWithMW(secret)

	token := generateTestToken(secret, false)
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	secret := []byte("testsecret")
	r := setupTestRouterWithMW(secret)

	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	secret := []byte("testsecret")
	r := setupTestRouterWithMW(secret)

	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "NotBearer token")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	secret := []byte("testsecret")
	r := setupTestRouterWithMW(secret)

	token := generateTestToken(secret, true)
	req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for expired token, got %d", w.Code)
	}
}
