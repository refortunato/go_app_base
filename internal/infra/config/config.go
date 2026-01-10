package config

type Conf struct {
	AppName              string `mapstructure:"SERVER_APP_NAME"`
	ImageName            string `mapstructure:"SERVER_APP_IMAGE_NAME"`
	ImageVersion         string `mapstructure:"SERVER_APP_IMAGE_VERSION"`
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
}
