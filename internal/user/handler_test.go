package user

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestRouter(handler *UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/api/v1/users", handler.CreateUser)
	router.GET("/api/v1/users", handler.GetAll)
	router.GET("/api/v1/users/:id", handler.GetById)
	router.PUT("/api/v1/users/:id", handler.UpdateUser)
	router.DELETE("/api/v1/users/:id", handler.DeleteUser)
	return router
}

func TestHandler_CreateUser_InvalidJSON(t *testing.T) {
	router := setupTestRouter(&UserHandler{})

	w := performRequest(router, "POST", "/api/v1/users", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_CreateUser_EmptyBody(t *testing.T) {
	router := setupTestRouter(&UserHandler{})

	w := performRequest(router, "POST", "/api/v1/users", "{}")

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	if !strings.Contains(w.Body.String(), "error") {
		t.Error("expected error message in response body")
	}
}

func TestHandler_CreateUser_MissingFields(t *testing.T) {
	router := setupTestRouter(&UserHandler{})

	body := `{"name": "Alice"}`
	req, _ := http.NewRequest("POST", "/api/v1/users", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_GetById_InvalidID(t *testing.T) {
	router := setupTestRouter(&UserHandler{})

	w := performRequest(router, "GET", "/api/v1/users/abc", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	if !strings.Contains(w.Body.String(), "invalid") {
		t.Errorf("expected 'invalid' in error message, got: %s", w.Body.String())
	}
}

func TestHandler_GetById_FloatID(t *testing.T) {
	router := setupTestRouter(&UserHandler{})

	w := performRequest(router, "GET", "/api/v1/users/1.5", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_UpdateUser_InvalidID(t *testing.T) {
	router := setupTestRouter(&UserHandler{})

	body := `{"name": "test", "email": "test@example.com"}`
	req, _ := http.NewRequest("PUT", "/api/v1/users/abc", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	if !strings.Contains(w.Body.String(), "invalid") {
		t.Errorf("expected 'invalid' in error message, got: %s", w.Body.String())
	}
}

func TestHandler_UpdateUser_InvalidJSON(t *testing.T) {
	router := setupTestRouter(&UserHandler{})

	req, _ := http.NewRequest("PUT", "/api/v1/users/1", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_UpdateUser_MissingJSONBody(t *testing.T) {
	router := setupTestRouter(&UserHandler{})

	w := performRequest(router, "PUT", "/api/v1/users/1", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_UpdateUser_EmptyBody(t *testing.T) {
	router := setupTestRouter(&UserHandler{})

	w := performRequest(router, "PUT", "/api/v1/users/1", "{}")

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestHandler_DeleteUser_InvalidID(t *testing.T) {
	router := setupTestRouter(&UserHandler{})

	w := performRequest(router, "DELETE", "/api/v1/users/abc", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	if !strings.Contains(w.Body.String(), "invalid") {
		t.Errorf("expected 'invalid' in error message, got: %s", w.Body.String())
	}
}

func TestHandler_DeleteUser_FloatID(t *testing.T) {
	router := setupTestRouter(&UserHandler{})

	w := performRequest(router, "DELETE", "/api/v1/users/1.5", nil)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}
