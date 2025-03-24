// Package crypto предоставляет функции для безопасного шифрования и дешифрования строк,
// а также для генерации криптографических ключей из паролей и соли. Пакет использует алгоритмы AES-GCM
// для шифрования и scrypt для генерации ключей, обеспечивая высокий уровень безопасности
// при обработке конфиденциальных данных.
//
// Основные возможности пакета:
//
// - DeriveKey: генерация криптографического ключа из пароля и соли.
// - Encrypt: шифрование строки с использованием AES-GCM.
// - Decrypt: расшифровка строки, зашифрованной с помощью Encrypt.
// - Обработка ошибок, связанных с недостаточной длиной зашифрованной строки.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"golang.org/x/crypto/scrypt"
	"io"
)

// ErrCiphertextTooShort указывает, что переданная зашифрованная строка
// (ciphertext) недостаточно длинная для корректной расшифровки.
// Как правило, требуется минимум 12 байт для nonce в AES-GCM.
var ErrCiphertextTooShort = errors.New("ciphertext too short")

// DeriveKey - Генерация ключа из мастер-пароля и соли
func DeriveKey(password, salt string) ([]byte, error) {
	if password == "" {
		return []byte(password), errors.New("empty password")
	}
	return scrypt.Key([]byte(password), []byte(salt), 32768, 8, 1, 32)
}

// Encrypt - Шифрование строки
func Encrypt(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12) // 12 байт - стандартный размер nonce для AES-GCM
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	astc, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := astc.Seal(nil, nonce, []byte(plaintext), nil)
	encrypted := append(nonce, ciphertext...)
	return hex.EncodeToString(encrypted), nil
}

// Decrypt - Расшифровка строки
func Decrypt(encrypted string, key []byte) (string, error) {
	data, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	if len(data) < 12 {
		return "", ErrCiphertextTooShort
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	astc, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := data[:12]
	ciphertext := data[12:]

	plaintext, err := astc.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
