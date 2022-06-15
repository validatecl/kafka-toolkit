mod:
	@echo "Vendoring..."
	@go mod vendor

test: 
	@echo "Ejecutando tests..."
	@go test -mod=vendor ./... -v

coverage: 
	@echo "Coverage..."
	@go test -mod=vendor ./... --coverprofile coverfile_out >> /dev/null
	@go tool cover -func coverfile_out

.PHONY: test coverage