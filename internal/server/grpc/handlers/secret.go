package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/romanp1989/gophkeeper/internal/server/config"
	"github.com/romanp1989/gophkeeper/internal/server/domain"
	"github.com/romanp1989/gophkeeper/internal/server/secret"
	"github.com/romanp1989/gophkeeper/pkg/consts"
	"github.com/romanp1989/gophkeeper/pkg/converter"
	"github.com/romanp1989/gophkeeper/pkg/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
)

type SecretHandler struct {
	proto.UnimplementedSecretsServer
	config        *config.Config
	secretService *secret.Service
	logger        *zap.Logger
}

func NewSecretHandler(config *config.Config, secretService *secret.Service, logger *zap.Logger) *SecretHandler {
	return &SecretHandler{
		config:        config,
		secretService: secretService,
		logger:        logger,
	}
}

func (s *SecretHandler) GetUserSecret(ctx context.Context, in *proto.GetUserSecretRequest) (*proto.GetUserSecretResponse, error) {
	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	secret, err := s.secretService.Get(ctx, in.Id, userID)
	if err != nil {
		if errors.Is(err, fmt.Errorf("secret not found id %w", in.Id)) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.GetUserSecretResponse{Secret: converter.SecretToProto(secret)}, nil
}

func (s *SecretHandler) GetUserSecretList(ctx context.Context, in *proto.GetUserSecretsResponse) (*proto.GetUserSecretsResponse, error) {
	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	secrets, err := s.secretService.GetUserSecrets(ctx, userID)
	if err != nil {
		if errors.Is(err, fmt.Errorf("user's secrets not found id %w", userID)) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.GetUserSecretsResponse{Secret: converter.SecretsToProto(secrets)}, nil
}

func (s *SecretHandler) SaveUserSecret(ctx context.Context, in *proto.SaveUserSecretRequest) (*emptypb.Empty, error) {
	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	secretEntity := converter.ProtoToSecret(in.Secret)
	secretEntity.UserID = userID

	if secretEntity.ID > 0 {
		_, err = s.secretService.Update(ctx, secretEntity)
	} else {
		_, err = s.secretService.Add(ctx, secretEntity)
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	//@TODO
	//var clientID uint64
	//clientID, err = extractClientID(ctx)
	//if err != nil {
	//	s.logger.Error("failed to extract client ID", zap.Error(err))
	//	return &emptypb.Empty{}, err
	//}

	//isUpdated := secretEntity.ID > 0
	//err = s.notificationHandler.notifyClients(userID, clientID, secretEntity.ID, isUpdated)
	//if err != nil {
	//	s.logger.Error("failed to notify clients", zap.Error(err))
	//}

	return &empty.Empty{}, nil
}

func (s *SecretHandler) DeleteUserSecret(ctx context.Context, in *proto.DeleteUserSecretRequest) (*emptypb.Empty, error) {
	var err error
	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = s.secretService.Delete(ctx, in.Id, userID)
	if err != nil {
		if errors.Is(err, fmt.Errorf("secret not found id %w", in.Id)) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

// extractUserID получает userID из контекста
func extractUserID(ctx context.Context) (domain.UserID, error) {
	uid := ctx.Value(consts.UserIDKeyCtx)
	userID, ok := uid.(domain.UserID)
	if !ok {
		return 0, errors.New("failed to extract user id from context")
	}
	return userID, nil
}

// extractClientID получает id клиента из контекста
func extractClientID(ctx context.Context) (uint64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, errors.New("failed to get metadata")
	}
	values := md.Get(consts.ClientIDHeader)
	if len(values) == 0 {
		return 0, errors.New("missing client id metadata")
	}
	id, err := strconv.Atoi(values[0])
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}
