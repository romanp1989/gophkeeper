package config

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name           string
		setupEnv       func()
		expectedError  string
		expectedConfig *Config
	}{
		{
			name: "All_Variables_Set_Correctly",
			setupEnv: func() {
				os.Setenv("GOPHKEEPER_ADDRESS", "127.0.0.1:5000")
			},
			expectedConfig: &Config{
				ServerAddress: "127.0.0.1:5000",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupEnv()
			viper.Reset()

			config, err := LoadConfig()

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Nil(t, config)
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedConfig, config)
			}
		})
	}
}
