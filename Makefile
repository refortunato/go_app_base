.PHONY: dev prod go-mod-tidy go-get swagger jaeger-ui down clean help

# Versão do Go usada em todos os containers
GO_VERSION := 1.25.5

help: ## Mostra esta mensagem de ajuda
	@echo "Comandos disponíveis:"
	@echo "  make dev          - Inicia ambiente de desenvolvimento (mysql + app-dev + jaeger)"
	@echo "  make prod         - Inicia ambiente de produção (mysql + app-prod + jaeger)"
	@echo "  make go-get       - Instala dependências Go via Docker (use: make go-get DEPS='pkg1 pkg2')"
	@echo "  make go-mod-tidy  - Executa 'go mod tidy' usando container Docker"
	@echo "  make swagger      - Gera documentação Swagger"
	@echo "  make jaeger-ui    - Abre Jaeger UI no navegador"
	@echo "  make down         - Para todos os containers"
	@echo "  make clean        - Para containers e remove volumes"

dev: ## Inicia ambiente de desenvolvimento
	@echo "🚀 Iniciando ambiente de desenvolvimento..."
	docker-compose up --build mysql jaeger app-dev

prod: ## Inicia ambiente de produção
	@echo "🚀 Iniciando ambiente de produção..."
	docker-compose up --build mysql jaeger app-api

go-get: ## Instala dependências Go (uso: make go-get DEPS='package1 package2')
	@if [ -z "$(DEPS)" ]; then \
		echo "❌ Erro: Especifique as dependências com DEPS='pkg1 pkg2'"; \
		echo "Exemplo: make go-get DEPS='go.opentelemetry.io/otel@v1.24.0'"; \
		exit 1; \
	fi
	@echo "📦 Instalando dependências Go via Docker..."
	docker run --rm \
		-v $(PWD):/app \
		-w /app \
		golang:$(GO_VERSION)-alpine \
		go get $(DEPS)
	@echo "✅ Dependências instaladas! Execute 'make go-mod-tidy' para limpar."

go-mod-tidy: ## Executa 'go mod tidy' usando container Docker
	@echo "📦 Executando go mod tidy..."
	docker run --rm \
		-v $(PWD):/app \
		-w /app \
		golang:$(GO_VERSION)-alpine \
		go mod tidy
	@echo "✅ Dependências atualizadas!"

swagger: ## Gera documentação Swagger
	@echo "📝 Gerando documentação Swagger..."
	@docker run --rm \
		-v $(PWD):/app \
		-w /app \
		golang:$(GO_VERSION)-alpine \
		sh -c "apk add --no-cache git && go install github.com/swaggo/swag/cmd/swag@latest && /go/bin/swag init -g cmd/server/main.go -o docs"
	@echo "✅ Swagger gerado! Acesse: http://localhost:8080/swagger/index.html"

jaeger-ui: ## Abre Jaeger UI no navegador
	@echo "🔍 Abrindo Jaeger UI..."
	@open http://localhost:16686 || xdg-open http://localhost:16686 || echo "Abra manualmente: http://localhost:16686"

down: ## Para todos os containers
	@echo "🛑 Parando containers..."
	docker-compose down

clean: ## Para containers e remove volumes
	@echo "🧹 Limpando containers e volumes..."
	docker-compose down -v
