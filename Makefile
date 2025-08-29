.PHONY: help build run-api run-web run-cli clean dev-setup install-cli

help:
	@echo "Available commands:"
	@echo "  build       - Build all components"
	@echo "  run-api     - Run the API server"
	@echo "  run-web     - Run the web development server"
	@echo "  run-cli     - Build and show CLI help"
	@echo "  install-cli - Build and install dbx CLI to ~/bin"
	@echo "  dev-setup   - Start development environment (postgres)"
	@echo "  clean       - Clean build artifacts"

build: build-api build-cli

build-api:
	cd api && go build -o ../bin/server cmd/server/main.go

build-cli:
	cd cli && go build -o ../bin/dbx cmd/dbx/main.go

run-api:
	cd api && go run cmd/server/main.go

run-web:
	cd web && npm run dev

run-cli: build-cli
	./bin/dbx --help

dev-setup:
	docker-compose up -d postgres

dev-down:
	docker-compose down

# Install CLI to user's local bin directory
install-cli: build-cli
	@echo "📦 Installing dbx CLI..."
	@mkdir -p ~/bin
	@cp bin/dbx ~/bin/dbx
	@chmod +x ~/bin/dbx
	@echo "✅ dbx installed to ~/bin/dbx"
	@echo ""
	@echo "To use dbx from anywhere, ensure ~/bin is in your PATH:"
	@echo "  echo 'export PATH=\"\$$HOME/bin:\$$PATH\"' >> ~/.bashrc"
	@echo "  echo 'export PATH=\"\$$HOME/bin:\$$PATH\"' >> ~/.zshrc"
	@echo ""
	@echo "Or reload your shell and try: dbx --help"

clean:
	rm -rf bin/
	rm -f ~/bin/dbx
	cd api && go clean
	cd cli && go clean
	cd web && rm -rf dist/ node_modules/.cache/

install-deps:
	cd api && go mod tidy
	cd cli && go mod tidy
	cd web && npm install

.PHONY: build-api build-cli