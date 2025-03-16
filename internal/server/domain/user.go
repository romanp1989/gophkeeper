package domain

import "time"

type UserID uint64

// User описывает структуру пользователя
type User struct {
	// ID уникальный идентификатор пользователя
	ID UserID `json:"id"`
	// Login содержит логин пользователя, используемый для входа в систему
	Login string `json:"login"`
	// Password содержит пароль пользователя. Это поле не включается в JSON представление
	Password string `json:"password"`
	// CreatedAt содержит временную метку создания аккаунта пользователя
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt содержит временную метку последнего обновления данных аккаунта пользователя
	UpdatedAt time.Time `json:"updated_at"`
}
