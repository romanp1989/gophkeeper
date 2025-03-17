package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/romanp1989/gophkeeper/certs"
	"github.com/romanp1989/gophkeeper/domain"
	"github.com/romanp1989/gophkeeper/internal/client/config"
	"github.com/romanp1989/gophkeeper/internal/client/grpc/interceptors"
	"github.com/romanp1989/gophkeeper/pkg/converter"
	"github.com/romanp1989/gophkeeper/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
	"math/rand/v2"
	"sync"
	"time"
)

type ClientGRPCInterface interface {
	Login(ctx context.Context, login, password string) (string, error)
	Register(ctx context.Context, login, password string) (string, error)
	LoadSecrets(ctx context.Context) ([]*domain.Secret, error)
	LoadSecret(ctx context.Context, ID uint64) (*domain.Secret, error)
	SaveSecret(ctx context.Context, secret *domain.Secret) error
	DeleteSecret(ctx context.Context, id uint64) error
	SetToken(token string)
	GetToken() string
	SetPassword(password string)
	GetPassword() string
}

type (
	// ClientGRPC управляет соединением с gRPC сервером и реализует методы для работы с серверными ресурсами.
	ClientGRPC struct {
		config        *config.Config
		UsersClient   proto.UsersClient
		SecretsClient proto.SecretsClient
		accessToken   string
		password      string
		clientID      uint64
		previews      sync.Map
	}
	// ReloadSecretList метка для обработчика
	ReloadSecretList struct{}
)

// NewClientGRPC создаёт новый экземпляр ClientGRPC с предварительной настройкой подключения к серверу.
func NewClientGRPC(cfg *config.Config) (ClientGRPCInterface, error) {
	var opts []grpc.DialOption

	newClient := ClientGRPC{
		config:   cfg,
		clientID: uint64(rand.IntN(math.MaxInt32)),
	}

	opts = append(
		opts,
		grpc.WithChainUnaryInterceptor(
			interceptors.Timeout(time.Second*5),
			interceptors.AddAuth(&newClient.accessToken, uint32(newClient.clientID)),
		),
	)

	opts = append(
		opts,
	)

	tlsCredential, err := loadTLSConfig("ca-cert.pem", "client-cert.pem", "client-key.pem")
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS config: %w", err)
	}

	opts = append(opts, grpc.WithTransportCredentials(tlsCredential))

	c, err := grpc.NewClient(cfg.ServerAddress, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}

	newClient.UsersClient = proto.NewUsersClient(c)
	newClient.SecretsClient = proto.NewSecretsClient(c)

	return &newClient, nil
}

// Login авторизует пользователя на сервере и получает токен доступа.
func (c *ClientGRPC) Login(ctx context.Context, login string, password string) (string, error) {
	req := &proto.LoginRequest{
		Login:    login,
		Password: password,
	}

	response, err := c.UsersClient.Login(ctx, req)
	if err != nil {
		return "", parseError(err)
	}

	c.accessToken = response.AccessToken

	return response.AccessToken, nil
}

// Register регистрирует нового пользователя и получает токен доступа.
func (c *ClientGRPC) Register(ctx context.Context, login string, password string) (string, error) {
	req := &proto.RegisterRequest{
		Login:    login,
		Password: password,
	}

	response, err := c.UsersClient.Register(ctx, req)
	if err != nil {
		return "", parseError(err)
	}

	c.accessToken = response.AccessToken

	return response.AccessToken, nil
}

// LoadSecrets загружает список секретов пользователя.
func (c *ClientGRPC) LoadSecrets(ctx context.Context) ([]*domain.Secret, error) {
	request := emptypb.Empty{}

	response, err := c.SecretsClient.GetUserSecrets(ctx, &request)
	if err != nil {
		return nil, parseError(err)
	}

	secrets := converter.ProtoToSecrets(response.Secrets)
	return secrets, nil
}

// LoadSecret загружает информацию о конкретном секрете.
func (c *ClientGRPC) LoadSecret(_ context.Context, ID uint64) (*domain.Secret, error) {
	request := &proto.GetUserSecretRequest{
		Id: ID,
	}

	response, err := c.SecretsClient.GetUserSecret(context.Background(), request)
	if err != nil {
		return nil, parseError(err)
	}

	secret := converter.ProtoToSecret(response.Secret)

	return secret, nil
}

// SaveSecret сохраняет или обновляет секрет пользователя на сервере.
func (c *ClientGRPC) SaveSecret(ctx context.Context, secret *domain.Secret) error {
	sec := &proto.Secret{
		Title:      secret.Title,
		Metadata:   secret.Metadata,
		SecretType: converter.TypeToProto(secret.SecretType),
		Payload:    secret.Payload,
		CreatedAt:  timestamppb.New(secret.CreatedAt),
		UpdatedAt:  timestamppb.New(secret.UpdatedAt),
	}

	if secret.ID > 0 {
		sec.Id = secret.ID
	}

	request := &proto.SaveUserSecretRequest{Secret: sec}
	_, err := c.SecretsClient.SaveUserSecret(ctx, request)

	return parseError(err)
}

// DeleteSecret удаляет секрет пользователя.
func (c *ClientGRPC) DeleteSecret(ctx context.Context, id uint64) error {
	request := &proto.DeleteUserSecretRequest{Id: id}
	_, err := c.SecretsClient.DeleteUserSecret(ctx, request)

	return parseError(err)
}

// SetToken устанавливает текущий токен доступа клиента.
func (c *ClientGRPC) SetToken(token string) {
	c.accessToken = token
}

// GetToken возвращает текущий токен доступа клиента.
func (c *ClientGRPC) GetToken() string {
	return c.accessToken
}

// SetPassword устанавливает текущий пароль клиента.
func (c *ClientGRPC) SetPassword(password string) {
	c.password = password
}

// GetPassword возвращает текущий пароль клиента.
func (c *ClientGRPC) GetPassword() string {
	return c.password
}

// loadTLSConfig загружает TLS конфигурацию для подключения к серверу.
func loadTLSConfig(caCertFile, clientCertFile, clientKeyFile string) (credentials.TransportCredentials, error) {
	caPem, err := certs.Cert.ReadFile(caCertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA cert: %w", err)
	}

	clientCertPEM, err := certs.Cert.ReadFile(clientCertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read client cert: %w", err)
	}

	clientKeyPEM, err := certs.Cert.ReadFile(clientKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read client key: %w", err)
	}

	clientCert, err := tls.X509KeyPair(clientCertPEM, clientKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to load x509 key pair: %w", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPem) {
		return nil, fmt.Errorf("failed to append CA cert to cert pool: %w", err)
	}

	tlcConfiguration := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}

	return credentials.NewTLS(tlcConfiguration), nil
}

// parseError анализирует ошибки от gRPC вызовов и конвертирует их в более понятный формат.
func parseError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.Unavailable:
		return errors.New("server unavailable")
	case codes.Unauthenticated:
		return errors.New("failed to authenticate")
	case codes.AlreadyExists:
		return errors.New("user already exists")
	default:
		return err
	}
}
