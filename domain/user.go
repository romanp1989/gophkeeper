package domain

import "time"

type UserID uint64

// User описывает структуру пользователя
type User struct {
	// Уникальный идентификатор пользователя
	ID UserID `json:"id"`
	// Логин пользователя, используемый для входа в систему
	Login string `json:"login"`
	// Пароль пользователя. Это поле не включается в JSON представление
	Password string `json:"password"`
	// Временная метка создания аккаунта пользователя
	CreatedAt time.Time `json:"created_at"`
	// Временная метка последнего обновления данных аккаунта пользователя
	UpdatedAt time.Time `json:"updated_at"`
}
