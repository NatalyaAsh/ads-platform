package handler

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NatalyaAsh/ads-platform/services/auth-user/internal/service"
	authv1 "github.com/NatalyaAsh/ads-platform/services/auth-user/proto_gen/auth_user/v1"
)

// AuthGRPCHandler реализует gRPC интерфейс AuthUserServiceServer
type AuthGRPCHandler struct {
	authv1.UnimplementedAuthUserServiceServer
	authService *service.AuthService
}

// NewAuthGRPCHandler создаёт новый gRPC хендлер
func NewAuthGRPCHandler(authService *service.AuthService) *AuthGRPCHandler {
	return &AuthGRPCHandler{
		authService: authService,
	}
}

// Register регистрация пользователя
func (h *AuthGRPCHandler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	resp, err := h.authService.Register(&service.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authv1.RegisterResponse{
		UserId:       uint32(resp.UserID),
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

// Login вход пользователя
func (h *AuthGRPCHandler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	resp, err := h.authService.Login(&service.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		if errors.Is(err, service.ErrUserBlocked) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authv1.LoginResponse{
		UserId:       uint32(resp.UserID),
		Role:         resp.Role,
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

// ValidateToken проверка токена
func (h *AuthGRPCHandler) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	claims, err := h.authService.ValidateToken(req.AccessToken)
	if err != nil {
		return &authv1.ValidateTokenResponse{
			Valid: false,
		}, nil
	}

	return &authv1.ValidateTokenResponse{
		Valid:  true,
		UserId: uint32(claims.UserID),
		Role:   claims.Role,
	}, nil
}

// RefreshToken обновление токенов
func (h *AuthGRPCHandler) RefreshToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	accessToken, refreshToken, err := h.authService.RefreshTokens(req.RefreshToken)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	return &authv1.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Logout выход пользователя
func (h *AuthGRPCHandler) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	err := h.authService.Logout(req.RefreshToken)
	if err != nil {
		return &authv1.LogoutResponse{Success: false}, nil
	}

	return &authv1.LogoutResponse{Success: true}, nil
}

// GetUser получение пользователя
func (h *AuthGRPCHandler) GetUser(ctx context.Context, req *authv1.GetUserRequest) (*authv1.GetUserResponse, error) {
	user, err := h.authService.GetUserByID(uint(req.UserId))
	if err != nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &authv1.GetUserResponse{
		Id:        uint32(user.ID),
		Email:     user.Email,
		Role:      string(user.Role),
		IsBlocked: user.IsBlocked,
		CreatedAt: user.CreatedAt.String(),
	}, nil
}
