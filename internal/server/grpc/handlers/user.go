package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/romanp1989/gophkeeper/domain"
	"github.com/romanp1989/gophkeeper/internal/server/token"
	"github.com/romanp1989/gophkeeper/pkg/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IUserService interface {
	RegisterUser(ctx context.Context, login string, password string) (*domain.User, error)
	LoginUser(ctx context.Context, login string, password string) (*domain.User, error)
}

type UserHandler struct {
	proto.UnimplementedUsersServer
	userService  IUserService
	tokenService token.Service
	logger       *zap.Logger
}

func NewUserHandler(userService IUserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
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
