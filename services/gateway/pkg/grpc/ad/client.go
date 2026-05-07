package ad

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	adv1 "github.com/NatalyaAsh/ads-platform/services/gateway/internal/pb/ad"
)

type Client struct {
	conn   *grpc.ClientConn
	client adv1.AdSearchServiceClient
}

func NewClient(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ad service: %w", err)
	}

	log.Printf("Connected to Ad Search Service at %s", addr)

	return &Client{
		conn:   conn,
		client: adv1.NewAdSearchServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetAd(ctx context.Context, id string) (*adv1.GetAdResponse, error) {
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid ad id: %w", err)
	}
	req := &adv1.GetAdRequest{Id: uint32(idUint)}
	return c.client.GetAd(ctx, req)
}

func (c *Client) ListCategories(ctx context.Context) (*adv1.ListCategoriesResponse, error) {
	req := &adv1.ListCategoriesRequest{}
	return c.client.ListCategories(ctx, req)
}

func (c *Client) CreateAd(ctx context.Context, title, description string, price float64, userID, categoryID string) (*adv1.CreateAdResponse, error) {
	userUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}
	catUint, err := strconv.ParseUint(categoryID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid category id: %w", err)
	}
	req := &adv1.CreateAdRequest{
		Title:       title,
		Description: description,
		Price:       price,
		UserId:      uint32(userUint),
		CategoryId:  uint32(catUint),
	}
	return c.client.CreateAd(ctx, req)
}

// UpdateAd обновляет объявление
func (c *Client) UpdateAd(ctx context.Context, id, title, description string, price float64, categoryID string) (*adv1.UpdateAdResponse, error) {
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid ad id: %w", err)
	}
	catUint, err := strconv.ParseUint(categoryID, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid category id: %w", err)
	}
	req := &adv1.UpdateAdRequest{
		Id:          uint32(idUint),
		Title:       title,
		Description: description,
		Price:       price,
		CategoryId:  uint32(catUint),
	}
	return c.client.UpdateAd(ctx, req)
}

// DeleteAd удаляет объявление
func (c *Client) DeleteAd(ctx context.Context, id string) (*adv1.DeleteAdResponse, error) {
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid ad id: %w", err)
	}
	req := &adv1.DeleteAdRequest{Id: uint32(idUint)}
	return c.client.DeleteAd(ctx, req)
}
