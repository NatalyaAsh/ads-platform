package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authv1 "github.com/NatalyaAsh/ads-platform/services/auth-user/proto_gen/auth_user/v1"
)

func main() {
	// Подключаемся к gRPC серверу
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := authv1.NewAuthUserServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. РЕГИСТРАЦИЯ
	log.Println("=== Testing Register ===")
	registerResp, err := client.Register(ctx, &authv1.RegisterRequest{
		Email:    "test3@example.com",
		Password: "password256",
	})
	if err != nil {
		log.Printf("Register error: %v", err)
	} else {
		log.Printf("✅ Register success!")
		log.Printf("   User ID: %d", registerResp.UserId)
		log.Printf("   Access Token: %s...", registerResp.AccessToken[:50])
		log.Printf("   Refresh Token: %s...", registerResp.RefreshToken[:50])
	}

	// 2. ЛОГИН
	log.Println("\n=== Testing Login ===")
	loginResp, err := client.Login(ctx, &authv1.LoginRequest{
		Email:    "test3@example.com",
		Password: "password256",
	})
	if err != nil {
		log.Printf("Login error: %v", err)
	} else {
		log.Printf("✅ Login success!")
		log.Printf("   User ID: %d", loginResp.UserId)
		log.Printf("   Role: %s", loginResp.Role)
		log.Printf("   Access Token: %s...", loginResp.AccessToken[:50])
		log.Printf("   Refresh Token: %s...", loginResp.RefreshToken[:50])
	}

	// 3. ВАЛИДАЦИЯ ТОКЕНА
	if loginResp != nil {
		log.Println("\n=== Testing ValidateToken ===")
		validateResp, err := client.ValidateToken(ctx, &authv1.ValidateTokenRequest{
			AccessToken: loginResp.AccessToken,
		})
		if err != nil {
			log.Printf("ValidateToken error: %v", err)
		} else {
			log.Printf("✅ Token valid: %v", validateResp.Valid)
			if validateResp.Valid {
				log.Printf("   User ID from token: %d", validateResp.UserId)
				log.Printf("   Role from token: %s", validateResp.Role)
			}
		}
	}

	// 4. ОБНОВЛЕНИЕ ТОКЕНОВ
	if loginResp != nil {
		log.Println("\n=== Testing RefreshToken ===")
		refreshResp, err := client.RefreshToken(ctx, &authv1.RefreshTokenRequest{
			RefreshToken: loginResp.RefreshToken,
		})
		if err != nil {
			log.Printf("RefreshToken error: %v", err)
		} else {
			log.Printf("✅ New tokens generated!")
			log.Printf("   New Access Token: %s...", refreshResp.AccessToken[:50])
			log.Printf("   New Refresh Token: %s...", refreshResp.RefreshToken[:50])
		}
	}

	// 5. ПОЛУЧЕНИЕ ПОЛЬЗОВАТЕЛЯ
	if loginResp != nil {
		log.Println("\n=== Testing GetUser ===")
		userResp, err := client.GetUser(ctx, &authv1.GetUserRequest{
			UserId: loginResp.UserId,
		})
		if err != nil {
			log.Printf("GetUser error: %v", err)
		} else {
			log.Printf("✅ User found!")
			log.Printf("   ID: %d", userResp.Id)
			log.Printf("   Email: %s", userResp.Email)
			log.Printf("   Role: %s", userResp.Role)
			log.Printf("   Is Blocked: %v", userResp.IsBlocked)
			log.Printf("   Created At: %s", userResp.CreatedAt)
		}
	}

	// 6. ВЫХОД (LOGOUT)
	if loginResp != nil {
		log.Println("\n=== Testing Logout ===")
		logoutResp, err := client.Logout(ctx, &authv1.LogoutRequest{
			RefreshToken: loginResp.RefreshToken,
		})
		if err != nil {
			log.Printf("Logout error: %v", err)
		} else {
			log.Printf("✅ Logout success: %v", logoutResp.Success)
		}
	}

	log.Println("\n=== All tests completed ===")
}
