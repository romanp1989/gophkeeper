package server

import (
	"database/sql"
	serverConfig "github.com/romanp1989/gophkeeper/internal/server/config"
	dbService "github.com/romanp1989/gophkeeper/internal/server/db"
	"github.com/romanp1989/gophkeeper/internal/server/grpc"
	logger2 "github.com/romanp1989/gophkeeper/internal/server/logger"
	"go.uber.org/zap"
	"log"
)

func main() {
	var db *sql.DB

	logger, err := logger2.NewLogger(false)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := serverConfig.NewConfig()
	if err != nil {
		logger.Fatal("Error loading config", zap.Error(err))
	}

	db, err = dbService.InitDB(cfg.Db, logger)
	if err != nil {
		logger.Fatal("Error initializing database", zap.Error(err))
	}

	server := grpc.NewServer(cfg, db, logger)
	err = server.Start()
	if err != nil {
		logger.Fatal("Error starting server", zap.Error(err))
	}
}
