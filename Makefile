.PHONY: dev prod go-mod-tidy down clean help

# Versão do Go usada em todos os containers
GO_VERSION := 1.25.5

help: ## Mostra esta mensagem de ajuda
	@echo "Comandos disponíveis:"
	@echo "  make dev          - Inicia ambiente de desenvolvimento (mysql + app-dev)"
	@echo "  make prod         - Inicia ambiente de produção (mysql + app-prod)"
	@echo "  make go-mod-tidy  - Executa 'go mod tidy' usando container Docker"
	@echo "  make down         - Para todos os containers"
	@echo "  make clean        - Para containers e remove volumes"

dev: ## Inicia ambiente de desenvolvimento
	@echo "🚀 Iniciando ambiente de desenvolvimento..."
	docker-compose up --build mysql app-dev

prod: ## Inicia ambiente de produção
	@echo "🚀 Iniciando ambiente de produção..."
	docker-compose up --build mysql app-api

go-mod-tidy: ## Executa 'go mod tidy' usando container Docker
	@echo "📦 Executando go mod tidy..."
	docker run --rm \
		-v $(PWD):/app \
		-w /app \
		golang:$(GO_VERSION)-alpine \
		go mod tidy
	@echo "✅ Dependências atualizadas!"

down: ## Para todos os containers
	@echo "🛑 Parando containers..."
	docker-compose down

clean: ## Para containers e remove volumes
	@echo "🧹 Limpando containers e volumes..."
	docker-compose down -v
