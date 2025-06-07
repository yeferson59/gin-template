package models

import (
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestUserTableName(t *testing.T) {
	user := User{}
	expected := "users"
	if user.TableName() != expected {
		t.Errorf("TableName() = %s; want %s", user.TableName(), expected)
	}
}

func TestUserFields(t *testing.T) {
	now := time.Now()
	user := User{
		ID:        1,
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "hashedpassword",
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: gorm.DeletedAt{},
	}

	if user.ID != 1 {
		t.Errorf("ID = %d; want 1", user.ID)
	}
	if user.Username != "testuser" {
		t.Errorf("Username = %s; want testuser", user.Username)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Email = %s; want test@example.com", user.Email)
	}
	if user.Password != "hashedpassword" {
		t.Errorf("Password = %s; want hashedpassword", user.Password)
	}
	if !user.CreatedAt.Equal(now) {
		t.Errorf("CreatedAt = %v; want %v", user.CreatedAt, now)
	}
	if !user.UpdatedAt.Equal(now) {
		t.Errorf("UpdatedAt = %v; want %v", user.UpdatedAt, now)
	}
}
