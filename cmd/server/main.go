package server

import (
	serverConfig "github.com/romanp1989/gophkeeper/internal/server/config"
	"github.com/romanp1989/gophkeeper/internal/server/grpc"
	logger2 "github.com/romanp1989/gophkeeper/internal/server/logger"
	"go.uber.org/zap"
	"log"
)

func main() {
	logger, err := logger2.NewLogger(false)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := serverConfig.NewConfig()
	if err != nil {
		logger.Fatal("Error loading config", zap.Error(err))
	}

	server := grpc.NewServer(cfg, logger)
	err = server.Start()
	if err != nil {
		logger.Fatal("Error starting server", zap.Error(err))
	}
}
