# Configuration

This directory contains the configuration loading and database connection setup for the application.

## Files

- `config.go` - Loads environment variables with `SERVER_APP_` prefix
- `db_connection.go` - Sets up MySQL database connection with connection pooling

## Environment Variables

All environment variables should be prefixed with `SERVER_APP_`:

- `SERVER_APP_NAME` - Application name (default: "go_app_base")
- `SERVER_APP_WEB_SERVER_PORT` - Web server port (default: "8080")
- `SERVER_APP_DB_DRIVER` - Database driver (default: "mysql")
- `SERVER_APP_DB_HOST` - Database host (default: "localhost")
- `SERVER_APP_DB_PORT` - Database port (default: "3316")
- `SERVER_APP_DB_USER` - Database user (default: "root")
- `SERVER_APP_DB_PASSWORD` - Database password (default: "root")
- `SERVER_APP_DB_NAME` - Database name (default: "go_app_base")
- `SERVER_APP_DB_MAX_OPEN_CONNECTIONS` - Max open connections (default: 20)
- `SERVER_APP_DB_MAX_IDLE_CONNECTIONS` - Max idle connections (default: 10)
- `SERVER_APP_DB_CONN_MAX_LIFETIME` - Connection max lifetime in hours (default: 1)
- `SERVER_APP_DB_CONN_MAX_IDLE_TIME` - Connection max idle time in minutes (default: 10)
- `SERVER_APP_DEBUG_MODE` - Debug mode (default: false)
