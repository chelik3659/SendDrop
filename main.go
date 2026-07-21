package main

import (
	"log"
	"os"
	"path/filepath"
	
	"SendDrop/internal/config"
	"SendDrop/internal/gui"
)

const (
	Version = "0.1.0-alpha"
	BuildDate = "2026-07-21"
)

func main() {
	log.Printf("🚀 SendDrop %s (Build: %s)", Version, BuildDate)
	
	// Инициализация конфига
	if err := config.Init(); err != nil {
		log.Printf("⚠️ Config init: %v", err)
	}
	
	// Создаём папку sharing
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	shareDir := filepath.Join(filepath.Dir(exePath), config.AppConfig.ShareDir)
	if err := os.MkdirAll(shareDir, 0755); err != nil {
		log.Fatal(err)
	}
	
	// Загружаем языки
	if err := config.LoadLanguages(); err != nil {
		log.Printf("⚠️ Error loading languages: %v", err)
	}
	
	// Запускаем GUI
	window := gui.NewWindow()
	window.Run()
}
