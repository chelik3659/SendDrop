#!/bin/bash

echo "🔧 Обновление SendDrop до v0.1.0 (альфа-фикс)"

# 1. Обновляем зависимости
echo "📦 Обновление зависимостей..."
go get github.com/skip2/go-qrcode
go mod tidy

# 2. Пересобираем
echo "🔨 Пересборка..."
CGO_ENABLED=1 go build -ldflags="-s -w" -o bin/SendDrop main.go

# 3. Копируем assets
cp -r assets bin/ 2>/dev/null || true

# 4. Проверяем
if [ -f "bin/SendDrop" ]; then
    echo "✅ Сборка успешна!"
    echo "📁 Бинарник: $(pwd)/bin/SendDrop"
    echo ""
    echo "🚀 Запуск..."
    ./bin/SendDrop
else
    echo "❌ Ошибка сборки!"
    exit 1
fi
