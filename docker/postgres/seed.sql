-- Script de seeder para poblar la base de datos GoAsync con datos de ejemplo
-- Este script se ejecuta después de init.sql

-- Insertar categorías de ejemplo
INSERT INTO categories (name, description, slug) VALUES
('Tecnología', 'Artículos sobre tecnología, programación y desarrollo', 'tecnologia'),
('Ciencia', 'Artículos sobre ciencia, investigación y descubrimientos', 'ciencia'),
('Salud', 'Artículos sobre salud, bienestar y medicina', 'salud'),
('Educación', 'Artículos sobre educación, aprendizaje y desarrollo personal', 'educacion'),
('Entretenimiento', 'Artículos sobre entretenimiento, cultura y ocio', 'entretenimiento'),
('Deportes', 'Artículos sobre deportes, fitness y actividades físicas', 'deportes'),
('Negocios', 'Artículos sobre negocios, emprendimiento y economía', 'negocios'),
('Viajes', 'Artículos sobre viajes, turismo y aventuras', 'viajes');

-- Insertar tags de ejemplo
INSERT INTO tags (name, slug, description) VALUES
('Go', 'go', 'Lenguaje de programación Go'),
('API', 'api', 'Interfaces de programación de aplicaciones'),
('Docker', 'docker', 'Plataforma de contenedores'),
('PostgreSQL', 'postgresql', 'Base de datos relacional'),
('Web Development', 'web-development', 'Desarrollo web'),
('Microservicios', 'microservicios', 'Arquitectura de microservicios'),
('Cloud Computing', 'cloud-computing', 'Computación en la nube'),
('DevOps', 'devops', 'Prácticas de desarrollo y operaciones'),
('Machine Learning', 'machine-learning', 'Aprendizaje automático'),
('Data Science', 'data-science', 'Ciencia de datos');

-- Insertar usuarios de ejemplo (password: 'password123')
INSERT INTO users (username, email, password_hash, first_name, last_name) VALUES
('admin', 'admin@goasync.com', crypt('password123', gen_salt('bf')), 'Admin', 'User'),
('johndoe', 'john.doe@example.com', crypt('password123', gen_salt('bf')), 'John', 'Doe'),
('janesmith', 'jane.smith@example.com', crypt('password123', gen_salt('bf')), 'Jane', 'Smith'),
('bobwilson', 'bob.wilson@example.com', crypt('password123', gen_salt('bf')), 'Bob', 'Wilson'),
('alicebrown', 'alice.brown@example.com', crypt('password123', gen_salt('bf')), 'Alice', 'Brown');

-- Insertar perfiles de usuario
INSERT INTO user_profiles (user_id, bio, avatar_url, phone, address) VALUES
((SELECT id FROM users WHERE username = 'admin'), 'Administrador del sistema', 'https://example.com/avatars/admin.jpg', '+1234567890', '123 Admin St, City, Country'),
((SELECT id FROM users WHERE username = 'johndoe'), 'Desarrollador Go apasionado por la tecnología', 'https://example.com/avatars/john.jpg', '+1234567891', '456 Developer Ave, Tech City, Country'),
((SELECT id FROM users WHERE username = 'janesmith'), 'Arquitecta de software especializada en microservicios', 'https://example.com/avatars/jane.jpg', '+1234567892', '789 Architect Blvd, Design City, Country'),
((SELECT id FROM users WHERE username = 'bobwilson'), 'DevOps engineer con experiencia en cloud computing', 'https://example.com/avatars/bob.jpg', '+1234567893', '321 DevOps Rd, Cloud City, Country'),
((SELECT id FROM users WHERE username = 'alicebrown'), 'Data scientist y experta en machine learning', 'https://example.com/avatars/alice.jpg', '+1234567894', '654 Data St, Science City, Country');

