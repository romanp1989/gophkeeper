package domain

import (
	"bytes"
	"fmt"
	"time"
)

// Secret описывает структуру для хранения конфиденциальных данных пользователя
type Secret struct {
	// Уникальный номер записи с конфиденциальной инфомарцией
	ID uint64 `db:"id" json:"id"`
	// Идентификатор пользователя, владельца секрета
	UserID UserID `db:"user_id"`
	// Время создания секрета
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	// Время последнего обновления секрета
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	// Заголовок секрета
	Title string `db:"title" json:"title"`
	// Метаданные, связанные с секретом
	Metadata string `db:"metadata" json:"metadata"`
	// Данные секрета в зашифрованном виде
	Payload []byte `db:"payload" json:"payload"`
	// Тип секрета
	SecretType string `db:"secret_type" json:"secret_type"`

	// Следующие поля не включаются в БД, используются только в методах.
	// Credentials - учетные данные, если SecretType = "credential"
	Credentials *Credentials `db:"-"`
	// Text - текст, если SecretType = "text"
	Text *Text `db:"-"`
	// Blob - бинарные данные, если SecretType = "blob"
	Blob *Blob `db:"-"`
	// Card - данные карты, если SecretType = "card"
	Card *Card `db:"-"`
}

func NewSecret(t SecretType) *Secret {
	return &Secret{SecretType: string(t)}
}

// SecretType тип секрета
type SecretType string

const (
	// CredSecret - секрет, содержащий учетные данные
	CredSecret SecretType = "credential"
	// TextSecret - секрет, содержащий текст
	TextSecret SecretType = "text"
	// BlobSecret - секрет, содержащий бинарные данные
	BlobSecret SecretType = "blob"
	// CardSecret - секрет, содержащий данные карты
	CardSecret SecretType = "card"
	// UnknownSecret - неизвестный тип секрета
	UnknownSecret SecretType = "unknown"
)

// Credentials описывает учетные данные пользователя
type Credentials struct {
	// Login - логин пользователя
	Login string `json:"login"`
	// Password - пароль пользователя
	Password string `json:"password"`
}

// Text описывает текстовую информацию
type Text struct {
	// Content - текстовое содержимое
	Content string `json:"content"`
}

// Blob описывает бинарные данные файла
type Blob struct {
	// FileName - имя файла
	FileName string `json:"file_name"`
	// FileBytes - байты файла
	FileBytes []byte `json:"file_bytes"`
}

// Card описывает данные банковской карты
type Card struct {
	// Number - номер карты
	Number string `json:"number"`
	// ExpYear - год истечения срока действия карты
	ExpYear uint32 `json:"exp_year"`
	// ExpMonth - месяц истечения срока действия карты
	ExpMonth uint32 `json:"exp_month"`
	// CVV - CVV-код карты
	CVV uint32 `json:"cvv"`
}

// ToClipboard форматирует информацию секрета для использования в буфере обмена.
func (s Secret) ToClipboard() string {
	var b bytes.Buffer

	switch SecretType(s.SecretType) {
	case CredSecret:
		b.WriteString(fmt.Sprintf("login: %s\n", s.Credentials.Login))
		b.WriteString(fmt.Sprintf("password: %s", s.Credentials.Password))
	case CardSecret:
		b.WriteString(fmt.Sprintf("Card Number: %s\n", s.Card.Number))
		b.WriteString(fmt.Sprintf("Exp: %d/%d", s.Card.ExpMonth, s.Card.ExpYear))
		b.WriteString(fmt.Sprintf("CVV: %d", s.Card.CVV))
	case TextSecret:
		b.WriteString(fmt.Sprintf("Text: %s\n", s.Text.Content))
	case BlobSecret:
		b.WriteString("File data cannot be moved to clipboard\n")
	}

	return b.String()
}
