#!/bin/bash

echo "🔨 Building SendDrop for Linux (Arch Linux)..."

# Устанавливаем переменные окружения для pkg-config
export PKG_CONFIG_PATH=/usr/lib/pkgconfig:$PKG_CONFIG_PATH

# Проверяем наличие Go
if ! command -v go &> /dev/null; then
    echo "❌ Go not found! Installing..."
    sudo pacman -S go
fi

# Устанавливаем зависимости
echo "📦 Downloading dependencies..."
go mod download
go mod tidy

# Создаём папку для бинарников
mkdir -p ../bin

# Сборка для Linux (текущая архитектура)
echo "🏗️ Building for Linux (native)..."
CGO_ENABLED=1 go build -ldflags="-s -w" -o ../bin/SendDrop ../main.go

# Проверяем успешность сборки
if [ -f "../bin/SendDrop" ]; then
    echo "✅ Build successful!"
    echo "📁 Binary location: ../bin/SendDrop"
    chmod +x ../bin/SendDrop
    ls -la ../bin/
else
    echo "❌ Build failed! Check errors above."
    exit 1
fi