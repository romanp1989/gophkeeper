package secret

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/romanp1989/gophkeeper/domain"
	"github.com/romanp1989/gophkeeper/tests/mocks"
	"testing"
)

func TestSecretService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockISecretRepository(ctrl)
	service := NewSecretService(mockRepo)

	ctx := context.Background()
	testSecret := domain.Secret{
		ID:         1,
		Title:      "Test Secret",
		UserID:     1,
		SecretType: string(domain.TextSecret),
	}

	tests := []struct {
		name      string
		testFunc  func(t *testing.T)
		expectErr bool
	}{
		{
			name: "Get_Success",
			testFunc: func(t *testing.T) {
				mockRepo.EXPECT().GetByID(ctx, uint64(1), uint64(1)).Return(testSecret, nil)

				secret, err := service.Get(ctx, 1, 1)
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if secret.ID != testSecret.ID {
					t.Errorf("Expected secret ID %v, got %v", testSecret.ID, secret.ID)
				}
			},
			expectErr: false,
		},
		{
			name: "Get_Fail_NotFound",
			testFunc: func(t *testing.T) {
				mockRepo.EXPECT().GetByID(ctx, uint64(1), uint64(1)).Return(nil, sql.ErrNoRows)

				_, err := service.Get(ctx, 1, 1)
				if err == nil || err.Error() != "secret not found (id=1)" {
					t.Errorf("Expected error 'secret not found (id=1)', got %v", err)
				}
			},
			expectErr: true,
		},
		{
			name: "Get_Fail",
			testFunc: func(t *testing.T) {
				mockRepo.EXPECT().GetByID(ctx, uint64(1), uint64(1)).Return(nil, errors.New("some error"))

				_, err := service.Get(ctx, 1, 1)
				if err == nil || err.Error() != "some error" {
					t.Errorf("Expected error 'some error', got %v", err)
				}
			},
			expectErr: true,
		},
		{
			name: "GetAllByUserID_Success",
			testFunc: func(t *testing.T) {
				secrets := []domain.Secret{testSecret}
				mockRepo.EXPECT().GetAllByUserID(ctx, uint64(1)).Return(secrets, nil)

				result, err := service.GetUserSecrets(ctx, 1)
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if len(result) != 1 {
					t.Errorf("Expected 1 secret, got %d", len(result))
				}
			},
			expectErr: false,
		},
		{
			name: "GetAllByUserID_Fail_EmptyList",
			testFunc: func(t *testing.T) {
				mockRepo.EXPECT().GetAllByUserID(ctx, uint64(1)).Return(nil, nil)

				_, err := service.GetUserSecrets(ctx, 1)
				if err == nil || err.Error() != "no secrets found" {
					t.Errorf("Expected error 'no secrets found', got %v", err)
				}
			},
			expectErr: true,
		},
		{
			name: "GetAllByUserID_Fail",
			testFunc: func(t *testing.T) {
				mockRepo.EXPECT().GetAllByUserID(ctx, uint64(1)).Return(nil, errors.New("some error"))

				_, err := service.GetUserSecrets(ctx, 1)
				if err == nil || err.Error() != "some error" {
					t.Errorf("Expected error 'some error', got %v", err)
				}
			},
			expectErr: true,
		},
		{
			name: "Add_Success",
			testFunc: func(t *testing.T) {
				mockRepo.EXPECT().Create(ctx, &testSecret).Return(uint64(1), nil)

				createdSecret, err := service.Add(ctx, &testSecret)
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if createdSecret.ID != 1 {
					t.Errorf("Expected secret ID 1, got %v", createdSecret.ID)
				}
			},
			expectErr: false,
		},
		{
			name: "Add_Fail",
			testFunc: func(t *testing.T) {
				mockRepo.EXPECT().Create(ctx, testSecret).Return(uint64(1), errors.New("some error"))

				_, err := service.Add(ctx, &testSecret)
				if err == nil || err.Error() != "failed to create secret: some error" {
					t.Errorf("Expected error 'failed to create secret: some error', got %v", err)
				}
			},
			expectErr: true,
		},
		{
			name: "Update_Success",
			testFunc: func(t *testing.T) {
				mockRepo.EXPECT().Update(ctx, testSecret).Return(nil)

				updatedSecret, err := service.Update(ctx, &testSecret)
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if updatedSecret.ID != testSecret.ID {
					t.Errorf("Expected secret ID %v, got %v", testSecret.ID, updatedSecret.ID)
				}
			},
			expectErr: false,
		},
		{
			name: "Update_Fail_NotFound",
			testFunc: func(t *testing.T) {
				mockRepo.EXPECT().Update(ctx, testSecret).Return(sql.ErrNoRows)

				_, err := service.Update(ctx, &testSecret)
				if err == nil || err.Error() != "secret not found (id=1)" {
					t.Errorf("Expected error 'secret not found (id=1)', got %v", err)
				}
			},
			expectErr: true,
		},

		{
			name: "Update_Fail",
			testFunc: func(t *testing.T) {
				mockRepo.EXPECT().Update(ctx, testSecret).Return(errors.New("some error"))

				_, err := service.Update(ctx, &testSecret)
				if err == nil || err.Error() != "failed to store secret: some error" {
					t.Errorf("Expected error 'failed to store secret: some error', got %v", err)
				}
			},
			expectErr: true,
		},
		{
			name: "Delete_Success",
			testFunc: func(t *testing.T) {
				mockRepo.EXPECT().Delete(ctx, uint64(1), uint64(1)).Return(nil)

				err := service.Delete(ctx, 1, 1)
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			},
			expectErr: false,
		},
		{
			name: "Delete_Fail",
			testFunc: func(t *testing.T) {
				mockRepo.EXPECT().Delete(ctx, uint64(1), uint64(1)).Return(fmt.Errorf("delete failed"))

				err := service.Delete(ctx, 1, 1)
				if err == nil || err.Error() != "delete failed" {
					t.Errorf("Expected error 'delete failed', got %v", err)
				}
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, tc.testFunc)
	}
}
