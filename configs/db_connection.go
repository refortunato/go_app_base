package configs

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/XSAM/otelsql"
	_ "github.com/go-sql-driver/mysql"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func NewMySQL(cfg *Conf) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=UTC", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Register instrumented driver if observability is enabled
	driverName := "mysql"
	if cfg.OtelEnabled {
		var err error
		driverName, err = otelsql.Register("mysql",
			otelsql.WithAttributes(
				semconv.DBSystemMySQL,
			),
			// Configure span options to avoid false errors
			otelsql.WithSpanOptions(otelsql.SpanOptions{
				DisableQuery:    false, // Keep query visible for debugging
				OmitRows:        true,  // Don't record row counts
				OmitConnPrepare: true,  // Skip prepare statement spans
				OmitConnQuery:   false, // Keep query spans
			}),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to register instrumented driver: %w", err)
		}
	}

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}

	// Configura o pool
	db.SetMaxOpenConns(cfg.DBMaxOpenConnections)                              // máximo de conexões abertas simultâneas
	db.SetMaxIdleConns(cfg.DBMaxIdleConnections)                              // conexões em idle (ociosas)
	db.SetConnMaxLifetime(time.Duration(cfg.DBConnMaxLifetime) * time.Hour)   // recicla conexões a cada X tempo
	db.SetConnMaxIdleTime(time.Duration(cfg.DBConnMaxIdleTime) * time.Minute) // idle máximo antes de destruir conexão

	// Testa conexão
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
