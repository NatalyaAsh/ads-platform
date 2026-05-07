package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"github.com/NatalyaAsh/ads-platform/services/gateway/pkg/grpc/ad"
	"github.com/NatalyaAsh/ads-platform/services/gateway/pkg/grpc/auth"
)

var adClient *ad.Client
var authClient *auth.Client

func init() {
	var err error

	// Подключаемся к Ad Search Service
	adAddr := os.Getenv("AD_SERVICE_ADDR")
	if adAddr == "" {
		adAddr = "localhost:50052"
	}
	adClient, err = ad.NewClient(adAddr)
	if err != nil {
		log.Fatalf("Failed to connect to ad service: %v", err)
	}

	// Подключаемся к Auth Service
	authAddr := os.Getenv("AUTH_SERVICE_ADDR")
	if authAddr == "" {
		authAddr = "localhost:50051"
	}
	authClient, err = auth.NewClient(authAddr)
	if err != nil {
		log.Fatalf("Failed to connect to auth service: %v", err)
	}
}

func main() {
	// Закрываем соединения при завершении
	defer adClient.Close()
	defer authClient.Close()

	// Тип Category
	categoryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Category",
		Fields: graphql.Fields{
			"id":   &graphql.Field{Type: graphql.ID},
			"name": &graphql.Field{Type: graphql.String},
			"slug": &graphql.Field{Type: graphql.String},
		},
	})

	// Тип Ad
	adType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Ad",
		Fields: graphql.Fields{
			"id":          &graphql.Field{Type: graphql.ID},
			"title":       &graphql.Field{Type: graphql.String},
			"description": &graphql.Field{Type: graphql.String},
			"price":       &graphql.Field{Type: graphql.Float},
			"userId":      &graphql.Field{Type: graphql.ID},
			"categoryId":  &graphql.Field{Type: graphql.ID},
			"status":      &graphql.Field{Type: graphql.String},
			"views":       &graphql.Field{Type: graphql.Int},
			"createdAt":   &graphql.Field{Type: graphql.String},
			"updatedAt":   &graphql.Field{Type: graphql.String},
			"category": &graphql.Field{
				Type: categoryType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					ad, ok := p.Source.(map[string]interface{})
					if !ok {
						return nil, nil
					}
					categoryID, ok := ad["categoryId"].(string)
					if !ok {
						return nil, nil
					}
					// Заглушка — реальную категорию можно подгрузить отдельным запросом
					return map[string]interface{}{
						"id":   categoryID,
						"name": "Category",
						"slug": "category",
					}, nil
				},
			},
		},
	})

	// Тип AuthPayload
	authPayloadType := graphql.NewObject(graphql.ObjectConfig{
		Name: "AuthPayload",
		Fields: graphql.Fields{
			"userId":       &graphql.Field{Type: graphql.ID},
			"role":         &graphql.Field{Type: graphql.String},
			"accessToken":  &graphql.Field{Type: graphql.String},
			"refreshToken": &graphql.Field{Type: graphql.String},
		},
	})

	// Query
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"ad": &graphql.Field{
				Type: adType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id := p.Args["id"].(string)
					resp, err := adClient.GetAd(context.Background(), id)
					if err != nil {
						return nil, err
					}
					return map[string]interface{}{
						"id":          resp.Id,
						"title":       resp.Title,
						"description": resp.Description,
						"price":       resp.Price,
						"userId":      resp.UserId,
						"categoryId":  resp.CategoryId,
						"status":      resp.Status,
						"views":       resp.Views,
						"createdAt":   resp.CreatedAt,
						"updatedAt":   resp.UpdatedAt,
					}, nil
				},
			},
			"categories": &graphql.Field{
				Type: graphql.NewList(categoryType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					resp, err := adClient.ListCategories(context.Background())
					if err != nil {
						return nil, err
					}
					var categories []interface{}
					for _, cat := range resp.Categories {
						categories = append(categories, map[string]interface{}{
							"id":   cat.Id,
							"name": cat.Name,
							"slug": cat.Slug,
						})
					}
					return categories, nil
				},
			},
		},
	})

	// Mutation
	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"login": &graphql.Field{
				Type: authPayloadType,
				Args: graphql.FieldConfigArgument{
					"email":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"password": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					email := p.Args["email"].(string)
					password := p.Args["password"].(string)
					resp, err := authClient.Login(context.Background(), email, password)
					if err != nil {
						return nil, err
					}
					return map[string]interface{}{
						"userId":       resp.UserId,
						"role":         resp.Role,
						"accessToken":  resp.AccessToken,
						"refreshToken": resp.RefreshToken,
					}, nil
				},
			},
			"register": &graphql.Field{
				Type: authPayloadType,
				Args: graphql.FieldConfigArgument{
					"email":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"password": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					email := p.Args["email"].(string)
					password := p.Args["password"].(string)
					resp, err := authClient.Register(context.Background(), email, password)
					if err != nil {
						return nil, err
					}
					return map[string]interface{}{
						"userId":       resp.UserId,
						"role":         "user",
						"accessToken":  resp.AccessToken,
						"refreshToken": resp.RefreshToken,
					}, nil
				},
			},
		},
	})

	// Схема
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
	if err != nil {
		log.Fatalf("Failed to create schema: %v", err)
	}

	// HTTP хендлер
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql", h)

	// Graceful shutdown
	go func() {
		log.Println("GraphQL server running on http://localhost:8080/graphql")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gateway...")
}
