package config

import (
	"os"
)

type Config struct {
	Port                string
	MongoURI            string
	MongoDB             string
	RedisAddr           string
	RedisPassword       string
	RedisTLS            bool
	AviationStackAPIKey string
}

func Load() *Config {
	return &Config{
		Port:                getEnv("PORT", "8080"),
		MongoURI:            getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:             getEnv("MONGO_DB", "skyquest"),
		RedisAddr:           getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:       getEnv("REDIS_PASSWORD", ""),
		RedisTLS:            getEnv("REDIS_TLS", "false") == "true",
		AviationStackAPIKey: getEnv("AVIATIONSTACK_API_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
