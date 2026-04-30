package server

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/NatalyaAsh/ads-platform/services/auth-user/internal/handler"
	"github.com/NatalyaAsh/ads-platform/services/auth-user/internal/service"
	authv1 "github.com/NatalyaAsh/ads-platform/services/auth-user/proto_gen/auth_user/v1"
)

func RunGRPCServer(port string, authService *service.AuthService) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	authHandler := handler.NewAuthGRPCHandler(authService)
	authv1.RegisterAuthUserServiceServer(grpcServer, authHandler)

	log.Printf("gRPC server listening on port %s", port)

	return grpcServer.Serve(lis)
}
