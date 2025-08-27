FROM postgres:15-alpine

# Instalar extensiones adicionales útiles
RUN apk add --no-cache \
    postgresql-contrib \
    postgresql-dev \
    gcc \
    musl-dev \
    make

# Copiar scripts de inicialización
COPY ./init.sql /docker-entrypoint-initdb.d/
COPY ./seed.sql /docker-entrypoint-initdb.d/

# Configurar variables de entorno
ENV POSTGRES_DB=goasync
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=password

# Exponer puerto
EXPOSE 5432

# Comando por defecto
CMD ["postgres"]
