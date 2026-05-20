package config

import (
	"os"
	"strings"
)

type Config struct {
	Port            string
	UploadRoot      string
	MinioEndpoint   string
	MinioAccessKey  string
	MinioSecretKey  string
	MinioBucket     string
	MinioUseSSL     bool
	MinioPublicBase string
}

func Load() Config {
	return Config{
		Port:            getEnv("UPLOAD_SERVICE_PORT", "9091"),
		UploadRoot:      getEnv("UPLOAD_ROOT", "./uploads"),
		MinioEndpoint:   getEnv("MINIO_ENDPOINT", "127.0.0.1:9000"),
		MinioAccessKey:  getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinioSecretKey:  getEnv("MINIO_SECRET_KEY", "minioadmin"),
		MinioBucket:     getEnv("MINIO_BUCKET", "gin-blog-upload"),
		MinioUseSSL:     getEnvBool("MINIO_USE_SSL", false),
		MinioPublicBase: getEnv("MINIO_PUBLIC_BASE", ""),
	}
}

func getEnv(key, defaultVal string) string {
	if val := strings.TrimSpace(os.Getenv(key)); val != "" {
		return val
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	val := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if val == "" {
		return defaultVal
	}
	return val == "1" || val == "true" || val == "yes"
}

