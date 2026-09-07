package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTableName(t *testing.T) {
	u := User{}
	if u.TableName() != "users" {
		t.Errorf("TableName() = %q, want %q", u.TableName(), "users")
	}
}

func TestToResponse(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	user := &User{
		ID:        1,
		Name:      "Alice",
		Email:     "alice@example.com",
		Password:  "secret123",
		AvatarKey: "avatars/123_avatar.png",
		CreatedAt: now,
		UpdatedAt: now,
	}

	resp := user.ToResponse()

	if resp.ID != user.ID {
		t.Errorf("ID = %d, want %d", resp.ID, user.ID)
	}
	if resp.Name != user.Name {
		t.Errorf("Name = %q, want %q", resp.Name, user.Name)
	}
	if resp.Email != user.Email {
		t.Errorf("Email = %q, want %q", resp.Email, user.Email)
	}
	if resp.CreatedAt != user.CreatedAt {
		t.Errorf("CreatedAt = %v, want %v", resp.CreatedAt, user.CreatedAt)
	}
	if resp.UpdatedAt != user.UpdatedAt {
		t.Errorf("UpdatedAt = %v, want %v", resp.UpdatedAt, user.UpdatedAt)
	}
}

func TestToResponseDoesNotExposePassword(t *testing.T) {
	user := &User{
		ID:       1,
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "supersecret",
	}

	resp := user.ToResponse()

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal UserResponse: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if _, exists := raw["password"]; exists {
		t.Error("UserResponse JSON should not contain 'password' field")
	}
}

func TestUserJSONExcludesPassword(t *testing.T) {
	user := User{
		ID:       1,
		Name:     "Charlie",
		Email:    "charlie@example.com",
		Password: "hidden_password",
	}

	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("failed to marshal User: %v", err)
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if _, exists := raw["password"]; exists {
		t.Error("User JSON should not contain 'password' field")
	}

	if raw["name"] != "Charlie" {
		t.Errorf("name = %v, want Charlie", raw["name"])
	}
	if raw["email"] != "charlie@example.com" {
		t.Errorf("email = %v, want charlie@example.com", raw["email"])
	}
}

func TestToResponsePreservesZeroValues(t *testing.T) {
	user := &User{}
	resp := user.ToResponse()

	if resp.ID != 0 {
		t.Errorf("ID = %d, want 0", resp.ID)
	}
	if resp.Name != "" {
		t.Errorf("Name = %q, want empty", resp.Name)
	}
	if resp.Email != "" {
		t.Errorf("Email = %q, want empty", resp.Email)
	}
}
