package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Conf holds all application configuration
type Conf struct {
	AppName              string `mapstructure:"SERVER_APP_NAME"`
	ImageName            string `mapstructure:"SERVER_APP_IMAGE_NAME"`
	ImageVersion         string `mapstructure:"SERVER_APP_IMAGE_VERSION"`
	Environment          string `mapstructure:"SERVER_APP_ENVIRONMENT"`
	DBDriver             string `mapstructure:"SERVER_APP_DB_DRIVER"`
	DBHost               string `mapstructure:"SERVER_APP_DB_HOST"`
	DBPort               string `mapstructure:"SERVER_APP_DB_PORT"`
	DBUser               string `mapstructure:"SERVER_APP_DB_USER"`
	DBPassword           string `mapstructure:"SERVER_APP_DB_PASSWORD"`
	DBName               string `mapstructure:"SERVER_APP_DB_NAME"`
	DBMaxOpenConnections int    `mapstructure:"SERVER_APP_DB_MAX_OPEN_CONNECTIONS"`
	DBMaxIdleConnections int    `mapstructure:"SERVER_APP_DB_MAX_IDLE_CONNECTIONS"`
	DBConnMaxLifetime    int    `mapstructure:"SERVER_APP_DB_CONN_MAX_LIFETIME"`  // in hours
	DBConnMaxIdleTime    int    `mapstructure:"SERVER_APP_DB_CONN_MAX_IDLE_TIME"` // in minutes
	WebServerPort        string `mapstructure:"SERVER_APP_WEB_SERVER_PORT"`
	DebugMode            bool   `mapstructure:"SERVER_APP_DEBUG_MODE"`
	SwaggerEnabled       bool   `mapstructure:"SERVER_APP_SWAGGER_ENABLED"`
	SwaggerUser          string `mapstructure:"SERVER_APP_SWAGGER_USER"`
	SwaggerPass          string `mapstructure:"SERVER_APP_SWAGGER_PASS"`
	// Observability configuration
	OtelEnabled     bool   `mapstructure:"SERVER_APP_OTEL_ENABLED"`
	OtelServiceName string `mapstructure:"SERVER_APP_OTEL_SERVICE_NAME"`
	JaegerEndpoint  string `mapstructure:"SERVER_APP_JAEGER_ENDPOINT"`
	// Optional batching configuration (leave empty for defaults)
	OtelBatchTimeout       int `mapstructure:"SERVER_APP_OTEL_BATCH_TIMEOUT"`         // Default: 5 seconds
	OtelMaxExportBatchSize int `mapstructure:"SERVER_APP_OTEL_MAX_EXPORT_BATCH_SIZE"` // Default: 512
	OtelMaxQueueSize       int `mapstructure:"SERVER_APP_OTEL_MAX_QUEUE_SIZE"`        // Default: 2048
	OtelExportTimeout      int `mapstructure:"SERVER_APP_OTEL_EXPORT_TIMEOUT"`        // Default: 30 seconds
}

func LoadConfig(path string) (*Conf, error) {
	// Carrega o .env se existir (ignora erro se não existir)
	err := godotenv.Load(path + "/.env")
	if err != nil {
		println(path + "/.env file not found")
		println("No .env file found, using environment variables", err.Error())
	}

	cfg := &Conf{
		AppName:                getEnv("SERVER_APP_NAME", "go_app_base"),
		ImageName:              getEnv("SERVER_APP_IMAGE_NAME", ""),
		ImageVersion:           getEnv("SERVER_APP_IMAGE_VERSION", ""),
		Environment:            getEnv("SERVER_APP_ENVIRONMENT", "development"),
		WebServerPort:          getEnv("SERVER_APP_WEB_SERVER_PORT", "8080"),
		DBDriver:               getEnv("SERVER_APP_DB_DRIVER", "mysql"),
		DBHost:                 getEnv("SERVER_APP_DB_HOST", "localhost"),
		DBPort:                 getEnv("SERVER_APP_DB_PORT", "3316"),
		DBUser:                 getEnv("SERVER_APP_DB_USER", "root"),
		DBPassword:             getEnv("SERVER_APP_DB_PASSWORD", "root"),
		DBName:                 getEnv("SERVER_APP_DB_NAME", "go_app_base"),
		DBMaxOpenConnections:   getEnvAsInt("SERVER_APP_DB_MAX_OPEN_CONNECTIONS", 20),
		DBMaxIdleConnections:   getEnvAsInt("SERVER_APP_DB_MAX_IDLE_CONNECTIONS", 10),
		DBConnMaxLifetime:      getEnvAsInt("SERVER_APP_DB_CONN_MAX_LIFETIME", 1),
		DBConnMaxIdleTime:      getEnvAsInt("SERVER_APP_DB_CONN_MAX_IDLE_TIME", 10),
		DebugMode:              getEnvAsBool("SERVER_APP_DEBUG_MODE", false),
		SwaggerEnabled:         getEnvAsBool("SERVER_APP_SWAGGER_ENABLED", false),
		SwaggerUser:            getEnv("SERVER_APP_SWAGGER_USER", ""),
		SwaggerPass:            getEnv("SERVER_APP_SWAGGER_PASS", ""),
		OtelEnabled:            getEnvAsBool("SERVER_APP_OTEL_ENABLED", false),
		OtelServiceName:        getEnv("SERVER_APP_OTEL_SERVICE_NAME", "go_app_base"),
		JaegerEndpoint:         getEnv("SERVER_APP_JAEGER_ENDPOINT", "jaeger:4318"),
		OtelBatchTimeout:       getEnvAsInt("SERVER_APP_OTEL_BATCH_TIMEOUT", 5),
		OtelMaxExportBatchSize: getEnvAsInt("SERVER_APP_OTEL_MAX_EXPORT_BATCH_SIZE", 512),
		OtelMaxQueueSize:       getEnvAsInt("SERVER_APP_OTEL_MAX_QUEUE_SIZE", 2048),
		OtelExportTimeout:      getEnvAsInt("SERVER_APP_OTEL_EXPORT_TIMEOUT", 30),
	}

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

// Observability configuration getters (implements observability.ConfigProvider)
func (c *Conf) GetOtelEnabled() bool {
	return c.OtelEnabled
}

func (c *Conf) GetOtelServiceName() string {
	return c.OtelServiceName
}

func (c *Conf) GetJaegerEndpoint() string {
	return c.JaegerEndpoint
}

func (c *Conf) GetEnvironment() string {
	return c.Environment
}

func (c *Conf) GetOtelBatchTimeout() int {
	return c.OtelBatchTimeout
}

func (c *Conf) GetOtelMaxExportBatchSize() int {
	return c.OtelMaxExportBatchSize
}

func (c *Conf) GetOtelMaxQueueSize() int {
	return c.OtelMaxQueueSize
}

func (c *Conf) GetOtelExportTimeout() int {
	return c.OtelExportTimeout
}
