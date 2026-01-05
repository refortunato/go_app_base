package configs

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	internalConf "github.com/refortunato/go_app_base/internal/infra/config"
)

func NewMySQL(cfg *internalConf.Conf) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=UTC", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open("mysql", dsn)
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
