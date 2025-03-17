package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/romanp1989/gophkeeper/domain"
	storageErrors "github.com/romanp1989/gophkeeper/pkg/errors"
	"github.com/romanp1989/gophkeeper/tests/mocks"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestUserService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockIUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	ctx := context.Background()
	tests := []struct {
		name      string
		testFunc  func(t *testing.T)
		expectErr bool
	}{
		{
			name: "RegisterUser_Success",
			testFunc: func(t *testing.T) {
				mockRepo.EXPECT().FindByLogin(ctx, "new_user").Return(nil, storageErrors.ErrNotFound).Times(1)
				mockRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(1, nil).Times(1)

				user, err := svc.RegisterUser(ctx, "new_user", "password123")
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if user.Login != "new_user" {
					t.Errorf("Expected login 'newuser', got %v", user.Login)
				}
			},
			expectErr: false,
		},
		{
			name: "RegisterUser_Fail_Create_User",
			testFunc: func(t *testing.T) {
				mockRepo.EXPECT().FindByLogin(ctx, "new_user").Return(nil, storageErrors.ErrNotFound).Times(1)
				mockRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(1, errors.New("some error")).Times(1)

				_, err := svc.RegisterUser(ctx, "new_user", "password123")
				if err == nil || err.Error() != "failed to create user: some error" {
					t.Errorf("Expected error 'failed to create user: some error', got %v", err)
				}
			},
			expectErr: true,
		},
		{
			name: "RegisterUser_Fail_UserExists",
			testFunc: func(t *testing.T) {
				mockRepo.EXPECT().FindByLogin(ctx, "existing_user").Return(&domain.User{Login: "existing_user"}, nil).Times(1)

				_, err := svc.RegisterUser(ctx, "existing_user", "password123")
				if err == nil || err.Error() != "user already exists (existing_user)" {
					t.Errorf("Expected error 'user already exists (existing_user)', got %v", err)
				}
			},
			expectErr: true,
		},
		{
			name: "RegisterUser_Fail_Fetch_User",
			testFunc: func(t *testing.T) {
				mockRepo.EXPECT().FindByLogin(ctx, "existing_user").Return(nil, errors.New("some error")).Times(1)

				_, err := svc.RegisterUser(ctx, "existing_user", "password123")
				if err == nil || err.Error() != "failed to fetch user: some error" {
					t.Errorf("Expected error 'failed to fetch user: some error', got %v", err)
				}
			},
			expectErr: true,
		},
		{
			name: "LoginUser_Success",
			testFunc: func(t *testing.T) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				mockRepo.EXPECT().FindByLogin(ctx, "valid_user").Return(&domain.User{Login: "valid_user", Password: string(hashedPassword)}, nil).Times(1)

				user, err := svc.LoginUser(ctx, "valid_user", "password123")
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if user.Login != "valid_user" {
					t.Errorf("Expected login 'valid_user', got %v", user.Login)
				}
			},
			expectErr: false,
		},
		{
			name: "LoginUser_Fail_WrongPassword",
			testFunc: func(t *testing.T) {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				mockRepo.EXPECT().FindByLogin(ctx, "valid_user").Return(&domain.User{Login: "valid_user", Password: string(hashedPassword)}, nil).Times(1)

				_, err := svc.LoginUser(ctx, "valid_user", "wrong_password")
				if err == nil || !errors.Is(err, ErrBadCredentials) {
					t.Errorf("Expected error 'bad auth credentials', got %v", err)
				}
			},
			expectErr: true,
		},
		{
			name: "LoginUser_Fail_Authenticate",
			testFunc: func(t *testing.T) {
				_, _ = bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
				mockRepo.EXPECT().FindByLogin(ctx, "valid_user").Return(nil, errors.New("some error")).Times(1)

				_, err := svc.LoginUser(ctx, "valid_user", "wrong_password")
				if err == nil || err.Error() != "failed to authenticate user: some error" {
					t.Errorf("Expected error 'failed to authenticate user: some error', got %v", err)
				}
			},
			expectErr: true,
		},
		{
			name: "LoginUser_Fail_UserNotFound",
			testFunc: func(t *testing.T) {
				mockRepo.EXPECT().FindByLogin(ctx, "nonexistent_user").Return(nil, sql.ErrNoRows).Times(1)

				_, err := svc.LoginUser(ctx, "nonexistent_user", "password123")
				if err == nil || !errors.Is(err, ErrBadCredentials) {
					t.Errorf("Expected error 'bad auth credentials', got %v", err)
				}
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, tc.testFunc)
	}
}
