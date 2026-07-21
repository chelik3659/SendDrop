# 📡 SendDrop P2P – Desktop

[![Version](https://img.shields.io/badge/version-0.2.0--alpha-blue)](https://github.com/uniko404/SendDrop/releases)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.21+-blue)](https://golang.org/)

**SendDrop P2P for Desktop** — децентрализованный обмен файлами по локальной сети без сервера. Устройства видят друг друга автоматически, файлы передаются напрямую. Никаких серверов, никаких настроек.

---

## ✨ Features

- 🚀 **Peer‑to‑peer** – каждый участник равноправен, нет единой точки отказа
- 📡 **Автоматическое обнаружение** – устройства видят друг друга в локальной сети (UDP broadcast)
- 📁 **Отправка файлов** – выберите устройство и отправьте файл
- 📂 **Приём файлов** – автоматическое сохранение в `sharing/`
- 🔒 **Без интернета** – всё работает по локальной сети
- 🌐 **Кроссплатформенность** – совместим с Android-версией SendDrop P2P

---

## 📥 Download

[![Download](https://img.shields.io/badge/download-latest-green)](https://github.com/uniko404/SendDrop/releases)

Скачайте бинарник для вашей ОС из раздела [Releases](https://github.com/uniko404/SendDrop/releases).

---

## 🚀 Quick Start

1. Скачайте бинарник для вашей ОС
2. Запустите: `./SendDrop` (Linux/macOS) или `SendDrop.exe` (Windows)
3. Приложение автоматически найдёт другие устройства в сети
4. Выберите устройство из списка
5. Нажмите **"📤 Отправить файл"** и выберите файл
6. Файл будет передан напрямую на выбранное устройство

---

## 🛠️ Build from Source

```bash
git clone https://github.com/uniko404/SendDrop.git
cd SendDrop
go mod download
CGO_ENABLED=1 go build -ldflags="-s -w" -o SendDrop main.go
