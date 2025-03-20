package user

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

func NewUserRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

// CreateUser Создание нового пользователя
func (r *Repository) CreateUser(ctx context.Context, user *domain.User) (domain.UserID, error) {
	var newUserID domain.UserID

	err := r.db.QueryRowContext(ctx,
		"INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id",
		user.Login,
		user.Password,
	).Scan(&newUserID)

	if err != nil {
		return 0, err
	}

	user.ID = newUserID

	return user.ID, nil
}

// FindByLogin Поиск пользователя по логину
func (r *Repository) FindByLogin(ctx context.Context, login string) (*domain.User, error) {
	u := domain.User{}

	err := r.db.QueryRowContext(ctx,
		"SELECT id, login, password FROM users WHERE login = $1", login,
	).Scan(&u.ID, &u.Login, &u.Password)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storageErrors.ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}
