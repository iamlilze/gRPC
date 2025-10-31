package service

import (
	"context"
	"testing"

	pb "github.com/iamlilze/gRPC/api/auth"
	"github.com/iamlilze/gRPC/internal/storage"
)

func TestRegister(t *testing.T) {
	store := storage.NewInMemoryStorage()
	service := NewAuthService(store)
	ctx := context.Background()

	tests := []struct {
		name        string
		req         *pb.RegisterRequest
		wantSuccess bool
		wantError   bool
	}{
		{
			name: "valid registration",
			req: &pb.RegisterRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantSuccess: true,
			wantError:   false,
		},
		{
			name: "duplicate username",
			req: &pb.RegisterRequest{
				Username: "testuser",
				Email:    "test2@example.com",
				Password: "password123",
			},
			wantSuccess: false,
			wantError:   false,
		},
		{
			name: "missing username",
			req: &pb.RegisterRequest{
				Username: "",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantSuccess: false,
			wantError:   false,
		},
		{
			name: "missing password",
			req: &pb.RegisterRequest{
				Username: "testuser2",
				Email:    "test@example.com",
				Password: "",
			},
			wantSuccess: false,
			wantError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.Register(ctx, tt.req)
			if (err != nil) != tt.wantError {
				t.Errorf("Register() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if resp.Success != tt.wantSuccess {
				t.Errorf("Register() success = %v, want %v, message = %v", resp.Success, tt.wantSuccess, resp.Message)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	store := storage.NewInMemoryStorage()
	service := NewAuthService(store)
	ctx := context.Background()

	// Register a user first
	_, err := service.Register(ctx, &pb.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	tests := []struct {
		name        string
		req         *pb.LoginRequest
		wantSuccess bool
		wantError   bool
	}{
		{
			name: "valid login",
			req: &pb.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			wantSuccess: true,
			wantError:   false,
		},
		{
			name: "invalid password",
			req: &pb.LoginRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			wantSuccess: false,
			wantError:   false,
		},
		{
			name: "non-existent user",
			req: &pb.LoginRequest{
				Username: "nonexistent",
				Password: "password123",
			},
			wantSuccess: false,
			wantError:   false,
		},
		{
			name: "missing username",
			req: &pb.LoginRequest{
				Username: "",
				Password: "password123",
			},
			wantSuccess: false,
			wantError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.Login(ctx, tt.req)
			if (err != nil) != tt.wantError {
				t.Errorf("Login() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if resp.Success != tt.wantSuccess {
				t.Errorf("Login() success = %v, want %v, message = %v", resp.Success, tt.wantSuccess, resp.Message)
			}
			if tt.wantSuccess && resp.Token == "" {
				t.Error("Login() expected non-empty token")
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	store := storage.NewInMemoryStorage()
	service := NewAuthService(store)
	ctx := context.Background()

	// Register and login a user
	regResp, _ := service.Register(ctx, &pb.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	})

	loginResp, _ := service.Login(ctx, &pb.LoginRequest{
		Username: "testuser",
		Password: "password123",
	})

	tests := []struct {
		name      string
		req       *pb.ValidateTokenRequest
		wantValid bool
	}{
		{
			name: "valid token",
			req: &pb.ValidateTokenRequest{
				Token: loginResp.Token,
			},
			wantValid: true,
		},
		{
			name: "invalid token",
			req: &pb.ValidateTokenRequest{
				Token: "invalid-token",
			},
			wantValid: false,
		},
		{
			name: "empty token",
			req: &pb.ValidateTokenRequest{
				Token: "",
			},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.ValidateToken(ctx, tt.req)
			if err != nil {
				t.Errorf("ValidateToken() error = %v", err)
				return
			}
			if resp.Valid != tt.wantValid {
				t.Errorf("ValidateToken() valid = %v, want %v", resp.Valid, tt.wantValid)
			}
			if tt.wantValid && resp.UserId != regResp.UserId {
				t.Errorf("ValidateToken() userId = %v, want %v", resp.UserId, regResp.UserId)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	store := storage.NewInMemoryStorage()
	service := NewAuthService(store)
	ctx := context.Background()

	// Register and login a user
	service.Register(ctx, &pb.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	})

	loginResp, _ := service.Login(ctx, &pb.LoginRequest{
		Username: "testuser",
		Password: "password123",
	})

	tests := []struct {
		name        string
		req         *pb.LogoutRequest
		wantSuccess bool
	}{
		{
			name: "valid logout",
			req: &pb.LogoutRequest{
				Token: loginResp.Token,
			},
			wantSuccess: true,
		},
		{
			name: "empty token",
			req: &pb.LogoutRequest{
				Token: "",
			},
			wantSuccess: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := service.Logout(ctx, tt.req)
			if err != nil {
				t.Errorf("Logout() error = %v", err)
				return
			}
			if resp.Success != tt.wantSuccess {
				t.Errorf("Logout() success = %v, want %v, message = %v", resp.Success, tt.wantSuccess, resp.Message)
			}
		})
	}

	// Verify token is invalidated after successful logout
	validateResp, _ := service.ValidateToken(ctx, &pb.ValidateTokenRequest{
		Token: loginResp.Token,
	})
	if validateResp.Valid {
		t.Error("Token should be invalid after logout")
	}
}
