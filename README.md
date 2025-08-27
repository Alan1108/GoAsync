# GoAsync API

Una API REST moderna construida en Go con Gin framework, diseñada para ser escalable y fácil de mantener.

## 🚀 Características

- **Framework**: Gin para routing HTTP
- **Logging**: Logrus para logging estructurado
- **Configuración**: Variables de entorno con godotenv
- **Containerización**: Docker y Docker Compose
- **Hot Reload**: Air para desarrollo con recarga automática
- **Estructura**: Arquitectura limpia y organizada

## 📋 Prerrequisitos

- Go 1.21 o superior
- Docker y Docker Compose (opcional)
- Make (opcional, para usar comandos del Makefile)

## 🛠️ Instalación

### Opción 1: Instalación local

1. Clona el repositorio:

```bash
git clone <tu-repositorio>
cd goasync
```

2. Instala las dependencias:

```bash
make deps
# o manualmente:
go mod tidy
go mod download
```

3. Copia el archivo de variables de entorno:

```bash
cp env.example .env
```

4. Ejecuta la aplicación:

```bash
make run
# o manualmente:
go run main.go
```

### Opción 2: Con Docker

1. Construye y ejecuta con Docker Compose:

```bash
docker-compose up --build
```

2. O solo la aplicación:

```bash
make docker-build
make docker-run
```

## 🚀 Desarrollo

### Comandos útiles

```bash
make help          # Muestra todos los comandos disponibles
make dev           # Ejecuta con hot reload (Air)
make test          # Ejecuta los tests
make lint          # Ejecuta el linter
make format        # Formatea el código
make clean         # Limpia archivos generados
```

### Hot Reload

Para desarrollo con recarga automática:

```bash
make dev
```

Esto instalará Air automáticamente y ejecutará la aplicación con hot reload.

## 📁 Estructura del proyecto

```
goasync/
├── main.go              # Punto de entrada de la aplicación
├── go.mod               # Dependencias de Go
├── Makefile             # Comandos útiles
├── Dockerfile           # Configuración de Docker
├── docker-compose.yml   # Orquestación de servicios
├── .air.toml           # Configuración de Air (hot reload)
├── env.example         # Variables de entorno de ejemplo
└── README.md           # Este archivo
```

## 🌐 Endpoints disponibles

- `GET /health` - Verificación de salud de la API
- `GET /api/v1/` - Información de la API v1

## 🔧 Configuración

### Variables de entorno

Crea un archivo `.env` basado en `env.example`:

```bash
# Configuración del servidor
PORT=8080
GIN_MODE=debug

# Configuración de logs
LOG_LEVEL=info
```

### Puertos

- **API**: 8080 (configurable via PORT)
- **PostgreSQL**: 5432
- **Redis**: 6379

## 🧪 Testing

```bash
make test              # Tests básicos
make test-coverage     # Tests con cobertura
```

## 📦 Build

```bash
make build             # Compila la aplicación
```

El binario se generará en `bin/goasync`.

## 🐳 Docker

### Construir imagen

```bash
make docker-build
```

### Ejecutar contenedor

```bash
make docker-run
```

### Con Docker Compose

```bash
docker-compose up --build
```

## 📝 Logs

La aplicación usa Logrus para logging estructurado. Los logs se muestran en consola en modo desarrollo.

## 🤝 Contribuir

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.

## 🆘 Soporte

Si tienes alguna pregunta o problema, por favor abre un issue en el repositorio.
