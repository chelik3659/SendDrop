package storage

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileInfo struct {
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	ModTime   time.Time `json:"mod_time"`
	Extension string    `json:"extension"`
}

type Manager struct {
	BaseDir string
}

func NewManager(baseDir string) (*Manager, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, err
	}
	return &Manager{BaseDir: baseDir}, nil
}

func (m *Manager) ListFiles() ([]FileInfo, error) {
	entries, err := os.ReadDir(m.BaseDir)
	if err != nil {
		return nil, err
	}
	
	var files []FileInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		info, err := entry.Info()
		if err != nil {
			continue
		}
		
		// Декодируем имя файла (если оно было закодировано)
		name := entry.Name()
		if decoded, err := url.QueryUnescape(name); err == nil {
			name = decoded
		}
		
		files = append(files, FileInfo{
			Name:      name,
			Size:      info.Size(),
			ModTime:   info.ModTime(),
			Extension: strings.ToLower(filepath.Ext(name)),
		})
	}
	return files, nil
}

func (m *Manager) SaveFile(filename string, reader io.Reader) error {
	// Кодируем имя для безопасности
	encodedName := url.QueryEscape(filename)
	filePath := filepath.Join(m.BaseDir, encodedName)
	
	// Если файл существует, добавляем суффикс
	if _, err := os.Stat(filePath); err == nil {
		ext := filepath.Ext(filename)
		name := strings.TrimSuffix(filename, ext)
		timestamp := time.Now().Format("20060102_150405")
		filename = fmt.Sprintf("%s_%s%s", name, timestamp, ext)
		encodedName = url.QueryEscape(filename)
		filePath = filepath.Join(m.BaseDir, encodedName)
	}
	
	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()
	
	_, err = io.Copy(dst, reader)
	return err
}

func (m *Manager) GetFile(filename string) (string, error) {
	encodedName := url.QueryEscape(filename)
	filePath := filepath.Join(m.BaseDir, encodedName)
	
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found: %s", filename)
	}
	return filePath, nil
}

func (m *Manager) DeleteFile(filename string) error {
	encodedName := url.QueryEscape(filename)
	filePath := filepath.Join(m.BaseDir, encodedName)
	return os.Remove(filePath)
}

func (m *Manager) GetFileSize(filename string) (int64, error) {
	encodedName := url.QueryEscape(filename)
	filePath := filepath.Join(m.BaseDir, encodedName)
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}
