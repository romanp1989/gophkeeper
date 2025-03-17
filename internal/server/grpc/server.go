package grpc

import (
	"context"
	"database/sql"
	serverConfig "github.com/romanp1989/gophkeeper/internal/server/config"
	"github.com/romanp1989/gophkeeper/internal/server/grpc/handlers"
	"github.com/romanp1989/gophkeeper/internal/server/grpc/interceptors"
	"github.com/romanp1989/gophkeeper/internal/server/secret"
	"github.com/romanp1989/gophkeeper/internal/server/token"
	"github.com/romanp1989/gophkeeper/internal/server/user"
	"github.com/romanp1989/gophkeeper/pkg/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	config     *serverConfig.Config
	grpcServer *grpc.Server
	logger     *zap.Logger
}

func NewServer(config *serverConfig.Config, db *sql.DB, logger *zap.Logger) *Server {
	grpcServer := grpcServerSetup(config, db, logger)
	return &Server{
		config:     config,
		grpcServer: grpcServer,
		logger:     logger,
	}
}

// grpcServerSetup Конфигурирование GRPC сервера
func grpcServerSetup(cfg *serverConfig.Config, db *sql.DB, logger *zap.Logger) *grpc.Server {
	tokenService := token.NewJwtService(cfg.Token)
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(interceptors.Authentication(tokenService)),
	}

	tlsCredentials, err := cfg.LoadTLSConfig("ca-cert.pem", "server-cert.pem", "server-key.pem")
	if err != nil {
		logger.Fatal("Failed to load TLS credentials", zap.Error(err))
	}

	opts = append(opts, grpc.Creds(tlsCredentials))

	server := grpc.NewServer(opts...)

	userRepository := user.NewUserRepository(db)
	secretRepository := secret.NewSecretRepository(db)

	proto.RegisterUsersServer(server, handlers.NewUserHandler(user.NewUserService(userRepository), logger))
	proto.RegisterSecretsServer(server, handlers.NewSecretHandler(secret.NewSecretService(secretRepository), logger))

	return server
}

// Start запускает GPRC сервер приложения
func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		s.logger.Fatal("failed to listen", zap.Error(err))
		return err
	}

	go func() {
		if err := s.grpcServer.Serve(listener); err != nil {
			s.logger.Fatal("failed to serve", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	sig := <-quit
	s.logger.Info("interrupt received signal", zap.String("signal", sig.String()))
	s.shutdown()

	return nil
}

// shutdown останавливает GRPC сервер,осуществляя graceful завершение
func (s *Server) shutdown() {
	s.logger.Info("Shutting down server...")
	stopped := make(chan struct{})
	stopCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	go func() {
		s.grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-stopped:
		s.logger.Info("Server shutdown successful")
	case <-stopCtx.Done():
		s.logger.Info("Shutdown timeout exceeded")
	}
}
