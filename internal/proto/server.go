package proto

import (
	"fmt"
	"github.com/WildEgor/e-shop-auth/internal/configs"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type GRPCServer struct {
	server      *grpc.Server
	appConfig   *configs.AppConfig
	authService *AuthService
}

func NewGRPCServer(
	appConfig *configs.AppConfig,
	authService *AuthService,
) *GRPCServer {
	return &GRPCServer{
		appConfig:   appConfig,
		authService: authService,
	}
}

func (s *GRPCServer) Init() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", s.appConfig.RPCPort))
	if err != nil {
		slog.Error("cannot listen port - ", err)
		panic(err)
	}

	serv := grpc.NewServer()

	RegisterAuthServiceServer(serv, s.authService)

	s.server = serv

	go func() {
		// Run gRPC server
		if err := serv.Serve(listener); err != nil {
			slog.Error("error serve grpc", err)
			panic(err)
		}
	}()

	return nil
}

func (s *GRPCServer) Stop() {
	s.server.GracefulStop()

	slog.Info("stop gRPC server")
}
