package model

import (
	"testing"
	"time"
)

// запуск : go test ./internal/model/... -v

func TestUserValidation(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name: "valid user",
			user: User{
				Email: "test@example.com",
				Role:  RoleUser,
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			user: User{
				Email: "invalid-email",
				Role:  RoleUser,
			},
			wantErr: true,
		},
		{
			name: "invalid role",
			user: User{
				Email: "test@example.com",
				Role:  "superuser",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRefreshTokenValidation(t *testing.T) {
	rt := &RefreshToken{
		UserID:    1,
		Token:     "some-token",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := rt.Validate(); err != nil {
		t.Errorf("Validate() should pass, got error: %v", err)
	}

	if rt.IsExpired() {
		t.Error("IsExpired() should return false for future date")
	}

	if !rt.IsValid() {
		t.Error("IsValid() should return true for valid token")
	}

	rt.Revoke()
	if rt.IsValid() {
		t.Error("IsValid() should return false for revoked token")
	}
}
