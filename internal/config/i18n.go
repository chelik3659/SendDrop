package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Locale struct {
	AppTitle          string `json:"app_title"`
	ServerOnline      string `json:"server_online"`
	ServerOffline     string `json:"server_offline"`
	StartServer       string `json:"start_server"`
	StopServer        string `json:"stop_server"`
	UploadFile        string `json:"upload_file"`
	CopyLink          string `json:"copy_link"`
	Download          string `json:"download"`
	Delete            string `json:"delete"`
	NoFiles           string `json:"no_files"`
	SupportAuthor     string `json:"support_author"`
	Settings          string `json:"settings"`
	Theme             string `json:"theme"`
	Language          string `json:"language"`
	Port              string `json:"port"`
	IPAddress         string `json:"ip_address"`
	FileUploaded      string `json:"file_uploaded"`
	FileDeleted       string `json:"file_deleted"`
	LinkCopied        string `json:"link_copied"`
	SelectFile        string `json:"select_file"`
	ConfirmDelete     string `json:"confirm_delete"`
	NewVersionFound   string `json:"new_version_found"`
	UpdateAvailable   string `json:"update_available"`
	CurrentVersion    string `json:"current_version"`
	LatestVersion     string `json:"latest_version"`
}

var (
	currentLocale *Locale
	translations  = make(map[string]*Locale)
)

func LoadLanguages() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	
	langDir := filepath.Join(filepath.Dir(exePath), "assets", "locales")
	
	// Загружаем английский
	enData, err := os.ReadFile(filepath.Join(langDir, "en.json"))
	if err != nil {
		return err
	}
	var en Locale
	if err := json.Unmarshal(enData, &en); err != nil {
		return err
	}
	translations["en"] = &en
	
	// Загружаем русский
	ruData, err := os.ReadFile(filepath.Join(langDir, "ru.json"))
	if err != nil {
		return err
	}
	var ru Locale
	if err := json.Unmarshal(ruData, &ru); err != nil {
		return err
	}
	translations["ru"] = &ru
	
	// Устанавливаем язык по умолчанию
	SetLanguage(AppConfig.Language)
	
	return nil
}

func SetLanguage(lang string) {
	if loc, ok := translations[lang]; ok {
		currentLocale = loc
	} else {
		currentLocale = translations["en"]
	}
	AppConfig.Language = lang
	Save()
}

func T() *Locale {
	if currentLocale == nil {
		// fallback
		currentLocale = translations["en"]
	}
	return currentLocale
}