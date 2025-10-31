package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	pb "github.com/iamlilze/gRPC/api/auth"
	"github.com/iamlilze/gRPC/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

// AuthService implements the gRPC AuthService
type AuthService struct {
	pb.UnimplementedAuthServiceServer
	storage storage.Storage
}

// NewAuthService creates a new auth service
func NewAuthService(store storage.Storage) *AuthService {
	return &AuthService{
		storage: store,
	}
}

// Register creates a new user account
func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Validate input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return &pb.RegisterResponse{
			Success: false,
			Message: "username, email, and password are required",
		}, nil
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return &pb.RegisterResponse{
			Success: false,
			Message: "failed to hash password",
		}, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &storage.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	// Store user
	if err := s.storage.CreateUser(user); err != nil {
		if err == storage.ErrUserAlreadyExists {
			return &pb.RegisterResponse{
				Success: false,
				Message: "user already exists",
			}, nil
		}
		return &pb.RegisterResponse{
			Success: false,
			Message: "failed to create user",
		}, fmt.Errorf("failed to create user: %w", err)
	}

	return &pb.RegisterResponse{
		UserId:  user.ID,
		Success: true,
		Message: "user registered successfully",
	}, nil
}

// Login authenticates a user and returns a token
func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// Validate input
	if req.Username == "" || req.Password == "" {
		return &pb.LoginResponse{
			Success: false,
			Message: "username and password are required",
		}, nil
	}

	// Get user
	user, err := s.storage.GetUserByUsername(req.Username)
	if err != nil {
		if err == storage.ErrUserNotFound {
			return &pb.LoginResponse{
				Success: false,
				Message: "invalid username or password",
			}, nil
		}
		return &pb.LoginResponse{
			Success: false,
			Message: "failed to authenticate",
		}, fmt.Errorf("failed to get user: %w", err)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return &pb.LoginResponse{
			Success: false,
			Message: "invalid username or password",
		}, nil
	}

	// Generate token
	token := uuid.New().String()
	if err := s.storage.StoreToken(user.ID, token); err != nil {
		return &pb.LoginResponse{
			Success: false,
			Message: "failed to generate token",
		}, fmt.Errorf("failed to store token: %w", err)
	}

	return &pb.LoginResponse{
		Token:   token,
		UserId:  user.ID,
		Success: true,
		Message: "login successful",
	}, nil
}

// ValidateToken checks if a token is valid
func (s *AuthService) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	if req.Token == "" {
		return &pb.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	// Validate token
	userID, err := s.storage.ValidateToken(req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	// Get user details
	user, err := s.storage.GetUserByID(userID)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:    true,
		UserId:   user.ID,
		Username: user.Username,
	}, nil
}

// Logout invalidates a user's token
func (s *AuthService) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	if req.Token == "" {
		return &pb.LogoutResponse{
			Success: false,
			Message: "token is required",
		}, nil
	}

	// Revoke token
	if err := s.storage.RevokeToken(req.Token); err != nil {
		return &pb.LogoutResponse{
			Success: false,
			Message: "failed to logout",
		}, fmt.Errorf("failed to revoke token: %w", err)
	}

	return &pb.LogoutResponse{
		Success: true,
		Message: "logout successful",
	}, nil
}
