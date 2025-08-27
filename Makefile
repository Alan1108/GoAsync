.PHONY: help build run test clean deps lint

# Variables
BINARY_NAME=goasync
MAIN_FILE=main.go

# Comando por defecto
help: ## Muestra esta ayuda
	@echo "Comandos disponibles:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

deps: ## Instala las dependencias
	go mod tidy
	go mod download

build: ## Compila la aplicación
	go build -o bin/$(BINARY_NAME) $(MAIN_FILE)

run: ## Ejecuta la aplicación
	go run $(MAIN_FILE)

dev: ## Ejecuta en modo desarrollo con hot reload
	@echo "Instalando air para hot reload..."
	go install github.com/cosmtrek/air@latest
	air

test: ## Ejecuta los tests
	go test -v ./...

test-coverage: ## Ejecuta los tests con cobertura
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean: ## Limpia los archivos generados
	rm -rf bin/
	rm -f coverage.out

lint: ## Ejecuta el linter
	@echo "Instalando golangci-lint..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run

format: ## Formatea el código
	go fmt ./...
	go vet ./...

docker-build: ## Construye la imagen Docker
	docker build -t $(BINARY_NAME) .

docker-run: ## Ejecuta el contenedor Docker
	docker run -p 8080:8080 $(BINARY_NAME)
