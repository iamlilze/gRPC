package storage

import (
	"testing"
)

func TestCreateUser(t *testing.T) {
	store := NewInMemoryStorage()

	user := &User{
		ID:           "user-1",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}

	err := store.CreateUser(user)
	if err != nil {
		t.Errorf("CreateUser() error = %v", err)
	}

	// Try to create duplicate user
	err = store.CreateUser(user)
	if err != ErrUserAlreadyExists {
		t.Errorf("CreateUser() duplicate error = %v, want %v", err, ErrUserAlreadyExists)
	}
}

func TestGetUserByUsername(t *testing.T) {
	store := NewInMemoryStorage()

	user := &User{
		ID:           "user-1",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}

	store.CreateUser(user)

	// Get existing user
	retrieved, err := store.GetUserByUsername("testuser")
	if err != nil {
		t.Errorf("GetUserByUsername() error = %v", err)
	}
	if retrieved.ID != user.ID {
		t.Errorf("GetUserByUsername() ID = %v, want %v", retrieved.ID, user.ID)
	}

	// Get non-existent user
	_, err = store.GetUserByUsername("nonexistent")
	if err != ErrUserNotFound {
		t.Errorf("GetUserByUsername() error = %v, want %v", err, ErrUserNotFound)
	}
}

func TestGetUserByID(t *testing.T) {
	store := NewInMemoryStorage()

	user := &User{
		ID:           "user-1",
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}

	store.CreateUser(user)

	// Get existing user
	retrieved, err := store.GetUserByID("user-1")
	if err != nil {
		t.Errorf("GetUserByID() error = %v", err)
	}
	if retrieved.Username != user.Username {
		t.Errorf("GetUserByID() Username = %v, want %v", retrieved.Username, user.Username)
	}

	// Get non-existent user
	_, err = store.GetUserByID("nonexistent")
	if err != ErrUserNotFound {
		t.Errorf("GetUserByID() error = %v, want %v", err, ErrUserNotFound)
	}
}

func TestTokenOperations(t *testing.T) {
	store := NewInMemoryStorage()

	userID := "user-1"
	token := "test-token"

	// Store token
	err := store.StoreToken(userID, token)
	if err != nil {
		t.Errorf("StoreToken() error = %v", err)
	}

	// Validate token
	retrievedUserID, err := store.ValidateToken(token)
	if err != nil {
		t.Errorf("ValidateToken() error = %v", err)
	}
	if retrievedUserID != userID {
		t.Errorf("ValidateToken() userID = %v, want %v", retrievedUserID, userID)
	}

	// Validate invalid token
	_, err = store.ValidateToken("invalid-token")
	if err != ErrInvalidToken {
		t.Errorf("ValidateToken() error = %v, want %v", err, ErrInvalidToken)
	}

	// Revoke token
	err = store.RevokeToken(token)
	if err != nil {
		t.Errorf("RevokeToken() error = %v", err)
	}

	// Validate revoked token
	_, err = store.ValidateToken(token)
	if err != ErrInvalidToken {
		t.Errorf("ValidateToken() after revoke error = %v, want %v", err, ErrInvalidToken)
	}
}
