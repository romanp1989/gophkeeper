package client

import (
	"github.com/romanp1989/gophkeeper/internal/client/config"
	"github.com/romanp1989/gophkeeper/internal/client/grpc"
	"github.com/romanp1989/gophkeeper/internal/client/tui/app"
	logger2 "github.com/romanp1989/gophkeeper/internal/server/logger"
	"go.uber.org/zap"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	logger, err := logger2.NewLogger(false)
	if err != nil {
		panic(err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("error loading config", zap.Error(err))
	}

	cfg.BuildVersion = buildVersion
	cfg.BuildDate = buildDate
	cfg.BuildCommit = buildCommit

	if err != nil {
		logger.Fatal("Error initializing model", zap.Error(err))
	}

	grpcClient, err := grpc.NewClientGRPC(cfg)
	if err != nil {
		logger.Fatal("Error initializing gRPC-client", zap.Error(err))
	}

	tuiApp := app.NewTuiApplication(grpcClient, cfg, logger)
	tuiApp.Start()
}
