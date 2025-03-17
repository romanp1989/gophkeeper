package secret

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/romanp1989/gophkeeper/domain"
)

type ISecretRepository interface {
	Create(ctx context.Context, secret *domain.Secret) (*domain.Secret, error)
	GetAllByUserID(ctx context.Context, userID domain.UserID) ([]*domain.Secret, error)
	GetByID(ctx context.Context, id uint64, userID domain.UserID) (*domain.Secret, error)
	Update(ctx context.Context, secret *domain.Secret) (*domain.Secret, error)
	Delete(ctx context.Context, id uint64, userID domain.UserID) error
}

type Service struct {
	repository ISecretRepository
}

func NewSecretService(repository ISecretRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Get(ctx context.Context, secretID uint64, userID domain.UserID) (*domain.Secret, error) {
	secret, err := s.repository.GetByID(ctx, secretID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("secret not found id %d", secretID)
		}
		return nil, err
	}

	return secret, nil
}

func (s *Service) GetUserSecrets(ctx context.Context, userID domain.UserID) ([]*domain.Secret, error) {
	secrets, err := s.repository.GetAllByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user's secrets not found id %d", userID)
		}
		return nil, err
	}

	return secrets, nil
}

func (s *Service) Add(ctx context.Context, secret *domain.Secret) (*domain.Secret, error) {
	var err error

	secret, err = s.repository.Create(ctx, secret)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret: %w", err)
	}
	return secret, nil
}

func (s *Service) Update(ctx context.Context, secret *domain.Secret) (*domain.Secret, error) {
	var err error

	secret, err = s.repository.Update(ctx, secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("secret not found id %d", secret.ID)
		}

		return nil, fmt.Errorf("failed to update secret: %d", err)
	}

	return secret, nil
}

func (s *Service) Delete(ctx context.Context, secretID uint64, userID domain.UserID) error {
	err := s.repository.Delete(ctx, secretID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("secret not found id %d", secretID)
		}
		return fmt.Errorf("failed to delete secret: %w", err)
	}

	return nil
}
