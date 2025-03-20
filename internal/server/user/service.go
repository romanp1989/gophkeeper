package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/romanp1989/gophkeeper/domain"
	storageErrors "github.com/romanp1989/gophkeeper/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// ErrBadCredentials определяет ошибку, возникающую при неверных учетных данных для аутентификации.
var ErrBadCredentials = errors.New("bad token credentials")

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) (domain.UserID, error)
	FindByLogin(ctx context.Context, login string) (*domain.User, error)
}

type Service struct {
	userRepository UserRepository
}

// NewUserService создает новый экземпляр UserService с заданным репозиторием
func NewUserService(userRepository UserRepository) *Service {
	return &Service{userRepository: userRepository}
}

// RegisterUser метод регистрации пользователя
func (s *Service) RegisterUser(ctx context.Context, login string, password string) (*domain.User, error) {
	var newUser *domain.User

	user, err := s.userRepository.FindByLogin(ctx, login)
	if err != nil && !errors.Is(err, storageErrors.ErrNotFound) {
		return newUser, fmt.Errorf("failed to find user by login: %w", err)
	}

	if user != nil {
		return newUser, fmt.Errorf("user already exists")
	}

	hashPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return newUser, fmt.Errorf("failed to hash password: %w", err)
	}

	newUser = &domain.User{
		Login:     login,
		Password:  string(hashPwd),
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	var userID domain.UserID
	userID, err = s.userRepository.CreateUser(ctx, newUser)
	if err != nil {
		return newUser, fmt.Errorf("failed to create user: %w", err)
	}

	newUser.ID = userID

	return newUser, nil
}

// Login метод авторизации пользователя
func (s *Service) LoginUser(ctx context.Context, login string, password string) (*domain.User, error) {
	user, err := s.userRepository.FindByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, storageErrors.ErrNotFound) {
			return nil, ErrBadCredentials
		}

		return nil, fmt.Errorf("failed to authenticate user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrBadCredentials
	}

	return user, nil
}
