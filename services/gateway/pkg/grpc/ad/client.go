package ad

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	adv1 "github.com/NatalyaAsh/ads-platform/services/gateway/proto_gen/ad_search/v1"
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
	req := &adv1.GetAdRequest{Id: id}
	return c.client.GetAd(ctx, req)
}

func (c *Client) ListCategories(ctx context.Context) (*adv1.ListCategoriesResponse, error) {
	req := &adv1.ListCategoriesRequest{}
	return c.client.ListCategories(ctx, req)
}

func (c *Client) CreateAd(ctx context.Context, title, description string, price float64, userID, categoryID string) (*adv1.CreateAdResponse, error) {
	req := &adv1.CreateAdRequest{
		Title:       title,
		Description: description,
		Price:       price,
		UserId:      userID,
		CategoryId:  categoryID,
	}
	return c.client.CreateAd(ctx, req)
}
