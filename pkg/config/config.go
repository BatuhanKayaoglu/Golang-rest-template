package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	MongoDB  MongoConfig
	JWT      JWTConfig
	APM      APMConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type MongoConfig struct {
	URI        string
	Database   string
	Collection string
}

type JWTConfig struct {
	Secret     string
	Expiration int
}

type APMConfig struct {
	ServiceName string
	ServerURL   string
	Environment string
	SecretToken string
	Active      bool
}

var cfg *Config

func Load() *Config {
	if cfg != nil {
		return cfg
	}

	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	jwtExp, _ := strconv.Atoi(getEnv("JWT_EXPIRATION", "24"))
	apmActive, _ := strconv.ParseBool(getEnv("ELASTIC_APM_ACTIVE", "true"))

	cfg = &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8001"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", ""),
			Name:     getEnv("POSTGRES_DB", "postgres"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		MongoDB: MongoConfig{
			URI:        getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database:   getEnv("MONGODB_DATABASE", "logs"),
			Collection: getEnv("MONGODB_COLLECTION", "api_logs"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", ""),
			Expiration: jwtExp,
		},
		APM: APMConfig{
			ServiceName: getEnv("ELASTIC_APM_SERVICE_NAME", "go-rest-api"),
			ServerURL:   getEnv("ELASTIC_APM_SERVER_URL", "http://localhost:8200"),
			Environment: getEnv("ELASTIC_APM_ENVIRONMENT", "development"),
			SecretToken: getEnv("ELASTIC_APM_SECRET_TOKEN", ""),
			Active:      apmActive,
		},
	}

	return cfg
}

func Get() *Config {
	if cfg == nil {
		return Load()
	}
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
