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
	AviationStackAPIKey string
}

func Load() *Config {
	return &Config{
		Port:                getEnv("PORT", "8080"),
		MongoURI:            getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:             getEnv("MONGO_DB", "skyquest"),
		RedisAddr:           getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:       getEnv("REDIS_PASSWORD", ""),
		AviationStackAPIKey: getEnv("AVIATIONSTACK_API_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
