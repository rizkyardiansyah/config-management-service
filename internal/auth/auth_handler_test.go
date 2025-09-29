package auth

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockAuthService struct {
	loginCalled bool
}

func (m *mockAuthService) Login(username, password string) (string, string, error) {
	m.loginCalled = true
	if password == "correct" {
		return "access123", "refresh123", nil
	}
	return "", "", http.ErrNoCookie // just return error
}

func (m *mockAuthService) Logout(userID string) error {
	return nil
}

func TestAuthHandler_Login_Success(t *testing.T) {
	svc := &mockAuthService{}
	h := NewAuthHandler(svc)

	body := bytes.NewBufferString(`{"username":"john","password":"correct"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Result().StatusCode)
	}
}

func TestAuthHandler_Login_Invalid(t *testing.T) {
	svc := &mockAuthService{}
	h := NewAuthHandler(svc)

	body := bytes.NewBufferString(`{"username":"john","password":"wrong"}`)
	req := httptest.NewRequest(http.MethodPost, "/login", body)
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Result().StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Result().StatusCode)
	}
}
