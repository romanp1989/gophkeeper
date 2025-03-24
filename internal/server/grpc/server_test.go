package grpc

import (
	"database/sql"
	"github.com/golang/mock/gomock"
	serverConfig "github.com/romanp1989/gophkeeper/internal/server/config"
	"github.com/romanp1989/gophkeeper/internal/server/db"
	"github.com/romanp1989/gophkeeper/internal/server/token"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	logger := zap.NewNop()
	cfg := &serverConfig.Config{
		Address: ":50051",
		Db: &db.Config{
			Dsn:             "dsn",
			MaxIdleConns:    1,
			MaxOpenConns:    1,
			MaxLifetimeConn: 10,
		},
		Token: &token.Config{
			Secret: "secret",
			Name:   "Authorization",
			Expire: time.Hour * 1,
		},
	}
	dbMock := &sql.DB{}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	srv := NewServer(cfg, dbMock, logger)

	assert.NotNil(t, srv)
	assert.Equal(t, cfg, srv.config)
	assert.NotNil(t, srv.grpcServer)
	assert.Equal(t, logger, srv.logger)
}

func Test_grpcServerSetup(t *testing.T) {
	logger := zap.NewNop()
	cfg := &serverConfig.Config{
		Address: ":50051",
		Db: &db.Config{
			Dsn:             "dsn",
			MaxIdleConns:    1,
			MaxOpenConns:    1,
			MaxLifetimeConn: 10,
		},
		Token: &token.Config{
			Secret: "secret",
			Name:   "Authorization",
			Expire: time.Hour * 1,
		},
	}
	dbMock := &sql.DB{}
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	server := grpcServerSetup(cfg, dbMock, logger)

	assert.NotNil(t, server)
}
