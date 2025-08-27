.PHONY: help build run test clean deps lint db-up db-down db-reset seed seed-docker seed-massive seed-small

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

# Comandos de base de datos
db-up: ## Levanta los servicios de base de datos con Docker Compose
	docker-compose up -d postgres

db-down: ## Detiene los servicios de base de datos
	docker-compose down

db-reset: ## Reinicia la base de datos (elimina volúmenes)
	docker-compose down -v
	docker-compose up -d postgres

db-logs: ## Muestra logs de la base de datos
	docker-compose logs -f postgres

db-connect: ## Conecta a la base de datos PostgreSQL
	docker-compose exec postgres psql -U postgres -d goasync

# Comandos de seeder
seed: ## Ejecuta el seeder localmente
	./scripts/seed-db.sh

seed-docker: ## Ejecuta el seeder en Docker
	./scripts/seed-db.sh --docker

seed-clean: ## Ejecuta el seeder limpiando la base de datos primero
	./scripts/seed-db.sh --clean

seed-massive: ## Ejecuta el seeder masivo (miles de registros)
	./scripts/seed-db.sh --massive

seed-massive-docker: ## Ejecuta el seeder masivo en Docker
	./scripts/seed-db.sh --massive --docker

seed-small: ## Ejecuta el seeder pequeño (solo datos básicos)
	./scripts/seed-db.sh --small

seed-small-docker: ## Ejecuta el seeder pequeño en Docker
	./scripts/seed-db.sh --small --docker

# Comandos de Docker
docker-build: ## Construye la imagen Docker
	docker build -t $(BINARY_NAME) .

docker-run: ## Ejecuta el contenedor Docker
	docker run -p 8080:8080 $(BINARY_NAME)

docker-compose-up: ## Levanta todos los servicios con Docker Compose
	docker-compose up --build

docker-compose-down: ## Detiene todos los servicios
	docker-compose down

# Comandos de desarrollo completo
dev-full: ## Levanta todos los servicios y ejecuta la aplicación
	@echo "🚀 Iniciando entorno de desarrollo completo..."
	@make db-up
	@echo "⏳ Esperando que la base de datos esté lista..."
	@sleep 10
	@make seed
	@echo "🎯 Ejecutando aplicación..."
	@make run

dev-massive: ## Levanta BD, ejecuta seeder masivo y ejecuta la aplicación
	@echo "🚀 Iniciando entorno de desarrollo con datos masivos..."
	@make db-up
	@echo "⏳ Esperando que la base de datos esté lista..."
	@sleep 10
	@make seed-massive
	@echo "🎯 Ejecutando aplicación..."
	@make run

dev-docker: ## Ejecuta todo en Docker
	@echo "🐳 Iniciando entorno Docker completo..."
	@make docker-compose-up

# Comandos de utilidad
db-stats: ## Muestra estadísticas de la base de datos
	@echo "📊 Estadísticas de la base de datos:"
	@docker-compose exec postgres psql -U postgres -d goasync -c "SELECT * FROM get_database_stats();"

db-backup: ## Crea un backup de la base de datos
	@echo "💾 Creando backup de la base de datos..."
	@docker-compose exec postgres pg_dump -U postgres goasync > backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "✅ Backup creado exitosamente"

db-restore: ## Restaura la base de datos desde un backup
	@echo "🔄 Restaurando base de datos desde backup..."
	@if [ -z "$(BACKUP_FILE)" ]; then \
		echo "❌ Especifica el archivo de backup: make db-restore BACKUP_FILE=backup.sql"; \
		exit 1; \
	fi
	@docker-compose exec -T postgres psql -U postgres goasync < $(BACKUP_FILE)
	@echo "✅ Base de datos restaurada exitosamente"
