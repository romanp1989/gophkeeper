package handlers

import (
	"context"
	"errors"
	"fmt"
	serverConfig "github.com/romanp1989/gophkeeper/internal/server/config"
	"github.com/romanp1989/gophkeeper/internal/server/token"
	"github.com/romanp1989/gophkeeper/internal/server/user"
	"github.com/romanp1989/gophkeeper/pkg/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	proto.UnimplementedUsersServer
	config       *serverConfig.Config
	userService  *user.UserService
	tokenService token.Service
	logger       *zap.Logger
}

func NewUserHandler(config *serverConfig.Config, userService *user.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		config:      config,
		userService: userService,
		logger:      logger,
	}
}

func (h *UserHandler) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	var tokenAuth string

	userEntity, err := h.userService.RegisterUser(ctx, req.Login, req.Password)
	if err != nil {
		if errors.Is(err, fmt.Errorf("user already exists %s", req.Login)) {
			return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("user already exists %s", req.Login))
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	tokenAuth, err = h.tokenService.BuildToken(userEntity.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed auth: %s", err.Error()))
	}
	return &proto.RegisterResponse{AccessToken: tokenAuth}, nil
}

func (h *UserHandler) Login(ctx context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	userEntity, err := h.userService.LoginUser(ctx, req.Login, req.Password)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	tokenAuth, err := h.tokenService.BuildToken(userEntity.ID)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed auth: %s", err.Error()))
	}
	return &proto.LoginResponse{AccessToken: tokenAuth}, nil
}
