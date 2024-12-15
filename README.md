# MercadoLibre Coupon Challenge

API implementada en Go para el challenge de MercadoLibre que optimiza el uso de cupones de descuento.

## Estructura del Proyecto

```
meli-coupon/
├── cmd/api/           # Punto de entrada de la aplicación
├── internal/          # Código privado de la aplicación
│   ├── domain/       # Modelos y errores
│   ├── handlers/     # Manejadores HTTP
│   ├── services/     # Lógica de negocio
│   └── repository/   # Acceso a datos externos
└── pkg/              # Código público reutilizable
    └── calculator/   # Algoritmo de combinaciones
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
git clone https://github.com/tu-usuario/meli-coupon.git

# Entrar al directorio
cd meli-coupon

# Instalar dependencias
go mod tidy

# Ejecutar tests
go test ./...

# Ejecutar la aplicación
go run cmd/api/main.go
```

## Uso

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

## Mejoras Futuras

1. Implementar caché para precios de productos
2. Agregar métricas y monitoring
3. Documentación con Swagger
4. Containerización con Docker
5. CI/CD pipeline

## Autor

[Tu Nombre]

## Licencia

MIT