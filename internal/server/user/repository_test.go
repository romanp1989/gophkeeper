package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/romanp1989/gophkeeper/domain"
	storageErrors "github.com/romanp1989/gophkeeper/pkg/errors"
	"testing"
	"time"
)

func TestUserRepository(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name      string
		testFunc  func(t *testing.T, repo IUserRepository, mock sqlmock.Sqlmock)
		expectErr bool
	}{
		{
			name: "CreateUser_Success",
			testFunc: func(t *testing.T, repo IUserRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO users \(login, password\) VALUES \(\$1, \$2\) RETURNING id`).
					WithArgs("new_user", "hashed_password").
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				user := domain.User{
					Login:    "new_user",
					Password: "hashed_password",
				}
				id, err := repo.CreateUser(ctx, &user)
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if id != 1 {
					t.Errorf("Expected ID 1, got %d", id)
				}
			},
			expectErr: false,
		},
		{
			name: "CreateUser_Fail_DatabaseError",
			testFunc: func(t *testing.T, repo IUserRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO users \(login, password\) VALUES \(\$1, \$2\) RETURNING id`).
					WithArgs("new_user", "hashed_password").
					WillReturnError(fmt.Errorf("database error"))

				user := domain.User{
					Login:    "new_user",
					Password: "hashed_password",
				}
				_, err := repo.CreateUser(ctx, &user)
				if err == nil || err.Error() != "database error" {
					t.Errorf("Expected error 'database error', got %v", err)
				}
			},
			expectErr: true,
		},
		{
			name: "FindByLogin_Success",
			testFunc: func(t *testing.T, repo IUserRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, login, created_at, password FROM users WHERE id = \$1`).
					WithArgs("existing_user").
					WillReturnRows(sqlmock.NewRows([]string{"id", "login", "created_at", "password"}).
						AddRow(1, "existing_user", time.Now(), "hashed_password"))

				user, err := repo.FindByLogin(ctx, "existing_user")
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if user.Login != "existing_user" {
					t.Errorf("Expected login 'existing_user', got %v", user.Login)
				}
			},
			expectErr: false,
		},
		{
			name: "FindByLogin_Fail_NotFound",
			testFunc: func(t *testing.T, repo IUserRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, login, created_at, password FROM users WHERE id = \$1`).
					WithArgs("not_existing_user").
					WillReturnError(sql.ErrNoRows)

				_, err := repo.FindByLogin(ctx, "not_existing_user")
				if err == nil || !errors.Is(err, storageErrors.ErrNotFound) {
					t.Errorf("Expected error 'ErrNotFound', got %v", err)
				}
			},
			expectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create sqlmock: %v", err)
			}
			defer db.Close()

			repo := NewUserRepository(db)

			tc.testFunc(t, repo, mock)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unmet SQL expectations: %v", err)
			}
		})
	}
}
