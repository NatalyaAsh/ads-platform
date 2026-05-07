package server

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/NatalyaAsh/ads-platform/services/ad-search/internal/handler"
	"github.com/NatalyaAsh/ads-platform/services/ad-search/internal/service"
	pb "github.com/NatalyaAsh/ads-platform/services/ad-search/proto_gen/ad_search/v1"
)

func RunGRPCServer(port string, adService *service.AdService) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	adHandler := handler.NewAdSearchHandler(adService)
	pb.RegisterAdSearchServiceServer(grpcServer, adHandler)

	log.Printf("gRPC server listening on port %s", port)

	return grpcServer.Serve(lis)
}
