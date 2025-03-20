package secret

import (
	"context"
	"database/sql"
	"errors"
	"github.com/romanp1989/gophkeeper/domain"
	storageErrors "github.com/romanp1989/gophkeeper/pkg/errors"
)

type Repository struct {
	db *sql.DB
}

func NewSecretRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, secret *domain.Secret) (*domain.Secret, error) {
	var insertedID uint64

	query := `INSERT INTO secrets (user_id, title, metadata, secret_type, payload) 
			VALUES ($1, $2, $3, $4, $5) 
			RETURNING secret_id`

	result := r.db.QueryRowContext(ctx, query, secret.UserID, secret.Title, secret.Metadata, secret.SecretType, secret.Payload)
	err := result.Scan(&insertedID)
	if err != nil {
		return nil, err
	}

	secret.ID = insertedID

	return secret, nil
}

func (r *Repository) GetAllByUserID(ctx context.Context, userID domain.UserID) ([]*domain.Secret, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM secrets WHERE user_id = ? ORDER BY updated_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	secrets := make([]*domain.Secret, 0)

	for rows.Next() {
		var secret domain.Secret

		err := rows.Scan(&secret.ID, &secret.Title, &secret.Metadata, &secret.Payload)
		if err != nil {
			return nil, err
		}

		secrets = append(secrets, &secret)
	}

	return secrets, nil
}

func (r *Repository) GetByID(ctx context.Context, id uint64, userID domain.UserID) (*domain.Secret, error) {
	var secret domain.Secret

	query := `SELECT * FROM secrets WHERE id = $1 AND user_id = $2`

	err := r.db.QueryRowContext(ctx, query, id, userID).Scan(secret.ID, &secret.Title, &secret.Metadata, &secret.Payload)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storageErrors.ErrNotFound
		}
		return nil, err
	}

	return &secret, nil
}

// Update обновление конфиденциальных данных
func (r *Repository) Update(ctx context.Context, secret *domain.Secret) (*domain.Secret, error) {
	var secretID uint64
	err := r.db.QueryRowContext(ctx, "SELECT id FROM secrets WHERE user_id = $1 AND id = $2 FOR UPDATE", secret.UserID, secret.ID).Scan(secretID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storageErrors.ErrNotFound
		}
		return nil, err
	}

	query := `UPDATE secrets SET title = $1, metadata = $2, payload = $3, updated_at = $4 WHERE id = $5`
	_, err = r.db.ExecContext(ctx, query, secret.Title, secret.Metadata, secret.Payload, secret.UpdatedAt, secret.ID)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

// Delete удаление конфиденциальных данных
func (r *Repository) Delete(ctx context.Context, id uint64, userID domain.UserID) error {
	var secretID uint64
	err := r.db.QueryRowContext(ctx, "SELECT id FROM secrets WHERE user_id = $1 AND id = $2 FOR UPDATE", userID, id).Scan(secretID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return storageErrors.ErrNotFound
		}
		return err
	}

	query := `DELETE FROM secrets WHERE id = $1 AND user_id = $2`
	_, err = r.db.ExecContext(ctx, query, id, secretID)

	return err
}
