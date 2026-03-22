# transaction_routine

Ensure your Docker is running.

## API Documentation (Swagger/OpenAPI)

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **OpenAPI YAML**: http://localhost:8080/openapi.yaml

## Tests

```bash
# Unit tests only (no database required)
go test ./... -short

# Integration tests (requires PostgreSQL - run docker-compose up first)
go test ./... -tags=integration

# All tests
go test ./...
```

Or using Make (from `scripts/` directory):
- `make test` - unit tests
- `make test-integration` - integration tests
- `make test-all` - all tests

// testes
// camada de auditoria e logs
// encrypt document_number
// ver como ficaria camada de cache
