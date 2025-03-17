package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/romanp1989/gophkeeper/domain"
	"github.com/romanp1989/gophkeeper/pkg/consts"
	"github.com/romanp1989/gophkeeper/pkg/converter"
	"github.com/romanp1989/gophkeeper/pkg/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ISecretService interface {
	Get(ctx context.Context, secretID uint64, userID domain.UserID) (*domain.Secret, error)
	GetUserSecrets(ctx context.Context, userID domain.UserID) ([]*domain.Secret, error)
	Add(ctx context.Context, secret *domain.Secret) (*domain.Secret, error)
	Update(ctx context.Context, secret *domain.Secret) (*domain.Secret, error)
	Delete(ctx context.Context, secretID uint64, userID domain.UserID) error
}

type SecretHandler struct {
	proto.UnimplementedSecretsServer
	secretService ISecretService
	logger        *zap.Logger
}

func NewSecretHandler(secretService ISecretService, logger *zap.Logger) *SecretHandler {
	return &SecretHandler{
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
		if errors.Is(err, fmt.Errorf("secret not found id %d", in.Id)) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.GetUserSecretResponse{Secret: converter.SecretToProto(secret)}, nil
}

func (s *SecretHandler) GetUserSecretList(ctx context.Context, _ *emptypb.Empty) (*proto.GetUserSecretsResponse, error) {
	userID, err := extractUserID(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	secrets, err := s.secretService.GetUserSecrets(ctx, userID)
	if err != nil {
		if errors.Is(err, fmt.Errorf("user's secrets not found id %d", userID)) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.GetUserSecretsResponse{Secrets: converter.SecretsToProto(secrets)}, nil
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
		if errors.Is(err, fmt.Errorf("secret not found id %d", in.Id)) {
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
