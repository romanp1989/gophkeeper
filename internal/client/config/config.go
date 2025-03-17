// Package config содержит конфигурации для клиента.
package config

import (
	"errors"
	"github.com/spf13/viper"
	"strings"
)

// Config представляет основную конфигурацию клиентского приложения.
type Config struct {
	BuildCommit   string // Информация о сборке (коммит)
	BuildDate     string // Информация о сборке (дата)
	BuildVersion  string // Информация о сборке (версия)
	ServerAddress string // Address определяет адрес сервера.
}

// LoadConfig инициализирует и возвращает новый экземпляр конфигурации.
// Ошибка возвращается, если обязательные конфигурационные параметры не заданы.
func LoadConfig() (*Config, error) {
	viper.SetDefault("address", "127.0.0.1:50051")
	viper.SetEnvPrefix("GOPHKEEPER")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	address := viper.GetString("address")
	if address == "" {
		return nil, errors.New("server address is not set: set GOPHKEEPER_ADDRESS environment variable")
	}

	return &Config{
		ServerAddress: address,
	}, nil
}
