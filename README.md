# MercadoLibre Coupon Challenge

API implementada en Go para el challenge de MercadoLibre que optimiza el uso de cupones de descuento.

## Estructura del Proyecto

```
meli-coupon/
├── cmd/api/            # Punto de entrada de la aplicación
├── internal/           # Código privado de la aplicación
│   ├── domain/         # Modelos y errores
│   ├── handlers/       # Manejadores HTTP
│   ├── services/       # Lógica de negocio
│   └── repository/     # Acceso a datos externos
├── pkg/                # Código público reutilizable
│   └── calculator/     # Algoritmo de combinaciones
└── fly.toml            # Configuración de Fly.io

```

## Características

- Endpoint POST `/coupon/`: Calcula la mejor combinación de productos para un cupón
- Endpoint GET `/coupon/stats`: Muestra los 5 productos más favoritados
- Manejo concurrente de estadísticas
- Tests unitarios
- Documentación completa
- Manejo de errores personalizado

## Requisitos

- Go 1.20 o superior
- Conexión a Internet (para acceder a la API de MercadoLibre)

## Instalación

```bash
# Clonar el repositorio
git clone https://github.com/fabriconiglio/meli-coupon-challenge.git

# Entrar al directorio
cd meli-coupon

# Instalar dependencias
go mod tidy

# Ejecutar tests
go test ./...

# Ejecutar la aplicación
go run cmd/api/main.go
```

## Uso Local

### Calcular mejor combinación de productos:
```bash
curl -X POST http://localhost:8080/coupon/ \
-H "Content-Type: application/json" \
-d '{
    "item_ids": ["MLA1", "MLA2", "MLA3"],
    "amount": 500
}'
```

### Obtener estadísticas:
```bash
curl http://localhost:8080/coupon/stats
```

## Decisiones de Diseño

1. **Clean Architecture**: Separación clara de responsabilidades en capas
2. **Dependency Injection**: Facilita testing y mantenimiento
3. **Concurrencia**: Uso de mutex para estadísticas thread-safe
4. **Error Handling**: Errores personalizados y descriptivos
5. **Testing**: Cobertura completa con mocks

# Ejecuta la aplicación como un contenedor de Docker
```bash
docker build -t meli-coupon .
docker run -p 8080:8080 meli-coupon
```

## Implementacion de cache
## Primera petición (sin caché)
```bash
curl -X POST http://localhost:8080/coupon/ \
-H "Content-Type: application/json" \
-d '{
    "item_ids": ["MLA1114676780", "MLA1114683943", "MLA1114656461"],
    "amount": 500000
}'
```

## Segunda petición (debería usar caché)
```bash
curl -X POST http://localhost:8080/coupon/ \
-H "Content-Type: application/json" \
-d '{
    "item_ids": ["MLA1114676780", "MLA1114683943", "MLA1114656461"],
    "amount": 500000
}'
```

## Agreguemos un endpoint de healthcheck y métricas:

## Healthcheck
```bash
curl http://localhost:8080/health
```

## Métricas
```bash
curl http://localhost:8080/metrics
```


## Deployment en Fly.io
## Requisitos

-Cuenta en Fly.io
-Fly CLI instalado

## Pasos para deployar

1) Instalar Fly CLI:

```bash
curl -L https://fly.io/install.sh | sh
```

2) Agregar Fly CLI al PATH:

```bash
export FLYCTL_INSTALL="/home/$USER/.fly"
export PATH="$FLYCTL_INSTALL/bin:$PATH"
```

3) Login en Fly:

```bash
flyctl auth login
```

4) Deployar la aplicación:
```bash
flyctl launch
```

## URL de la aplicación:
```bash
https://meli-coupon-autumn-leaf-7889.fly.dev
```

# Healthcheck y Métricas
## Healthcheck
```bash
https://meli-coupon-autumn-leaf-7889.fly.dev/health
```

## Métricas
```bash
https://meli-coupon-autumn-leaf-7889.fly.dev/metrics
```

## Autor

Fabrizzio Coniglio

## Licencia

MIT