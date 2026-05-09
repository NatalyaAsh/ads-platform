package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"

	"github.com/NatalyaAsh/ads-platform/services/gateway/pkg/grpc/ad"
	"github.com/NatalyaAsh/ads-platform/services/gateway/pkg/grpc/auth"
)

var adClient *ad.Client
var authClient *auth.Client

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func init() {
	var err error

	adAddr := os.Getenv("AD_SERVICE_ADDR")
	if adAddr == "" {
		adAddr = "localhost:50052"
	}
	adClient, err = ad.NewClient(adAddr)
	if err != nil {
		log.Fatalf("Failed to connect to ad service: %v", err)
	}

	authAddr := os.Getenv("AUTH_SERVICE_ADDR")
	if authAddr == "" {
		authAddr = "localhost:50051"
	}
	authClient, err = auth.NewClient(authAddr)
	if err != nil {
		log.Fatalf("Failed to connect to auth service: %v", err)
	}
}

// Middleware для извлечения токена из заголовка Authorization
func withAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		token := ""
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}
		ctx := context.WithValue(r.Context(), "token", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Вспомогательная функция для получения userID из токена
func getUserIDFromContext(ctx context.Context) (uint32, error) {
	token, ok := ctx.Value("token").(string)
	if !ok || token == "" {
		return 0, fmt.Errorf("missing or empty token")
	}
	resp, err := authClient.ValidateToken(ctx, token)
	if err != nil {
		return 0, fmt.Errorf("token validation failed: %w", err)
	}
	if !resp.Valid {
		return 0, fmt.Errorf("invalid token")
	}
	return resp.UserId, nil
}

func main() {
	defer adClient.Close()
	defer authClient.Close()

	// ========== ТИПЫ ==========
	categoryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Category",
		Fields: graphql.Fields{
			"id":   &graphql.Field{Type: graphql.ID},
			"name": &graphql.Field{Type: graphql.String},
			"slug": &graphql.Field{Type: graphql.String},
		},
	})

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
					ad, _ := p.Source.(map[string]interface{})
					categoryID, _ := ad["categoryId"].(string)
					return map[string]interface{}{
						"id":   categoryID,
						"name": "Category",
						"slug": "category",
					}, nil
				},
			},
		},
	})

	authPayloadType := graphql.NewObject(graphql.ObjectConfig{
		Name: "AuthPayload",
		Fields: graphql.Fields{
			"userId":       &graphql.Field{Type: graphql.ID},
			"role":         &graphql.Field{Type: graphql.String},
			"accessToken":  &graphql.Field{Type: graphql.String},
			"refreshToken": &graphql.Field{Type: graphql.String},
		},
	})

	createAdInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "CreateAdInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"title":       &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"description": &graphql.InputObjectFieldConfig{Type: graphql.String},
			"price":       &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Float)},
			"categoryId":  &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.ID)},
		},
	})

	updateAdInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "UpdateAdInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"title":       &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"description": &graphql.InputObjectFieldConfig{Type: graphql.String},
			"price":       &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.Float)},
			"categoryId":  &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.ID)},
		},
	})

	// ========== QUERY ==========
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
					return adClient.GetAd(context.Background(), id)
				},
			},
			"categories": &graphql.Field{
				Type: graphql.NewList(categoryType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					resp, err := adClient.ListCategories(context.Background())
					if err != nil {
						return nil, err
					}
					var out []interface{}
					for _, cat := range resp.Categories {
						out = append(out, map[string]interface{}{
							"id":   cat.Id,
							"name": cat.Name,
							"slug": cat.Slug,
						})
					}
					return out, nil
				},
			},
			"userAds": &graphql.Field{
				Type: graphql.NewList(adType),
				Args: graphql.FieldConfigArgument{
					"userId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID := p.Args["userId"].(string)
					ads, err := adClient.GetUserAds(context.Background(), userID)
					if err != nil {
						return nil, err
					}
					var result []interface{}
					for _, ad := range ads {
						result = append(result, map[string]interface{}{
							"id":          ad.Id,
							"title":       ad.Title,
							"description": ad.Description,
							"price":       ad.Price,
							"userId":      ad.UserId,
							"categoryId":  ad.CategoryId,
							"status":      ad.Status,
							"views":       ad.Views,
							"createdAt":   ad.CreatedAt,
							"updatedAt":   ad.UpdatedAt,
						})
					}
					return result, nil
				},
			},
			"ads": &graphql.Field{
				Type: graphql.NewList(adType),
				Args: graphql.FieldConfigArgument{
					"limit":  &graphql.ArgumentConfig{Type: graphql.Int},
					"offset": &graphql.ArgumentConfig{Type: graphql.Int},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					limit, _ := p.Args["limit"].(int)
					offset, _ := p.Args["offset"].(int)
					if limit == 0 {
						limit = 20
					}
					ads, err := adClient.ListAds(context.Background(), int32(limit), int32(offset))
					if err != nil {
						return nil, err
					}
					return ads, nil
				},
			},
		},
	})

	// ========== MUTATION ==========
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
			"createAd": &graphql.Field{
				Type: adType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(createAdInputType)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					userID, err := getUserIDFromContext(p.Context)
					log.Printf("🔍 Создаём объявление от userId = %d", userID)
					if err != nil {
						return nil, err
					}
					input := p.Args["input"].(map[string]interface{})
					title := input["title"].(string)
					description, _ := input["description"].(string)
					price := input["price"].(float64)
					categoryID := input["categoryId"].(string)

					resp, err := adClient.CreateAd(context.Background(), title, description, price, fmt.Sprintf("%d", userID), categoryID)
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
			"updateAd": &graphql.Field{
				Type: adType,
				Args: graphql.FieldConfigArgument{
					"id":    &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(updateAdInputType)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					//userID, err := getUserIDFromContext(p.Context)
					_, err := getUserIDFromContext(p.Context)
					if err != nil {
						return nil, err
					}
					id := p.Args["id"].(string)
					input := p.Args["input"].(map[string]interface{})
					title := input["title"].(string)
					description, _ := input["description"].(string)
					price := input["price"].(float64)
					categoryID := input["categoryId"].(string)

					// Вызов gRPC метода UpdateAd (требует прав: владелец или админ)
					// Для простоты пока не проверяем права, но можно передать userID и проверить внутри сервиса
					resp, err := adClient.UpdateAd(context.Background(), id, title, description, price, categoryID)
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
						"updatedAt":   resp.UpdatedAt,
					}, nil
				},
			},
			"deleteAd": &graphql.Field{
				Type: graphql.Boolean,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					_, err := getUserIDFromContext(p.Context) // проверяем авторизацию
					if err != nil {
						return false, err
					}
					id := p.Args["id"].(string)
					_, err = adClient.DeleteAd(context.Background(), id)
					return err == nil, err
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

	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// // Применяем middleware для извлечения токена
	// http.Handle("/graphql", withAuth(h))
	// Оборачиваем в CORS и авторизацию
	http.Handle("/graphql", corsMiddleware(withAuth(h)))

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
