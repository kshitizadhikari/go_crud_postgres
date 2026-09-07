package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func performRequest(r http.Handler, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		data, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(data)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestCreateUserRequestValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		body       CreateUserRequest
		wantStatus int
	}{
		{
			name:       "valid request",
			body:       CreateUserRequest{Name: "Alice", Email: "alice@example.com", Password: "secret123"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing name",
			body:       CreateUserRequest{Email: "bob@example.com", Password: "secret123"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing email",
			body:       CreateUserRequest{Name: "Bob", Password: "secret123"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid email format",
			body:       CreateUserRequest{Name: "Bob", Email: "not-an-email", Password: "secret123"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "password too short",
			body:       CreateUserRequest{Name: "Bob", Email: "bob@example.com", Password: "123"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty body",
			body:       CreateUserRequest{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.POST("/api/v1/users", func(c *gin.Context) {
				var req CreateUserRequest
				if err := c.ShouldBind(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"ok": true})
			})

			w := performRequest(router, "POST", "/api/v1/users", tt.body)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d, body: %s", w.Code, tt.wantStatus, w.Body.String())
			}
		})
	}
}

func TestUpdateUserRequestValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		body       UpdateUserRequest
		wantStatus int
	}{
		{
			name:       "valid request",
			body:       UpdateUserRequest{Name: "Alice Updated", Email: "alice-new@example.com"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "missing name",
			body:       UpdateUserRequest{Email: "alice@example.com"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "missing email",
			body:       UpdateUserRequest{Name: "Alice"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid email",
			body:       UpdateUserRequest{Name: "Alice", Email: "bad"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "empty body",
			body:       UpdateUserRequest{},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.PUT("/api/v1/users/:id", func(c *gin.Context) {
				var req UpdateUserRequest
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}
				c.JSON(http.StatusOK, gin.H{"ok": true})
			})

			w := performRequest(router, "PUT", "/api/v1/users/1", tt.body)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d, body: %s", w.Code, tt.wantStatus, w.Body.String())
			}
		})
	}
}
