package gui

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	
	"SendDrop/internal/config"
	"SendDrop/internal/qrcode"
	"SendDrop/internal/server"
	"SendDrop/internal/storage"
)

type Window struct {
	app       fyne.App
	window    fyne.Window
	server    *server.Server
	storage   *storage.Manager
	startBtn  *widget.Button
	status    *widget.Label
	qrImage   *canvas.Image
	fileList  *widget.List
	files     []storage.FileInfo
	uploadBtn *widget.Button
	copyBtn   *widget.Button
	qrData    []byte
}

func NewWindow() *Window {
	a := app.NewWithID("com.senddrop.app")
	w := a.NewWindow("SendDrop")
	w.Resize(fyne.NewSize(650, 550))
	w.CenterOnScreen()
	
	shareDir := filepath.Join(filepath.Dir(os.Args[0]), config.AppConfig.ShareDir)
	storageManager, _ := storage.NewManager(shareDir)
	srv := server.NewServer(config.AppConfig.Port, config.AppConfig.IP, shareDir)
	
	win := &Window{
		app:     a,
		window:  w,
		server:  srv,
		storage: storageManager,
		files:   []storage.FileInfo{},
	}
	
	win.buildUI()
	win.refreshFileList()
	
	return win
}

func (w *Window) buildUI() {
	// Заголовок
	title := widget.NewLabelWithStyle("📁 SendDrop", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	
	// Статус
	w.status = widget.NewLabel("🔴 Server: OFFLINE")
	w.status.Alignment = fyne.TextAlignCenter
	w.status.TextStyle = fyne.TextStyle{Bold: true}
	
	// QR код
	w.qrImage = canvas.NewImageFromResource(nil)
	w.qrImage.SetMinSize(fyne.NewSize(150, 150))
	w.qrImage.FillMode = canvas.ImageFillContain
	
	qrContainer := container.NewCenter(w.qrImage)
	
	// Список файлов
	w.fileList = widget.NewList(
		func() int { return len(w.files) },
		func() fyne.CanvasObject {
			return container.NewHBox(
				widget.NewIcon(theme.FileIcon()),
				widget.NewLabel("filename"),
				layout.NewSpacer(),
				widget.NewButtonWithIcon("", theme.DownloadIcon(), nil),
				widget.NewButtonWithIcon("", theme.DeleteIcon(), nil),
			)
		},
		func(id int, obj fyne.CanvasObject) {
			if id >= len(w.files) {
				return
			}
			file := w.files[id]
			hbox := obj.(*fyne.Container)
			
			label := hbox.Objects[1].(*widget.Label)
			label.SetText(fmt.Sprintf("%s (%s)", file.Name, formatFileSize(file.Size)))
			
			downloadBtn := hbox.Objects[3].(*widget.Button)
			downloadBtn.OnTapped = func() {
				w.downloadFile(file.Name)
			}
			
			deleteBtn := hbox.Objects[4].(*widget.Button)
			deleteBtn.OnTapped = func() {
				w.deleteFile(file.Name)
			}
		},
	)
	
	// Кнопки
	w.startBtn = widget.NewButtonWithIcon("▶ START SERVER", theme.MediaPlayIcon(), func() {
		w.toggleServer()
	})
	w.startBtn.Importance = widget.HighImportance
	
	w.uploadBtn = widget.NewButtonWithIcon("📤 UPLOAD FILE", theme.UploadIcon(), func() {
		w.uploadFile()
	})
	w.uploadBtn.Importance = widget.MediumImportance
	w.uploadBtn.Disable()
	
	w.copyBtn = widget.NewButtonWithIcon("📋 COPY LINK", theme.ContentCopyIcon(), func() {
		w.copyLink()
	})
	w.copyBtn.Importance = widget.LowImportance
	w.copyBtn.Disable()
	
	// Кнопка поддержки
	supportBtn := widget.NewHyperlink("❤️ Support Author", nil)
	supportBtn.OnTapped = func() {
		// Открыть ссылку
	}
	supportBtn.Alignment = fyne.TextAlignCenter
	
	// Версия
	versionLabel := widget.NewLabel("v0.1.0-alpha")
	versionLabel.Alignment = fyne.TextAlignTrailing
	versionLabel.TextStyle = fyne.TextStyle{Italic: true}
	
	// Компоновка
	content := container.NewBorder(
		container.NewVBox(
			title,
			widget.NewSeparator(),
			w.status,
			qrContainer,
			widget.NewSeparator(),
		),
		container.NewHBox(
			supportBtn,
			layout.NewSpacer(),
			versionLabel,
		),
		nil,
		nil,
		container.NewVBox(
			container.NewHBox(w.startBtn, w.uploadBtn, w.copyBtn),
			widget.NewSeparator(),
			widget.NewLabelWithStyle("📂 Files:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			widget.NewLabel("💡 Drag & Drop files here to upload"),
			container.NewVScroll(w.fileList),
		),
	)
	
	w.window.SetContent(content)
}

func (w *Window) toggleServer() {
	if w.server.IsOnline() {
		w.server.Stop()
		w.startBtn.SetText("▶ START SERVER")
		w.startBtn.SetIcon(theme.MediaPlayIcon())
		w.status.SetText("🔴 Server: OFFLINE")
		w.uploadBtn.Disable()
		w.copyBtn.Disable()
		w.qrImage.Resource = nil
		w.qrImage.Refresh()
	} else {
		w.server.Start()
		w.startBtn.SetText("⏹ STOP SERVER")
		w.startBtn.SetIcon(theme.MediaStopIcon())
		w.status.SetText(fmt.Sprintf("🟢 Server: ONLINE at %s", w.server.GetURL()))
		w.uploadBtn.Enable()
		w.copyBtn.Enable()
		
		// Генерируем QR
		qrData, err := qrcode.GenerateQR(w.server.GetURL(), 200)
		if err == nil && len(qrData) > 0 {
			w.qrData = qrData
			resource := fyne.NewStaticResource("qr.png", qrData)
			w.qrImage.Resource = resource
			w.qrImage.Refresh()
		}
	}
}

func (w *Window) uploadFile() {
	dialog.ShowFileOpen(func(r fyne.URIReadCloser, err error) {
		if err != nil || r == nil {
			return
		}
		defer r.Close()
		
		filename := filepath.Base(r.URI().Path())
		data, err := io.ReadAll(r)
		if err != nil {
			dialog.ShowError(err, w.window)
			return
		}
		
		if err := w.storage.SaveFile(filename, bytes.NewReader(data)); err != nil {
			dialog.ShowError(err, w.window)
			return
		}
		
		w.refreshFileList()
		dialog.ShowInformation("✅ Success", "File uploaded: "+filename, w.window)
	}, w.window)
}

func (w *Window) downloadFile(filename string) {
	if !w.server.IsOnline() {
		dialog.ShowInformation("❌ Error", "Server is offline!", w.window)
		return
	}
	
	filePath, err := w.storage.GetFile(filename)
	if err != nil {
		dialog.ShowError(err, w.window)
		return
	}
	
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil || writer == nil {
			return
		}
		defer writer.Close()
		
		data, err := os.ReadFile(filePath)
		if err != nil {
			dialog.ShowError(err, w.window)
			return
		}
		
		if _, err := writer.Write(data); err != nil {
			dialog.ShowError(err, w.window)
		}
	}, w.window)
}

func (w *Window) deleteFile(filename string) {
	dialog.ShowConfirm("🗑 Delete", "Delete file: "+filename+"?", func(confirmed bool) {
		if !confirmed {
			return
		}
		
		if err := w.storage.DeleteFile(filename); err != nil {
			dialog.ShowError(err, w.window)
			return
		}
		
		w.refreshFileList()
	}, w.window)
}

func (w *Window) copyLink() {
	if !w.server.IsOnline() {
		return
	}
	url := w.server.GetURL()
	w.window.Clipboard().SetContent(url)
	dialog.ShowInformation("📋 Copied", "Link copied to clipboard:\n"+url, w.window)
}

func (w *Window) refreshFileList() {
	files, err := w.storage.ListFiles()
	if err != nil {
		return
	}
	w.files = files
	w.fileList.Refresh()
}

func (w *Window) Run() {
	w.window.ShowAndRun()
}

// formatFileSize - форматирует размер в понятный вид
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
