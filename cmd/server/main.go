package main

import (
	"fmt"
	"os"

	"github.com/refortunato/go_app_base/configs"
	"github.com/refortunato/go_app_base/internal/infra/dependencies"
	"github.com/refortunato/go_app_base/internal/infra/web/webserver"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	println("Config loaded:", cfg.DBHost, cfg.DBName, cfg.DBPort)

	db, err := configs.NewMySQL(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := dependencies.InitDependencies(db, cfg); err != nil {
		panic(err)
	}

	// Determina qual serviço iniciar baseado nos argumentos
	mode := "api" // padrão
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	switch mode {
	case "api":
		fmt.Println("Starting API server...")
		webserver.Start(cfg.WebServerPort)

	case "rabbitmq":
		fmt.Println("Starting RabbitMQ consumer...")
		// TODO: Implementar consumidor RabbitMQ

	case "kafka":
		fmt.Println("Starting Kafka consumer...")
		// TODO: Implementar consumidor Kafka

	case "grpc":
		fmt.Println("Starting gRPC server...")
		// TODO: Implementar servidor gRPC

	default:
		fmt.Printf("Unknown mode: %s\n", mode)
		fmt.Println("Available modes: api (default), rabbitmq, kafka, grpc")
		os.Exit(1)
	}
}
