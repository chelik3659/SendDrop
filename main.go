package main

import (
	"log"
	"os"
	"path/filepath"

	"SendDrop/internal/gui"
)

const version = "0.2.0-alpha"

func main() {
	log.Printf("🚀 SendDrop P2P %s", version)

	// Используем домашнюю папку пользователя
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Cannot get home directory:", err)
	}

	shareDir := filepath.Join(home, "SendDrop")
	if err := os.MkdirAll(shareDir, 0755); err != nil {
		log.Fatal("Cannot create share directory:", err)
	}

	log.Printf("📁 Файлы сохраняются в: %s", shareDir)

	window := gui.NewWindow(shareDir)
	window.Run()
}
