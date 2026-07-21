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
	// Создаём папку для общих файлов
	exePath, _ := os.Executable()
	shareDir := filepath.Join(filepath.Dir(exePath), "sharing")
	os.MkdirAll(shareDir, 0755)

	window := gui.NewWindow()
	window.Run()
}
