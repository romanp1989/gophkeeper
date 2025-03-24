package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/romanp1989/gophkeeper/domain"
	"github.com/romanp1989/gophkeeper/pkg/consts"
	"github.com/romanp1989/gophkeeper/pkg/proto"
	"github.com/romanp1989/gophkeeper/tests/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"testing"
)

func TestSecretHandler_GetUserSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockISecretService(ctrl)
	logger := zap.NewNop()
	handler := NewSecretHandler(mockService, logger)

	tests := []struct {
		name      string
		setupMock func()
		ctx       context.Context
		input     *proto.GetUserSecretRequest
		expectErr string
	}{
		{
			name: "Success",
			setupMock: func() {
				mockService.EXPECT().Get(gomock.Any(), uint64(1), uint64(123)).Return(&domain.Secret{ID: 1}, nil).Times(1)
			},
			ctx: context.WithValue(context.Background(), consts.UserIDKeyCtx, uint64(123)),
			input: &proto.GetUserSecretRequest{
				Id: 1,
			},
			expectErr: "",
		},
		{
			name: "Error_NotFound",
			setupMock: func() {
				mockService.EXPECT().Get(gomock.Any(), uint64(1), uint64(123)).Return(nil, errors.New("secret not found")).Times(1)
			},
			ctx: context.WithValue(context.Background(), consts.UserIDKeyCtx, uint64(123)),
			input: &proto.GetUserSecretRequest{
				Id: 1,
			},
			expectErr: "rpc error: code = Internal desc = secret not found",
		},
		{
			name:      "Error_MissingUserID",
			setupMock: func() {},
			ctx:       context.Background(),
			input: &proto.GetUserSecretRequest{
				Id: 1,
			},
			expectErr: "rpc error: code = Internal desc = failed to extract user id from context",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			_, err := handler.GetUserSecret(tc.ctx, tc.input)
			if tc.expectErr != "" {
				assert.EqualError(t, err, tc.expectErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func TestSecretHandler_GetUserSecretList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockISecretService(ctrl)
	logger := zap.NewNop()
	handler := NewSecretHandler(mockService, logger)

	tests := []struct {
		name      string
		setupMock func()
		ctx       context.Context
		input     *emptypb.Empty
		expectErr string
	}{
		{
			name: "Success",
			setupMock: func() {
				mockService.EXPECT().GetUserSecrets(gomock.Any(), uint64(123)).Return([]domain.Secret{}, nil).Times(1)
			},
			ctx: metadata.NewIncomingContext(
				context.WithValue(context.Background(), consts.UserIDKeyCtx, uint64(123)),
				metadata.New(nil),
			),
			input:     &emptypb.Empty{},
			expectErr: "",
		},
		{
			name:      "Error_MissingUserID",
			setupMock: func() {},
			ctx:       context.Background(),
			input:     &emptypb.Empty{},
			expectErr: "rpc error: code = Internal desc = failed to extract user id from context",
		},
		{
			name: "Error_Internal",
			setupMock: func() {
				mockService.EXPECT().GetUserSecrets(gomock.Any(), uint64(123)).Return(nil, errors.New("internal error")).Times(1)
			},
			ctx: metadata.NewIncomingContext(
				context.WithValue(context.Background(), consts.UserIDKeyCtx, uint64(123)),
				metadata.New(nil),
			),
			input:     &emptypb.Empty{},
			expectErr: "rpc error: code = Internal desc = internal error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			_, err := handler.GetUserSecrets(tc.ctx, tc.input)
			if tc.expectErr != "" {
				assert.EqualError(t, err, tc.expectErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSecretHandler_SaveUserSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockISecretService(ctrl)
	logger := zap.NewNop()
	handler := NewSecretHandler(mockService, logger)

	tests := []struct {
		name      string
		setupMock func()
		ctx       context.Context
		input     *proto.SaveUserSecretRequest
		expectErr string
	}{
		{
			name: "Success_Create",
			setupMock: func() {
				mockService.EXPECT().Add(gomock.Any(), gomock.Any()).Return(&domain.Secret{ID: 1}, nil).Times(1)
			},
			ctx: metadata.NewIncomingContext(
				context.WithValue(context.Background(), consts.UserIDKeyCtx, uint64(123)),
				metadata.New(map[string]string{consts.ClientIDHeader: "456"}),
			),
			input: &proto.SaveUserSecretRequest{
				Secret: &proto.Secret{
					Title:      "Test Secret",
					Metadata:   "Test Metadata",
					SecretType: proto.SecretType_SECRET_TYPE_TEXT,
				},
			},
			expectErr: "",
		},
		{
			name:      "Error_MissingUserID",
			setupMock: func() {},
			ctx:       context.Background(),
			input: &proto.SaveUserSecretRequest{
				Secret: &proto.Secret{},
			},
			expectErr: "rpc error: code = Internal desc = failed to extract user id from context",
		},
		{
			name: "Error_CreateSecret",
			setupMock: func() {
				mockService.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil, errors.New("create error")).Times(1)
			},
			ctx: metadata.NewIncomingContext(
				context.WithValue(context.Background(), consts.UserIDKeyCtx, uint64(123)),
				metadata.New(map[string]string{consts.ClientIDHeader: "456"}),
			),
			input: &proto.SaveUserSecretRequest{
				Secret: &proto.Secret{},
			},
			expectErr: "rpc error: code = Internal desc = create error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			_, err := handler.SaveUserSecret(tc.ctx, tc.input)
			if tc.expectErr != "" {
				assert.EqualError(t, err, tc.expectErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSecretHandler_DeleteUserSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockISecretService(ctrl)
	logger := zap.NewNop()
	handler := NewSecretHandler(mockService, logger)

	tests := []struct {
		name      string
		setupMock func()
		ctx       context.Context
		input     *proto.DeleteUserSecretRequest
		expectErr string
	}{
		{
			name: "Success",
			setupMock: func() {
				mockService.EXPECT().Delete(gomock.Any(), uint64(1), uint64(123)).Return(nil).Times(1)
			},
			ctx: metadata.NewIncomingContext(
				context.WithValue(context.Background(), consts.UserIDKeyCtx, uint64(123)),
				metadata.New(nil),
			),
			input:     &proto.DeleteUserSecretRequest{Id: 1},
			expectErr: "",
		},
		{
			name: "Error_NotFound",
			setupMock: func() {
				mockService.EXPECT().Delete(gomock.Any(), uint64(1), uint64(123)).Return(fmt.Errorf("secret not found (id=%d)", 1)).Times(1)
			},
			ctx: metadata.NewIncomingContext(
				context.WithValue(context.Background(), consts.UserIDKeyCtx, uint64(123)),
				metadata.New(nil),
			),
			input:     &proto.DeleteUserSecretRequest{Id: 1},
			expectErr: "rpc error: code = Internal desc = secret not found (id=1)",
		},
		{
			name:      "Error_MissingUserID",
			setupMock: func() {},
			ctx:       context.Background(),
			input:     &proto.DeleteUserSecretRequest{Id: 1},
			expectErr: "rpc error: code = Internal desc = failed to extract user id from context",
		},
		{
			name: "Error_Internal",
			setupMock: func() {
				mockService.EXPECT().Delete(gomock.Any(), uint64(1), uint64(123)).Return(errors.New("internal error")).Times(1)
			},
			ctx: metadata.NewIncomingContext(
				context.WithValue(context.Background(), consts.UserIDKeyCtx, uint64(123)),
				metadata.New(nil),
			),
			input:     &proto.DeleteUserSecretRequest{Id: 1},
			expectErr: "rpc error: code = Internal desc = internal error",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMock()

			_, err := handler.DeleteUserSecret(tc.ctx, tc.input)
			if tc.expectErr != "" {
				assert.EqualError(t, err, tc.expectErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
