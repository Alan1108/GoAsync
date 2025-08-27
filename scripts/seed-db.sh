#!/bin/bash

# Script para ejecutar el seeder de la base de datos GoAsync
# Uso: ./scripts/seed-db.sh [opciones]

set -e

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Función para mostrar mensajes
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Función para mostrar ayuda
show_help() {
    echo "Uso: $0 [opciones]"
    echo ""
    echo "Opciones:"
    echo "  -h, --help          Mostrar esta ayuda"
    echo "  -e, --env-file      Archivo de variables de entorno (default: .env)"
    echo "  -c, --clean         Limpiar base de datos antes de insertar datos"
    echo "  -v, --verbose       Modo verbose"
    echo "  -d, --docker        Ejecutar seeder en contenedor Docker"
    echo "  -m, --massive       Ejecutar seeder masivo (miles de registros)"
    echo "  -s, --small         Ejecutar seeder pequeño (solo datos básicos)"
    echo ""
    echo "Ejemplos:"
    echo "  $0                    # Ejecutar seeder con configuración por defecto"
    echo "  $0 --clean           # Limpiar DB y ejecutar seeder"
    echo "  $0 --docker          # Ejecutar en contenedor Docker"
    echo "  $0 --massive         # Ejecutar seeder masivo"
    echo "  $0 --massive --docker # Seeder masivo en Docker"
    echo "  $0 -e .env.local     # Usar archivo .env.local"
}

# Variables por defecto
ENV_FILE=".env"
CLEAN_DB=false
VERBOSE=false
USE_DOCKER=false
SEEDER_MODE="default" # default, massive, small

# Parsear argumentos
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -e|--env-file)
            ENV_FILE="$2"
            shift 2
            ;;
        -c|--clean)
            CLEAN_DB=true
            shift
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -d|--docker)
            USE_DOCKER=true
            shift
            ;;
        -m|--massive)
            SEEDER_MODE="massive"
            shift
            ;;
        -s|--small)
            SEEDER_MODE="small"
            shift
            ;;
        *)
            log_error "Opción desconocida: $1"
            show_help
            exit 1
            ;;
    esac
done

# Verificar si estamos en el directorio correcto
if [[ ! -f "go.mod" ]]; then
    log_error "Este script debe ejecutarse desde el directorio raíz del proyecto"
    exit 1
fi

# Verificar si existe el archivo de variables de entorno
if [[ ! -f "$ENV_FILE" ]]; then
    log_warning "Archivo $ENV_FILE no encontrado, usando variables del sistema"
fi

# Función para verificar dependencias
check_dependencies() {
    log_info "Verificando dependencias..."
    
    if ! command -v go &> /dev/null; then
        log_error "Go no está instalado o no está en el PATH"
        exit 1
    fi
    
    if ! command -v docker &> /dev/null && [[ "$USE_DOCKER" == true ]]; then
        log_error "Docker no está instalado o no está en el PATH"
        exit 1
    fi
    
    log_success "Dependencias verificadas"
}

# Función para instalar dependencias de Go
install_go_deps() {
    log_info "Instalando dependencias de Go..."
    
    if [[ "$VERBOSE" == true ]]; then
        go mod tidy -v
        go mod download
    else
        go mod tidy
        go mod download
    fi
    
    log_success "Dependencias de Go instaladas"
}

# Función para ejecutar seeder localmente
run_seeder_local() {
    log_info "Ejecutando seeder localmente (modo: $SEEDER_MODE)..."
    
    # Construir argumentos del seeder
    seeder_args=""
    case $SEEDER_MODE in
        "massive")
            seeder_args="--massive"
            ;;
        "small")
            seeder_args="--small"
            ;;
        *)
            seeder_args=""
            ;;
    esac
    
    if [[ "$VERBOSE" == true ]]; then
        go run cmd/seeder/main.go $seeder_args
    else
        go run cmd/seeder/main.go $seeder_args 2>&1 | grep -E "(INFO|SUCCESS|ERROR|WARNING|🎉|✅|🧹|📂|👥|🏷️|📝|💬|🔗|📊)"
    fi
    
    if [ $? -eq 0 ]; then
        log_success "Seeder ejecutado localmente"
    else
        log_error "Error ejecutando seeder localmente"
        return 1
    fi
}

