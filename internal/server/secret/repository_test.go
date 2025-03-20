package secret

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/romanp1989/gophkeeper/domain"
	"testing"
	"time"
)

func TestSecretRepository(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		testFunc  func(t *testing.T, repo SecretRepository, mock sqlmock.Sqlmock)
		expectErr bool
	}{
		{
			name: "GetByID_Success",
			testFunc: func(t *testing.T, repo SecretRepository, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "title", "metadata", "secret_type", "payload", "created_at", "updated_at"}).
					AddRow(1, 1, "Test Secret", "Metadata", "text", []byte("payload"), time.Now(), time.Now())

				mock.ExpectQuery(`SELECT \* FROM secrets WHERE id = \$1 AND user_id = \$2`).
					WithArgs(1, 1).
					WillReturnRows(rows)

				secret, err := repo.GetByID(ctx, 1, 1)
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if secret.ID != 1 || secret.Title != "Test Secret" {
					t.Errorf("Unexpected secret data: %+v", secret)
				}
			},
			expectErr: false,
		},
		{
			name: "GetByID_Fail_NotFound",
			testFunc: func(t *testing.T, repo SecretRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM secrets WHERE id = \$1 AND user_id = \$2`).
					WithArgs(1, 1).
					WillReturnError(sql.ErrNoRows)

				_, err := repo.GetByID(ctx, 1, 1)
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			},
			expectErr: true,
		},
		{
			name: "GetAllByUserID_Success",
			testFunc: func(t *testing.T, repo SecretRepository, mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "title", "metadata", "secret_type", "payload", "created_at", "updated_at"}).
					AddRow(1, 1, "Secret 1", "Metadata 1", "text", []byte("payload1"), time.Now(), time.Now()).
					AddRow(2, 1, "Secret 2", "Metadata 2", "text", []byte("payload2"), time.Now(), time.Now())

				mock.ExpectQuery(`SELECT \* FROM secrets WHERE user_id = \$1 ORDER BY updated_at DESC`).
					WithArgs(1).
					WillReturnRows(rows)

				secrets, err := repo.GetAllByUserID(ctx, 1)
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if len(secrets) != 2 {
					t.Errorf("Expected 2 secrets, got %d", len(secrets))
				}
			},
			expectErr: false,
		},
		{
			name: "GetAllByUserID_Fail_QueryError",
			testFunc: func(t *testing.T, repo SecretRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM secrets WHERE user_id = \$1 ORDER BY updated_at DESC`).
					WithArgs(1).
					WillReturnError(fmt.Errorf("database error"))

				_, err := repo.GetAllByUserID(ctx, 1)
				if err == nil || err.Error() != "database error" {
					t.Errorf("Expected error 'database error', got %v", err)
				}
			},
			expectErr: true,
		},
		{
			name: "Create_Success",
			testFunc: func(t *testing.T, repo SecretRepository, mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO secrets \(user_id, title, metadata, secret_type, payload\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING id`).
					WithArgs(1, "Test Secret", "Metadata", "text", []byte("payload")).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				secret := &domain.Secret{
					UserID:     1,
					Title:      "Test Secret",
					Metadata:   "Metadata",
					SecretType: "text",
					Payload:    []byte("payload"),
				}
				insertedSecret, err := repo.Create(ctx, secret)
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if insertedSecret.ID != 1 {
					t.Errorf("Expected ID 1, got %v", insertedSecret.ID)
				}
			},
			expectErr: false,
		},
		{
			name: "Update_Success",
			testFunc: func(t *testing.T, repo SecretRepository, mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT 1 FROM secrets WHERE id = \$1 FOR UPDATE`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))
				mock.ExpectExec(`UPDATE secrets SET updated_at = \$1, title = \$2, metadata = \$3, secret_type = \$4, payload = \$5 WHERE id = \$6`).
					WithArgs(sqlmock.AnyArg(), "Updated Title", "Updated Metadata", "text", []byte("updated payload"), 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()

				secret := &domain.Secret{
					ID:         1,
					Title:      "Updated Title",
					Metadata:   "Updated Metadata",
					SecretType: "text",
					Payload:    []byte("updated payload"),
					UpdatedAt:  time.Now(),
				}
				secret, err := repo.Update(ctx, secret)
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			},
			expectErr: false,
		},
		{
			name: "Delete_Success",
			testFunc: func(t *testing.T, repo SecretRepository, mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM secrets WHERE id = \$1 AND user_id = \$2`).
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))

				err := repo.Delete(ctx, 1, 1)
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			},
			expectErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("Failed to create sqlmock: %v", err)
			}
			defer db.Close()

			repo := NewSecretRepository(db)

			tc.testFunc(t, repo, mock)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unmet SQL expectations: %v", err)
			}
		})
	}
}
