package middleware

import (
	"backend/config"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const testKey = "test-secret-key"

func init() {
	config.Key = testKey
}

func makeToken(role string, subject uint, exp time.Time) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"subject": float64(subject),
		"role":    role,
		"exp":     exp.Unix(),
	})
	signed, _ := token.SignedString([]byte(testKey))
	return signed
}

func setupRouter(role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/protected", AuthMiddleware(role), func(c *gin.Context) {
		id, _ := c.Get("ID")
		c.JSON(http.StatusOK, gin.H{"id": id})
	})
	return r
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	r := setupRouter("user")
	token := makeToken("user", 42, time.Now().Add(time.Hour))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "key", Value: token})
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAuthMiddleware_MissingCookie(t *testing.T) {
	r := setupRouter("user")

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	r := setupRouter("user")

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "key", Value: "not.a.valid.token"})
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	r := setupRouter("user")
	token := makeToken("user", 1, time.Now().Add(-time.Hour))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "key", Value: token})
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_WrongRole(t *testing.T) {
	r := setupRouter("admin")
	token := makeToken("user", 1, time.Now().Add(time.Hour))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "key", Value: token})
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_SetsIDInContext(t *testing.T) {
	r := setupRouter("user")
	token := makeToken("user", 42, time.Now().Add(time.Hour))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "key", Value: token})
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if body == "" {
		t.Error("expected response body with ID")
	}
}

func TestAuthMiddleware_TrainerRole(t *testing.T) {
	r := setupRouter("trainer")
	token := makeToken("trainer", 5, time.Now().Add(time.Hour))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.AddCookie(&http.Cookie{Name: "key", Value: token})
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}
