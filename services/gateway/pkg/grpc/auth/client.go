package auth

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authv1 "github.com/NatalyaAsh/ads-platform/services/gateway/proto_gen/auth_user/v1"
)

type Client struct {
	conn   *grpc.ClientConn
	client authv1.AuthUserServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	log.Printf("Connected to Auth Service at %s", addr)

	return &Client{
		conn:   conn,
		client: authv1.NewAuthUserServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Login(ctx context.Context, email, password string) (*authv1.LoginResponse, error) {
	req := &authv1.LoginRequest{
		Email:    email,
		Password: password,
	}
	return c.client.Login(ctx, req)
}

func (c *Client) Register(ctx context.Context, email, password string) (*authv1.RegisterResponse, error) {
	req := &authv1.RegisterRequest{
		Email:    email,
		Password: password,
	}
	return c.client.Register(ctx, req)
}

func (c *Client) ValidateToken(ctx context.Context, token string) (*authv1.ValidateTokenResponse, error) {
	req := &authv1.ValidateTokenRequest{
		AccessToken: token,
	}
	return c.client.ValidateToken(ctx, req)
}
