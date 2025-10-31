package storage

import (
	"errors"
	"sync"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidToken      = errors.New("invalid token")
)

// User represents a user in the system
type User struct {
	ID           string
	Username     string
	Email        string
	PasswordHash string
}

// Storage provides an interface for user storage
type Storage interface {
	CreateUser(user *User) error
	GetUserByUsername(username string) (*User, error)
	GetUserByID(id string) (*User, error)
	StoreToken(userID, token string) error
	ValidateToken(token string) (string, error)
	RevokeToken(token string) error
}

// InMemoryStorage is an in-memory implementation of Storage
type InMemoryStorage struct {
	mu      sync.RWMutex
	users   map[string]*User  // username -> User
	userIDs map[string]*User  // userID -> User
	tokens  map[string]string // token -> userID
}

// NewInMemoryStorage creates a new in-memory storage
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		users:   make(map[string]*User),
		userIDs: make(map[string]*User),
		tokens:  make(map[string]string),
	}
}

// CreateUser stores a new user
func (s *InMemoryStorage) CreateUser(user *User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[user.Username]; exists {
		return ErrUserAlreadyExists
	}

	s.users[user.Username] = user
	s.userIDs[user.ID] = user
	return nil
}

// GetUserByUsername retrieves a user by username
func (s *InMemoryStorage) GetUserByUsername(username string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[username]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// GetUserByID retrieves a user by ID
func (s *InMemoryStorage) GetUserByID(id string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.userIDs[id]
	if !exists {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// StoreToken stores a token for a user
func (s *InMemoryStorage) StoreToken(userID, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.tokens[token] = userID
	return nil
}

// ValidateToken checks if a token is valid and returns the user ID
func (s *InMemoryStorage) ValidateToken(token string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userID, exists := s.tokens[token]
	if !exists {
		return "", ErrInvalidToken
	}
	return userID, nil
}

// RevokeToken removes a token
func (s *InMemoryStorage) RevokeToken(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.tokens, token)
	return nil
}
