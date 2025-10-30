package config

import (
	"os"
	"strconv"
)

type Config struct {
	DB   DBConfig
	NATS NATSConfig
	HTTP HTTPConfig
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type NATSConfig struct {
	ClusterID string
	ClientID  string
	URL       string
}

type HTTPConfig struct {
	Port string
}

func Load() *Config {
	return &Config{
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "postgres"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "order_user"),
			Password: getEnv("DB_PASSWORD", "order_password"),
			Name:     getEnv("DB_NAME", "order_service"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		NATS: NATSConfig{
			ClusterID: getEnv("NATS_CLUSTER_ID", "test-cluster"),
			ClientID:  getEnv("NATS_CLIENT_ID", "order-service"),
			URL:       getEnv("NATS_URL", "nats://nats:4222"),
		},
		HTTP: HTTPConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}