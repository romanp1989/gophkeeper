package config

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/romanp1989/gophkeeper/certs"
	"github.com/romanp1989/gophkeeper/internal/server/db"
	"github.com/romanp1989/gophkeeper/internal/server/token"
	"github.com/spf13/viper"
	"google.golang.org/grpc/credentials"
	"strings"
	"time"
)

// Config представляет основную конфигурацию клиентского приложения.
type Config struct {
	Address string        // Address определяет адрес сервера.
	Db      *db.Config    // Db конфиг подключения к PostgreSQL.
	Token   *token.Config // Token конфиг JWT токена для авторизации
}

// NewConfig инициализирует и возвращает новый экземпляр конфигурации.
// Ошибка возвращается, если обязательные конфигурационные параметры не заданы.
func NewConfig() (*Config, error) {
	viper.SetEnvPrefix("GOPHKEEPER")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	address := viper.GetString("address")
	if address == "" {
		return nil, errors.New("server address is not set: set GOPHKEEPER_ADDRESS environment variable")
	}

	dsn := viper.GetString("postgres-dsn")
	if dsn == "" {
		return nil, errors.New("PostgreSQL DSN is not set: set GOPHKEEPER_POSTGRES_DSN environment variable")
	}
	dbConfig := &db.Config{
		Dsn:             dsn,
		MaxIdleConns:    1,
		MaxOpenConns:    10,
		MaxLifetimeConn: time.Minute * 1,
	}

	secretKey := viper.GetString("secret-key")
	if secretKey == "" {
		return nil, errors.New("secret key for signing JWT is not set: set GOPHKEEPER_SECRET_KEY environment variable")
	}

	tokenConfig := &token.Config{
		Secret: secretKey,
		Name:   "Authorization",
		Expire: 24 * time.Hour,
	}

	return &Config{
		Address: address,
		Db:      dbConfig,
		Token:   tokenConfig,
	}, nil
}

// LoadTLSConfig загружает TLS конфигурацию для сервера из указанных сертификата и ключа
func (c *Config) LoadTLSConfig(caCertFile, serverCertFile, serverKeyFile string) (credentials.TransportCredentials, error) {
	caPem, err := certs.Cert.ReadFile(caCertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA cert: %w", err)
	}

	serverCertPEM, err := certs.Cert.ReadFile(serverCertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read server cert: %w", err)
	}

	serverKeyPEM, err := certs.Cert.ReadFile(serverKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read server key: %w", err)
	}

	serverCert, err := tls.X509KeyPair(serverCertPEM, serverKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to load x509 key pair: %w", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPem) {
		return nil, fmt.Errorf("failed to append CA cert to cert pool: %w", err)
	}

	tlsCfg := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(tlsCfg), nil
}
