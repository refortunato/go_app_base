package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	internalConf "github.com/refortunato/go_app_base/internal/infra/config"
)

func LoadConfig(path string) (*internalConf.Conf, error) {
	// Carrega o .env se existir (ignora erro se não existir)
	err := godotenv.Load(path + "/.env")
	if err != nil {
		println(path + "/.env file not found")
		println("No .env file found, using environment variables", err.Error())
	}

	cfg := &internalConf.Conf{
		AppName:              getEnv("SERVER_APP_NAME", "go_app_base"),
		WebServerPort:        getEnv("SERVER_APP_WEB_SERVER_PORT", "8080"),
		DBDriver:             getEnv("SERVER_APP_DB_DRIVER", "mysql"),
		DBHost:               getEnv("SERVER_APP_DB_HOST", "localhost"),
		DBPort:               getEnv("SERVER_APP_DB_PORT", "3316"),
		DBUser:               getEnv("SERVER_APP_DB_USER", "root"),
		DBPassword:           getEnv("SERVER_APP_DB_PASSWORD", "root"),
		DBName:               getEnv("SERVER_APP_DB_NAME", "go_app_base"),
		DBMaxOpenConnections: getEnvAsInt("SERVER_APP_DB_MAX_OPEN_CONNECTIONS", 20),
		DBMaxIdleConnections: getEnvAsInt("SERVER_APP_DB_MAX_IDLE_CONNECTIONS", 10),
		DBConnMaxLifetime:    getEnvAsInt("SERVER_APP_DB_CONN_MAX_LIFETIME", 1),
		DBConnMaxIdleTime:    getEnvAsInt("SERVER_APP_DB_CONN_MAX_IDLE_TIME", 10),
		DebugMode:            getEnvAsBool("SERVER_APP_DEBUG_MODE", false),
	}

	println("Configuration loaded:")
	println("DBPort: ", cfg.DBPort)
	println("DBHost: ", cfg.DBHost)
	return cfg, nil
}

// Funções auxiliares para pegar variáveis com valor default
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	if valStr := os.Getenv(key); valStr != "" {
		if val, err := strconv.Atoi(valStr); err == nil {
			return val
		}
	}
	return defaultVal
}

func getEnvAsBool(key string, defaultVal bool) bool {
	if valStr := os.Getenv(key); valStr != "" {
		if val, err := strconv.ParseBool(valStr); err == nil {
			return val
		}
	}
	return defaultVal
}
