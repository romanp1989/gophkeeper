// Package storage предоставляет интерфейсы и реализации для управления хранилищем секретов.
package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/romanp1989/gophkeeper/domain"
	"github.com/romanp1989/gophkeeper/internal/client/crypto"
	"github.com/romanp1989/gophkeeper/internal/client/grpc"
)

// Storage описывает интерфейс для базовых операций с хранилищем секретов.
type Storage interface {
	Get(ctx context.Context, id uint64) (*domain.Secret, error)
	GetAll(ctx context.Context) ([]*domain.Secret, error)
	Create(ctx context.Context, secret *domain.Secret) error
	Update(ctx context.Context, secret *domain.Secret) error
	Delete(ctx context.Context, id uint64) error
	String() string
}

// RemoteStorage реализует хранилище секретов, используя удаленный сервис через gRPC.
type RemoteStorage struct {
	client    grpc.ClientGRPCInterface
	deriveKey []byte
}

// NewRemoteStorage создает новый экземпляр RemoteStorage с предварительно вычисленным ключом шифрования.
func NewRemoteStorage(client grpc.ClientGRPCInterface) (*RemoteStorage, error) {
	deriveKey, err := crypto.DeriveKey(client.GetPassword(), "")
	if err != nil {
		return nil, err
	}
	return &RemoteStorage{
		client:    client,
		deriveKey: deriveKey,
	}, nil
}

// Get извлекает секрет по его идентификатору, расшифровывает его и возвращает.
func (store *RemoteStorage) Get(_ context.Context, id uint64) (*domain.Secret, error) {
	secret, err := store.client.LoadSecret(context.Background(), id)
	if err != nil {
		return nil, err
	}

	err = store.decryptPayload(secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

// GetAll извлекает все секреты пользователя, расшифровывает их и возвращает.
func (store *RemoteStorage) GetAll(_ context.Context) ([]*domain.Secret, error) {
	secrets, err := store.client.LoadSecrets(context.Background())
	if err != nil {
		return nil, err
	}

	for _, s := range secrets {
		err = store.decryptPayload(s)
		if err != nil {
			return nil, err
		}
	}

	return secrets, nil
}

// Create создает новый секрет в хранилище, предварительно зашифровав его.
func (store *RemoteStorage) Create(_ context.Context, secret *domain.Secret) (err error) {
	err = store.encryptPayload(secret)
	if err != nil {
		return
	}

	err = store.client.SaveSecret(context.Background(), secret)
	return err
}

// Update обновляет существующий секрет, предварительно зашифровав его.
func (store *RemoteStorage) Update(_ context.Context, secret *domain.Secret) (err error) {
	err = store.encryptPayload(secret)
	if err != nil {
		return
	}

	err = store.client.SaveSecret(context.Background(), secret)
	return err
}

// Delete удаляет секрет по его идентификатору.
func (store *RemoteStorage) Delete(_ context.Context, id uint64) (err error) {
	err = store.client.DeleteSecret(context.Background(), id)
	return err
}

func (store *RemoteStorage) String() string {
	return "remote storage"
}

// encryptPayload шифрует данные секрета перед сохранением.
func (store *RemoteStorage) encryptPayload(secret *domain.Secret) (err error) {
	data, err := marshalSecret(secret)
	if err != nil {
		return fmt.Errorf("encryptPayload(): error serializing data: %w", err)
	}

	encryptedData, err := crypto.Encrypt(string(data), store.deriveKey)
	if err != nil {
		return fmt.Errorf("encryptPayload(): error encrypting Data: %w", err)
	}

	secret.Payload = []byte(encryptedData)
	return nil
}

// decryptPayload расшифровывает данные секрета после извлечения.
func (store *RemoteStorage) decryptPayload(secret *domain.Secret) (err error) {
	decryptedData, err := crypto.Decrypt(string(secret.Payload), store.deriveKey)
	if err != nil {
		return fmt.Errorf("decryptPayload: failed to decrypt data: %w", err)

	}

	err = unmarshalSecret(secret, []byte(decryptedData))
	if err != nil {
		return fmt.Errorf("decryptPayload: failed to unmarshal data: %w", err)
	}

	return nil
}

// marshalSecret кодирует данные секрета в JSON.
func marshalSecret(secret *domain.Secret) ([]byte, error) {
	var (
		data []byte
		err  error
	)

	switch domain.SecretType(secret.SecretType) {
	case domain.CredSecret:
		data, err = json.Marshal(secret.Credentials)
	case domain.TextSecret:
		data, err = json.Marshal(secret.Text)
	case domain.CardSecret:
		data, err = json.Marshal(secret.Card)
	case domain.BlobSecret:
		data, err = json.Marshal(secret.Blob)
	}

	return data, err
}

// unmarshalSecret кодирует данные секрета из JSON.
func unmarshalSecret(secret *domain.Secret, data []byte) error {
	var err error

	switch domain.SecretType(secret.SecretType) {
	case domain.CredSecret:
		err = json.Unmarshal(data, &secret.Credentials)
	case domain.TextSecret:
		err = json.Unmarshal(data, &secret.Text)
	case domain.CardSecret:
		err = json.Unmarshal(data, &secret.Card)
	case domain.BlobSecret:
		err = json.Unmarshal(data, &secret.Blob)
	}

	return err
}
