package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config представляет конфигурацию приложения
type Config struct {
	Emulator  EmulatorConfig  `mapstructure:"emulator"`
	Database  DatabaseConfig  `mapstructure:"database"`
	WebSocket WebSocketConfig `mapstructure:"websocket"`
	Bots      BotsConfig      `mapstructure:"bots"`
	Logging   LoggingConfig   `mapstructure:"logging"`
}

// EmulatorConfig конфигурация эмулятора
type EmulatorConfig struct {
	Port  int    `mapstructure:"port"`
	Host  string `mapstructure:"host"`
	Debug bool   `mapstructure:"debug"`
}

// DatabaseConfig конфигурация базы данных
type DatabaseConfig struct {
	URL            string `mapstructure:"url"`
	MaxConnections int    `mapstructure:"max_connections"`
}

// WebSocketConfig конфигурация WebSocket
type WebSocketConfig struct {
	HeartbeatInterval string `mapstructure:"heartbeat_interval"`
	MaxConnections    int    `mapstructure:"max_connections"`
}

// BotsConfig конфигурация ботов
type BotsConfig struct {
	WebhookTimeout string `mapstructure:"webhook_timeout"`
	MaxConnections int    `mapstructure:"max_connections"`
}

// LoggingConfig конфигурация логирования
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	File   string `mapstructure:"file"`
}

// Load загружает конфигурацию из файла и переменных окружения
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	// Переменные окружения
	viper.SetEnvPrefix("EMULATOR")
	viper.AutomaticEnv()

	// Значения по умолчанию
	setDefaults()

	// Чтение конфигурационного файла
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("ошибка чтения конфигурационного файла: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфигурации: %w", err)
	}

	return &config, nil
}

// setDefaults устанавливает значения по умолчанию
func setDefaults() {
	viper.SetDefault("emulator.port", 3001)
	viper.SetDefault("emulator.host", "localhost")
	viper.SetDefault("emulator.debug", true)

	viper.SetDefault("database.url", "sqlite:///data/emulator.db")
	viper.SetDefault("database.max_connections", 10)

	viper.SetDefault("websocket.heartbeat_interval", "30s")
	viper.SetDefault("websocket.max_connections", 1000)

	viper.SetDefault("bots.webhook_timeout", "30s")
	viper.SetDefault("bots.max_connections", 100)

	viper.SetDefault("logging.level", "debug")
	viper.SetDefault("logging.format", "console")
	viper.SetDefault("logging.file", "logs/emulator.log")
}

// GetHeartbeatInterval возвращает интервал heartbeat как Duration
func (c *Config) GetHeartbeatInterval() (time.Duration, error) {
	return time.ParseDuration(c.WebSocket.HeartbeatInterval)
}

// GetWebhookTimeout возвращает timeout webhook как Duration
func (c *Config) GetWebhookTimeout() (time.Duration, error) {
	return time.ParseDuration(c.Bots.WebhookTimeout)
}
