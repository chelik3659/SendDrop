.PHONY: build clean run install-deps

BINARY_NAME=SendDrop
BIN_DIR=../bin

# Установка зависимостей
install-deps:
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy

# Сборка для текущей платформы
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	mkdir -p $(BIN_DIR)
	go build -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY_NAME) ../main.go
	chmod +x $(BIN_DIR)/$(BINARY_NAME)
	@echo "✅ Build complete: $(BIN_DIR)/$(BINARY_NAME)"

# Кросс-компиляция для всех платформ
build-all:
	@echo "🔨 Building for all platforms..."
	mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY_NAME)_linux_amd64 ../main.go
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY_NAME)_linux_arm64 ../main.go
	GOOS=linux GOARCH=arm go build -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY_NAME)_linux_arm ../main.go
	GOOS=linux GOARCH=386 go build -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY_NAME)_linux_386 ../main.go
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY_NAME)_windows_amd64.exe ../main.go
	GOOS=windows GOARCH=386 go build -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY_NAME)_windows_386.exe ../main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY_NAME)_darwin_amd64 ../main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY_NAME)_darwin_arm64 ../main.go
	chmod +x $(BIN_DIR)/$(BINARY_NAME)*
	@echo "✅ Build complete! Check $(BIN_DIR)/"

# Запуск приложения
run:
	@echo "🚀 Running $(BINARY_NAME)..."
	cd .. && ./bin/$(BINARY_NAME)

# Очистка
clean:
	@echo "🧹 Cleaning..."
	rm -rf $(BIN_DIR)
	rm -f ../sharing/*

# Помощь
help:
	@echo "Available commands:"
	@echo "  make install-deps  - Install Go dependencies"
	@echo "  make build         - Build for current platform"
	@echo "  make build-all     - Build for all platforms"
	@echo "  make run          - Build and run"
	@echo "  make clean        - Clean build artifacts"