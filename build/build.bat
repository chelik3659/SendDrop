@echo off
echo Building SendDrop for Windows...

REM Установка зависимостей
go mod download

REM Создание папки для бинарников
if not exist "bin" mkdir bin

REM Сборка для Windows (64-bit)
echo Building Windows 64-bit...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-H=windowsgui -s -w" -o bin/SendDrop.exe ../main.go

REM Сборка для Windows (32-bit)
echo Building Windows 32-bit...
set GOOS=windows
set GOARCH=386
go build -ldflags="-H=windowsgui -s -w" -o bin/SendDrop_386.exe ../main.go

echo Build complete! Check bin folder.
pause