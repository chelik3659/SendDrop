#!/bin/bash

echo "🚀 Starting SendDrop..."

# Проверяем, собран ли бинарник
if [ ! -f "../bin/SendDrop" ]; then
    echo "❌ Binary not found! Building first..."
    cd ../build
    ./build.sh
    cd ..
fi

# Запускаем приложение
cd ..
./bin/SendDrop