-- Insertar posts de ejemplo
INSERT INTO posts (title, slug, content, excerpt, author_id, category_id, status, published_at) VALUES
(
    'Introducción a Go: El lenguaje del futuro',
    'introduccion-a-go-el-lenguaje-del-futuro',
    'Go es un lenguaje de programación desarrollado por Google que combina la simplicidad de Python con el rendimiento de C++. En este artículo exploraremos sus características principales, ventajas y casos de uso.

## ¿Por qué Go?

Go fue diseñado para resolver problemas específicos que los desarrolladores enfrentan en el desarrollo de software moderno:

- **Simplicidad**: Sintaxis clara y fácil de aprender
- **Rendimiento**: Compilación rápida y ejecución eficiente
- **Concurrencia**: Goroutines y channels para programación concurrente
- **Garbage Collection**: Gestión automática de memoria
- **Compilación estática**: Binarios únicos sin dependencias externas

## Características principales

### 1. Goroutines
Las goroutines son funciones que se ejecutan de forma concurrente. Son ligeras y eficientes:

```go
func main() {
    go func() {
        fmt.Println("Ejecutando en background")
    }()
    
    fmt.Println("Función principal")
    time.Sleep(time.Second)
}
```

### 2. Channels
Los channels permiten la comunicación entre goroutines:

```go
func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        fmt.Printf("worker %d processing job %d\n", id, j)
        time.Sleep(time.Second)
        results <- j * 2
    }
}
```

### 3. Interfaces
Go usa interfaces implícitas, lo que hace el código más flexible:

```go
type Animal interface {
    Speak() string
}

type Dog struct{}
type Cat struct{}

func (d Dog) Speak() string { return "Woof!" }
func (c Cat) Speak() string { return "Meow!" }
```

## Casos de uso

Go es ideal para:
- **Microservicios**: APIs y servicios web
- **Herramientas CLI**: Aplicaciones de línea de comandos
- **Sistemas distribuidos**: Aplicaciones que requieren alta concurrencia
- **DevOps**: Herramientas de automatización y CI/CD

## Conclusión

Go es un lenguaje moderno que resuelve muchos de los problemas del desarrollo de software actual. Su simplicidad, rendimiento y soporte nativo para concurrencia lo hacen ideal para aplicaciones modernas y escalables.',
    'Go es un lenguaje de programación moderno que combina simplicidad y rendimiento. Descubre por qué es ideal para el desarrollo de software actual.',
    (SELECT id FROM users WHERE username = 'johndoe'),
    (SELECT id FROM categories WHERE slug = 'tecnologia'),
    'published',
    CURRENT_TIMESTAMP - INTERVAL '5 days'
),
(
    'Construyendo APIs RESTful con Go y Gin',
    'construyendo-apis-restful-con-go-y-gin',
    'En este artículo aprenderemos a construir APIs RESTful robustas y escalables usando Go y el framework Gin.

## ¿Qué es Gin?

Gin es un framework web HTTP para Go que ofrece un rendimiento excelente y una API intuitiva. Es conocido por su velocidad y facilidad de uso.

## Configuración inicial

Primero, inicializamos nuestro módulo Go:

```bash
go mod init myapi
go get github.com/gin-gonic/gin
```

## Estructura básica de una API

```go
package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func main() {
    r := gin.Default()
    
    // Rutas
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "pong",
        })
    })
    
    r.Run(":8080")
}
```

## Middleware personalizado

Gin permite crear middleware personalizado fácilmente:

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Token requerido"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

## Validación de datos

Usando la librería validator:

```go
type User struct {
    Username string `json:"username" binding:"required,min=3,max=20"`
    Email    string `json:"email" binding:"required,email"`
    Age      int    `json:"age" binding:"gte=0,lte=130"`
}

func CreateUser(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    // Procesar usuario...
}
```

## Manejo de errores

```go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        
        if len(c.Errors) > 0 {
            c.JSON(http.StatusInternalServerError, gin.H{
                "errors": c.Errors.String(),
            })
        }
    }
}
```

## Testing

```go
func TestPingEndpoint(t *testing.T) {
    gin.SetMode(gin.TestMode)
    r := gin.New()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "pong"})
    })
    
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/ping", nil)
    r.ServeHTTP(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
}
```

## Conclusión

Go y Gin proporcionan una base sólida para construir APIs RESTful modernas y escalables. La combinación de rendimiento, simplicidad y herramientas disponibles hace de Go una excelente elección para el desarrollo de APIs.',
    'Aprende a construir APIs RESTful robustas usando Go y Gin. Descubre las mejores prácticas y patrones para crear servicios web escalables.',
    (SELECT id FROM users WHERE username = 'janesmith'),
    (SELECT id FROM categories WHERE slug = 'tecnologia'),
    'published',
    CURRENT_TIMESTAMP - INTERVAL '3 days'
),
(
    'Docker para desarrolladores: Una guía práctica',
    'docker-para-desarrolladores-una-guia-practica',
    'Docker ha revolucionado la forma en que desarrollamos, desplegamos y ejecutamos aplicaciones. En esta guía práctica aprenderemos los conceptos fundamentales y cómo aplicarlos en el desarrollo diario.

## ¿Qué es Docker?

Docker es una plataforma que permite empaquetar aplicaciones en contenedores ligeros y portables. Cada contenedor incluye todo lo necesario para ejecutar la aplicación: código, runtime, herramientas del sistema, bibliotecas y configuraciones.

## Conceptos básicos

### Imagen
Una imagen es un template de solo lectura que contiene el código de la aplicación, runtime, bibliotecas y configuraciones.

### Contenedor
Un contenedor es una instancia ejecutable de una imagen. Puedes crear múltiples contenedores a partir de la misma imagen.

### Dockerfile
Un Dockerfile es un script que contiene instrucciones para construir una imagen Docker.

## Creando tu primer Dockerfile

```dockerfile
# Usar imagen base de Node.js
FROM node:18-alpine

# Establecer directorio de trabajo
WORKDIR /app

# Copiar archivos de dependencias
COPY package*.json ./

# Instalar dependencias
RUN npm ci --only=production

# Copiar código fuente
COPY . .

# Exponer puerto
EXPOSE 3000

# Comando por defecto
CMD ["npm", "start"]
```

## Comandos básicos

```bash
# Construir imagen
docker build -t myapp .

# Ejecutar contenedor
docker run -p 3000:3000 myapp

# Listar contenedores
docker ps

# Detener contenedor
docker stop <container_id>

# Eliminar contenedor
docker rm <container_id>
```

## Docker Compose

Para aplicaciones más complejas, Docker Compose permite definir y ejecutar aplicaciones multi-contenedor:

```yaml
version: "3.8"
services:
  app:
    build: .
    ports:
      - "3000:3000"
    environment:
      - NODE_ENV=production
    depends_on:
      - db
  
  db:
    image: postgres:15
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

## Mejores prácticas

1. **Usar imágenes base oficiales**: Siempre usa imágenes oficiales y específicas
2. **Multi-stage builds**: Reduce el tamaño final de la imagen
3. **No ejecutar como root**: Crea usuarios no privilegiados
4. **Optimizar capas**: Agrupa comandos RUN para reducir capas
5. **Usar .dockerignore**: Excluye archivos innecesarios

## Debugging

```bash
# Entrar al contenedor
docker exec -it <container_id> /bin/bash

# Ver logs
docker logs <container_id>

# Inspeccionar contenedor
docker inspect <container_id>
```

## Conclusión

Docker simplifica el desarrollo y despliegue de aplicaciones, proporcionando un entorno consistente y reproducible. Su adopción en el desarrollo moderno es esencial para cualquier desarrollador.',
    'Descubre cómo Docker puede mejorar tu flujo de trabajo de desarrollo. Aprende conceptos fundamentales y mejores prácticas para usar contenedores.',
    (SELECT id FROM users WHERE username = 'bobwilson'),
    (SELECT id FROM categories WHERE slug = 'tecnologia'),
    'published',
    CURRENT_TIMESTAMP - INTERVAL '1 day'
),
(
    'PostgreSQL vs MySQL: ¿Cuál elegir para tu proyecto?',
    'postgresql-vs-mysql-cual-elegir-para-tu-proyecto',
    'La elección de la base de datos es una decisión crítica en cualquier proyecto. En este artículo compararemos PostgreSQL y MySQL para ayudarte a tomar la decisión correcta.

## PostgreSQL: El elefante robusto

PostgreSQL es una base de datos objeto-relacional de código abierto conocida por su robustez, extensibilidad y cumplimiento de estándares SQL.

### Ventajas de PostgreSQL

- **Cumplimiento de estándares**: Soporte completo para SQL estándar
- **Tipos de datos avanzados**: Arrays, JSON, UUID, geometrías
- **Extensibilidad**: Permite crear tipos de datos personalizados
- **Integridad referencial**: Soporte robusto para constraints
- **Concurrencia**: MVCC (Multi-Version Concurrency Control)
- **Extensiones**: Sistema rico de extensiones

### Casos de uso ideales

- Aplicaciones empresariales complejas
- Sistemas GIS y geoespaciales
- Análisis de datos y data warehousing
- Aplicaciones que requieren integridad de datos

## MySQL: El más popular

MySQL es la base de datos de código abierto más popular del mundo, conocida por su simplicidad y rendimiento.

### Ventajas de MySQL

- **Simplicidad**: Fácil de configurar y usar
- **Rendimiento**: Optimizado para aplicaciones web
- **Comunidad**: Gran comunidad y documentación
- **Herramientas**: Muchas herramientas de administración disponibles
- **Replicación**: Replicación maestro-esclavo robusta

### Casos de uso ideales

- Aplicaciones web simples a medianas
- Sitios web y blogs
- Aplicaciones que requieren alta velocidad de lectura
- Proyectos con recursos limitados

## Comparación técnica

| Característica | PostgreSQL | MySQL |
|----------------|------------|-------|
| Tipos de datos | Muy ricos | Básicos |
| Transacciones | ACID completo | ACID (InnoDB) |
| Concurrencia | MVCC | Row-level locking |
| Extensibilidad | Alta | Limitada |
| Rendimiento | Bueno | Excelente (lecturas) |
| Escalabilidad | Horizontal y vertical | Principalmente vertical |

## Ejemplos de código

### PostgreSQL

```sql
-- Crear tipo personalizado
CREATE TYPE user_status AS ENUM ('active', 'inactive', 'suspended');

-- Usar JSON
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    data JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Consulta con JSON
SELECT data->>'name' as name, data->>'email' as email
FROM users
WHERE data->>'status' = 'active';
```

### MySQL

```sql
-- Crear tabla
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    status ENUM('active', 'inactive', 'suspended') DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Índices
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_email ON users(email);
```

## Factores a considerar

### 1. Complejidad del proyecto
- **Proyectos simples**: MySQL puede ser suficiente
- **Proyectos complejos**: PostgreSQL ofrece más flexibilidad

### 2. Equipo de desarrollo
- **Desarrolladores junior**: MySQL es más fácil de aprender
- **Desarrolladores senior**: PostgreSQL aprovecha mejor las características avanzadas

### 3. Requisitos de rendimiento
- **Alta velocidad de lectura**: MySQL
- **Operaciones complejas**: PostgreSQL

### 4. Escalabilidad
- **Crecimiento vertical**: Ambos son buenos
- **Crecimiento horizontal**: PostgreSQL tiene ventajas

## Migración entre bases de datos

Migrar entre PostgreSQL y MySQL puede ser complejo debido a las diferencias en tipos de datos y sintaxis. Herramientas como:

- **pgloader**: Para migrar de MySQL a PostgreSQL
- **AWS DMS**: Para migraciones en la nube
- **Scripts personalizados**: Para casos específicos

## Conclusión

La elección entre PostgreSQL y MySQL depende de:

- **PostgreSQL**: Para proyectos que requieren robustez, extensibilidad y cumplimiento de estándares
- **MySQL**: Para proyectos que priorizan simplicidad, rendimiento y facilidad de uso

Ambas son excelentes opciones, pero cada una brilla en diferentes escenarios. Evalúa cuidadosamente los requisitos de tu proyecto antes de tomar la decisión.',
    'Compara PostgreSQL y MySQL para elegir la base de datos correcta para tu proyecto. Analiza ventajas, casos de uso y factores técnicos.',
    (SELECT id FROM users WHERE username = 'alicebrown'),
    (SELECT id FROM categories WHERE slug = 'tecnologia'),
    'published',
    CURRENT_TIMESTAMP
);

-- Asociar tags con posts
INSERT INTO post_tags (post_id, tag_id) VALUES
((SELECT id FROM posts WHERE slug = 'introduccion-a-go-el-lenguaje-del-futuro'), (SELECT id FROM tags WHERE slug = 'go')),
((SELECT id FROM posts WHERE slug = 'introduccion-a-go-el-lenguaje-del-futuro'), (SELECT id FROM tags WHERE slug = 'web-development')),
((SELECT id FROM posts WHERE slug = 'construyendo-apis-restful-con-go-y-gin'), (SELECT id FROM tags WHERE slug = 'go')),
((SELECT id FROM posts WHERE slug = 'construyendo-apis-restful-con-go-y-gin'), (SELECT id FROM tags WHERE slug = 'api')),
((SELECT id FROM posts WHERE slug = 'construyendo-apis-restful-con-go-y-gin'), (SELECT id FROM tags WHERE slug = 'web-development')),
((SELECT id FROM posts WHERE slug = 'docker-para-desarrolladores-una-guia-practica'), (SELECT id FROM tags WHERE slug = 'docker')),
((SELECT id FROM posts WHERE slug = 'docker-para-desarrolladores-una-guia-practica'), (SELECT id FROM tags WHERE slug = 'devops')),
((SELECT id FROM posts WHERE slug = 'postgresql-vs-mysql-cual-elegir-para-tu-proyecto'), (SELECT id FROM tags WHERE slug = 'postgresql')),
((SELECT id FROM posts WHERE slug = 'postgresql-vs-mysql-cual-elegir-para-tu-proyecto'), (SELECT id FROM tags WHERE slug = 'data-science'));

-- Insertar comentarios de ejemplo
INSERT INTO comments (post_id, author_id, content, is_approved) VALUES
((SELECT id FROM posts WHERE slug = 'introduccion-a-go-el-lenguaje-del-futuro'), (SELECT id FROM users WHERE username = 'janesmith'), 'Excelente artículo! Go realmente es un lenguaje increíble para microservicios.', true),
((SELECT id FROM posts WHERE slug = 'introduccion-a-go-el-lenguaje-del-futuro'), (SELECT id FROM users WHERE username = 'bobwilson'), 'Muy útil para principiantes. ¿Podrías hacer un artículo sobre testing en Go?', true),
((SELECT id FROM posts WHERE slug = 'construyendo-apis-restful-con-go-y-gin'), (SELECT id FROM users WHERE username = 'johndoe'), 'Gin es mi framework favorito para Go. Muy bien explicado.', true),
((SELECT id FROM posts WHERE slug = 'docker-para-desarrolladores-una-guia-practica'), (SELECT id FROM users WHERE username = 'alicebrown'), 'Docker ha cambiado completamente mi flujo de trabajo. Gran guía!', true);

-- Insertar logs de actividad de ejemplo
INSERT INTO activity_logs (user_id, action, resource_type, resource_id, details, ip_address) VALUES
((SELECT id FROM users WHERE username = 'admin'), 'user_login', 'user', (SELECT id FROM users WHERE username = 'admin'), '{"ip": "192.168.1.100", "user_agent": "Mozilla/5.0..."}', '192.168.1.100'),
((SELECT id FROM users WHERE username = 'johndoe'), 'post_created', 'post', (SELECT id FROM posts WHERE slug = 'introduccion-a-go-el-lenguaje-del-futuro'), '{"title": "Introducción a Go: El lenguaje del futuro"}', '192.168.1.101'),
((SELECT id FROM users WHERE username = 'janesmith'), 'comment_added', 'comment', (SELECT id FROM comments WHERE content LIKE '%Excelente artículo%'), '{"post_id": "' || (SELECT id FROM posts WHERE slug = 'introduccion-a-go-el-lenguaje-del-futuro') || '"}', '192.168.1.102');

-- Crear vistas útiles
CREATE OR REPLACE VIEW active_posts AS
SELECT 
    p.id,
    p.title,
    p.slug,
    p.excerpt,
    p.status,
    p.published_at,
    u.username as author_username,
    c.name as category_name,
    COUNT(co.id) as comment_count
FROM posts p
LEFT JOIN users u ON p.author_id = u.id
LEFT JOIN categories c ON p.category_id = c.id
LEFT JOIN comments co ON p.id = co.post_id AND co.is_approved = true
WHERE p.status = 'published'
GROUP BY p.id, u.username, c.name
ORDER BY p.published_at DESC;

CREATE OR REPLACE VIEW user_stats AS
SELECT 
    u.id,
    u.username,
    u.email,
    COUNT(p.id) as posts_count,
    COUNT(co.id) as comments_count,
    u.created_at as joined_at
FROM users u
LEFT JOIN posts p ON u.id = p.author_id
LEFT JOIN comments co ON u.id = co.author_id
GROUP BY u.id, u.username, u.email, u.created_at;

-- Crear índices adicionales para mejorar consultas
CREATE INDEX IF NOT EXISTS idx_posts_created_at ON posts(created_at);
CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments(created_at);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

-- Crear función para estadísticas de la base de datos
CREATE OR REPLACE FUNCTION get_database_stats()
RETURNS TABLE (
    total_users BIGINT,
    total_posts BIGINT,
    total_comments BIGINT,
    total_categories BIGINT,
    total_tags BIGINT
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        (SELECT COUNT(*) FROM users)::BIGINT,
        (SELECT COUNT(*) FROM posts)::BIGINT,
        (SELECT COUNT(*) FROM comments)::BIGINT,
        (SELECT COUNT(*) FROM categories)::BIGINT,
        (SELECT COUNT(*) FROM tags)::BIGINT;
END;
$$ LANGUAGE plpgsql;
