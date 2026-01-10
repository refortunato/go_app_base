package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/refortunato/go_app_base/cmd/server/container"
	"github.com/refortunato/go_app_base/configs"
	"github.com/refortunato/go_app_base/internal/infra/web/routes"
	"github.com/refortunato/go_app_base/internal/shared/web/server"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := configs.NewMySQL(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Initialize dependency container
	c, err := container.New(db, cfg)
	if err != nil {
		panic(err)
	}

	// Determina qual serviço iniciar baseado nos argumentos
	mode := "api" // padrão
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	// Canal para capturar sinais de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Canal para erros de inicialização
	serverErr := make(chan error, 1)

	var srv server.Server

	switch mode {
	case "api":
		fmt.Println("Starting API server...")
		srv = server.NewGinServerWithRoutes(cfg.WebServerPort, routes.RegisterRoutes(c))

		// Inicia o servidor em uma goroutine
		go func() {
			if err := srv.Start(); err != nil {
				serverErr <- fmt.Errorf("API server error: %w", err)
			}
		}()

	case "rabbitmq":
		fmt.Println("Starting RabbitMQ consumer...")
		// TODO: Implementar consumidor RabbitMQ
		// server = rabbitmq.NewConsumer(cfg)
		// go func() {
		//     if err := server.Start(); err != nil {
		//         serverErr <- fmt.Errorf("RabbitMQ consumer error: %w", err)
		//     }
		// }()
		fmt.Println("RabbitMQ consumer not implemented yet")
		os.Exit(1)

	case "kafka":
		fmt.Println("Starting Kafka consumer...")
		// TODO: Implementar consumidor Kafka
		// server = kafka.NewConsumer(cfg)
		// go func() {
		//     if err := server.Start(); err != nil {
		//         serverErr <- fmt.Errorf("Kafka consumer error: %w", err)
		//     }
		// }()
		fmt.Println("Kafka consumer not implemented yet")
		os.Exit(1)

	case "grpc":
		fmt.Println("Starting gRPC server...")
		// TODO: Implementar servidor gRPC
		// server = grpc.NewServer(cfg)
		// go func() {
		//     if err := server.Start(); err != nil {
		//         serverErr <- fmt.Errorf("gRPC server error: %w", err)
		//     }
		// }()
		fmt.Println("gRPC server not implemented yet")
		os.Exit(1)

	default:
		fmt.Printf("Unknown mode: %s\n", mode)
		fmt.Println("Available modes: api (default), rabbitmq, kafka, grpc")
		os.Exit(1)
	}

	// Aguarda sinal de interrupção ou erro do servidor
	select {
	case err := <-serverErr:
		fmt.Printf("Server error: %v\n", err)
		os.Exit(1)
	case sig := <-quit:
		fmt.Printf("\nReceived signal: %v\n", sig)
		fmt.Println("Initiating graceful shutdown...")

		// Cria um contexto com timeout para o shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Executa o shutdown gracioso
		if srv != nil {
			if err := srv.Shutdown(ctx); err != nil {
				fmt.Printf("Error during shutdown: %v\n", err)
				os.Exit(1)
			}
		}

		// Fecha a conexão com o banco de dados
		if err := db.Close(); err != nil {
			fmt.Printf("Error closing database: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Server stopped gracefully")
	}
}
