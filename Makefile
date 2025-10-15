.PHONY: help build run dev clean test docker-up docker-down migrate-up migrate-down logs

help:
	@echo "Comandos disponíveis:"
	@echo "  build        - Compila a aplicação"
	@echo "  run          - Executa a aplicação compilada"
	@echo "  dev          - Executa em modo desenvolvimento com hot reload"
	@echo "  clean        - Remove arquivos temporários"
	@echo "  test         - Executa os testes"
	@echo "  docker-up    - Inicia todos os serviços Docker"
	@echo "  docker-down  - Para todos os serviços Docker"
	@echo "  migrate-up   - Executa migrações do banco"
	@echo "  migrate-down - Reverte migrações do banco"
	@echo "  migrate-version - Verifica status das migrações"
	@echo "  logs         - Mostra logs da aplicação"

build:
	@echo "Compilando aplicação..."
	go build -o ./tmp/main ./cmd

run: build
	@echo "Executando aplicação..."
	./tmp/main

dev:
	@echo "Iniciando modo desenvolvimento com hot reload..."
	docker-compose up app

clean:
	@echo "Limpando arquivos temporários..."
	rm -rf ./tmp
	go clean

test:
	@echo "Executando testes..."
	go test ./...

docker-up:
	@echo "Iniciando serviços Docker..."
	docker-compose up -d

docker-down:
	@echo "Parando serviços Docker..."
	docker-compose down

migrate-up:
	@echo "Executando migrações..."
	MIGRATION_ACTION=up docker-compose up migrate

migrate-down:
	@echo "Revertendo migrações..."
	MIGRATION_ACTION=down docker-compose up migrate

migrate-version:
	@echo "Verificando status das migrações..."
	MIGRATION_ACTION=version docker-compose up migrate

logs:
	@echo "Mostrando logs da aplicação..."
	docker-compose logs -f app

setup:
	@echo "Configurando projeto..."
	@echo "1. Copiando .env.example para .env (se existir)"
	@if [ -f .env.example ]; then cp .env.example .env; echo "Arquivo .env criado"; else echo "Arquivo .env.example não encontrado"; fi
	@echo "2. Criando diretório tmp se não existir"
	@mkdir -p tmp
	@echo "3. Executando go mod tidy"
	go mod tidy
	@echo "Setup concluído!"

start: setup docker-up migrate-up
	@echo "Projeto iniciado com sucesso!"
	@echo "Acesse:"
	@echo "  - Aplicação: http://localhost:80"
	@echo "  - PgAdmin: http://localhost:8080"