# Función para ejecutar seeder en Docker
run_seeder_docker() {
    log_info "Ejecutando seeder en Docker (modo: $SEEDER_MODE)..."
    
    # Verificar si el contenedor de la base de datos está ejecutándose
    if ! docker ps | grep -q "goasync-postgres"; then
        log_error "El contenedor de PostgreSQL no está ejecutándose"
        log_info "Ejecuta 'make db-up' primero"
        exit 1
    fi
    
    # Ejecutar seeder en contenedor
    docker run --rm \
        --network goasync_goasync-network \
        --env-file "$ENV_FILE" \
        -e DB_HOST=postgres \
        -e DB_PORT=5432 \
        -e DB_NAME=goasync \
        -e DB_USER=postgres \
        -e DB_PASSWORD=password \
        -e DB_SSLMODE=disable \
        goasync-app \
        go run cmd/seeder/main.go
    
    log_success "Seeder ejecutado en Docker"
}

# Función para verificar conexión a la base de datos
check_db_connection() {
    log_info "Verificando conexión a la base de datos..."
    
    # Obtener variables de entorno
    source "$ENV_FILE" 2>/dev/null || true
    
    DB_HOST="${DB_HOST:-localhost}"
    DB_PORT="${DB_PORT:-5432}"
    DB_NAME="${DB_NAME:-goasync}"
    DB_USER="${DB_USER:-postgres}"
    DB_PASSWORD="${DB_PASSWORD:-password}"
    
    # Verificar si PostgreSQL está ejecutándose
    if command -v pg_isready &> /dev/null; then
        if pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" &> /dev/null; then
            log_success "Conexión a PostgreSQL verificada"
            return 0
        fi
    fi
    
    # Verificar si Docker está ejecutándose
    if docker ps | grep -q "goasync-postgres"; then
        log_success "Contenedor PostgreSQL ejecutándose en Docker"
        return 0
    fi
    
    log_warning "No se pudo verificar la conexión a PostgreSQL"
    log_info "Asegúrate de que PostgreSQL esté ejecutándose o usa --docker"
    return 1
}

# Función para mostrar información del modo de seeder
show_seeder_info() {
    case $SEEDER_MODE in
        "massive")
            log_info "🚀 Modo MASIVO seleccionado"
            echo "   - 1000 usuarios"
            echo "   - 15 categorías"
            echo "   - 100+ tags"
            echo "   - 5000 posts"
            echo "   - 15000 comentarios"
            echo "   - 25000 relaciones post-tag"
            echo "   - 10000 logs de actividad"
            echo "   - Total: ~60,000+ registros"
            ;;
        "small")
            log_info "📝 Modo PEQUEÑO seleccionado"
            echo "   - 5 usuarios básicos"
            echo "   - 8 categorías"
            echo "   - 10 tags"
            echo "   - 4 posts de ejemplo"
            echo "   - Comentarios básicos"
            echo "   - Total: ~30 registros"
            ;;
        *)
            log_info "⚡ Modo por defecto seleccionado"
            echo "   - Datos básicos para desarrollo"
            ;;
    esac
}

# Función principal
main() {
    log_info "🚀 Iniciando seeder de base de datos GoAsync"
    
    # Mostrar información del modo seleccionado
    show_seeder_info
    
    # Verificar dependencias
    check_dependencies
    
    # Instalar dependencias de Go
    install_go_deps
    
    # Verificar conexión a la base de datos
    if ! check_db_connection; then
        if [[ "$USE_DOCKER" == false ]]; then
            log_warning "¿Quieres ejecutar el seeder en Docker? Usa --docker"
        fi
    fi
    
    # Ejecutar seeder
    if [[ "$USE_DOCKER" == true ]]; then
        run_seeder_docker
    else
        run_seeder_local
    fi
    
    log_success "🎉 Seeder completado exitosamente!"
    
    # Mostrar estadísticas finales
    if [[ "$SEEDER_MODE" == "massive" ]]; then
        log_info "📊 Base de datos poblada con miles de registros"
        log_info "💡 Usa 'make db-connect' para explorar los datos"
        log_info "🌐 Accede a pgAdmin en http://localhost:5050"
    fi
}

# Ejecutar función principal
main "$@"
