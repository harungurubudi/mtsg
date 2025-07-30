package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		envVars map[string]string
		want    *Config
		wantErr bool
	}{
		{
			name:    "load with defaults",
			envVars: map[string]string{},
			want: &Config{
				Server: ServerConfig{
					Port:         "8080",
					ReadTimeout:  30 * time.Second,
					WriteTimeout: 30 * time.Second,
					IdleTimeout:  60 * time.Second,
					Environment:  "development",
				},
				Redis: RedisConfig{
					Host:     "localhost",
					Port:     6379,
					Password: "",
					DB:       0,
					PoolSize: 10,
				},
			},
			wantErr: false,
		},
		{
			name: "load with environment variables",
			envVars: map[string]string{
				"MTSG_SERVER_PORT":        "9090",
				"MTSG_SERVER_ENVIRONMENT": "production",
				"MTSG_REDIS_HOST":         "redis.example.com",
				"MTSG_REDIS_PORT":         "6380",
			},
			want: &Config{
				Server: ServerConfig{
					Port:         "9090",
					ReadTimeout:  30 * time.Second,
					WriteTimeout: 30 * time.Second,
					IdleTimeout:  60 * time.Second,
					Environment:  "production",
				},
				Redis: RedisConfig{
					Host:     "redis.example.com",
					Port:     6380,
					Password: "",
					DB:       0,
					PoolSize: 10,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}
			defer func() {
				// Cleanup environment variables
				for key := range tt.envVars {
					os.Unsetenv(key)
				}
			}()

			got, err := Load()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Server: ServerConfig{
					Port:         "8080",
					ReadTimeout:  30 * time.Second,
					WriteTimeout: 30 * time.Second,
					IdleTimeout:  60 * time.Second,
					Environment:  "development",
				},
				Redis: RedisConfig{
					Host:     "localhost",
					Port:     6379,
					Password: "",
					DB:       0,
					PoolSize: 10,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestServerConfig(t *testing.T) {
	config := &ServerConfig{
		Port:         "8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
		Environment:  "development",
	}

	assert.Equal(t, "8080", config.Port)
	assert.Equal(t, 30*time.Second, config.ReadTimeout)
	assert.Equal(t, 30*time.Second, config.WriteTimeout)
	assert.Equal(t, 60*time.Second, config.IdleTimeout)
	assert.Equal(t, "development", config.Environment)
}

func TestRedisConfig(t *testing.T) {
	config := &RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "secret",
		DB:       1,
		PoolSize: 20,
	}

	assert.Equal(t, "localhost", config.Host)
	assert.Equal(t, 6379, config.Port)
	assert.Equal(t, "secret", config.Password)
	assert.Equal(t, 1, config.DB)
	assert.Equal(t, 20, config.PoolSize)
}

func TestLoadWithInvalidEnvironmentVariable(t *testing.T) {
	// Set an invalid port (non-numeric)
	os.Setenv("MTSG_REDIS_PORT", "invalid")
	defer os.Unsetenv("MTSG_REDIS_PORT")

	_, err := Load()
	assert.Error(t, err)
}
