package gui

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"SendDrop/internal/discovery"
	"SendDrop/internal/transfer"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type Window struct {
	app        fyne.App
	window     fyne.Window
	discovery  *discovery.Discovery
	transfer   *transfer.Transfer
	peerList   *widget.List
	peers      []discovery.Peer
	peerLabel  *widget.Label
	fileList   *widget.List
	files      []string
	shareDir   string
	selectedIP string
	statusLabel *widget.Label
}

func NewWindow(shareDir string) *Window {
	a := app.NewWithID("com.senddrop.desktop")
	w := a.NewWindow("SendDrop P2P")
	w.Resize(fyne.NewSize(750, 550))
	w.CenterOnScreen()

	os.MkdirAll(shareDir, 0755)

	win := &Window{
		app:      a,
		window:   w,
		shareDir: shareDir,
		peers:    []discovery.Peer{},
		files:    []string{},
	}
	win.transfer = transfer.NewTransfer()
	win.transfer.SetOnFileReceive(win.onFileReceived)
	go win.transfer.StartServer(shareDir)

	win.discovery = discovery.NewDiscovery("Desktop")
	win.discovery.SetCallbacks(win.onPeerAdd, win.onPeerRemove)
	go win.discovery.Start()

	win.buildUI()
	win.refreshFileList()
	return win
}

func (w *Window) buildUI() {
	title := widget.NewLabelWithStyle("📡 SendDrop P2P", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	w.statusLabel = widget.NewLabel("🔍 Поиск устройств...")

	w.peerLabel = widget.NewLabel("Устройства в сети:")
	w.peerList = widget.NewList(
		func() int { return len(w.peers) },
		func() fyne.CanvasObject {
			return widget.NewLabel("device")
		},
		func(id int, obj fyne.CanvasObject) {
			if id >= len(w.peers) {
				return
			}
			obj.(*widget.Label).SetText(w.peers[id].Name + " (" + w.peers[id].IP + ")")
		},
	)
	w.peerList.OnSelected = func(id int) {
		if id < len(w.peers) {
			w.selectedIP = w.peers[id].IP
			w.peerLabel.SetText("Выбрано: " + w.peers[id].Name + " (" + w.peers[id].IP + ")")
			w.statusLabel.SetText("✅ Выбрано устройство: " + w.peers[id].Name)
		}
	}

	w.fileList = widget.NewList(
		func() int { return len(w.files) },
		func() fyne.CanvasObject {
			return widget.NewLabel("file")
		},
		func(id int, obj fyne.CanvasObject) {
			if id >= len(w.files) {
				return
			}
			obj.(*widget.Label).SetText(w.files[id])
		},
	)

	sendBtn := widget.NewButton("📤 Отправить файл", func() {
		if w.selectedIP == "" {
			dialog.ShowInformation("Ошибка", "Сначала выберите устройство из списка", w.window)
			return
		}
		// Открываем диалог выбора файла
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w.window)
				return
			}
			if reader == nil {
				return
			}
			defer reader.Close()

			// Читаем файл полностью
			data, err := io.ReadAll(reader)
			if err != nil {
				dialog.ShowError(err, w.window)
				return
			}
			if len(data) == 0 {
				dialog.ShowInformation("Ошибка", "Файл пуст", w.window)
				return
			}

			filename := filepath.Base(reader.URI().Path())
			w.statusLabel.SetText("⏳ Отправка " + filename + " на " + w.selectedIP + "...")

			err = w.transfer.SendFile(w.selectedIP, filename, data)
			if err != nil {
				w.statusLabel.SetText("❌ Ошибка отправки")
				dialog.ShowError(err, w.window)
			} else {
				w.statusLabel.SetText("✅ Файл отправлен: " + filename)
				dialog.ShowInformation("Успех", "Файл отправлен на "+w.selectedIP, w.window)
			}
		}, w.window)
	})

	refreshBtn := widget.NewButton("🔄 Обновить список устройств", func() {
		w.statusLabel.SetText("🔄 Поиск устройств...")
		// Принудительно обновляем список
		w.peers = []discovery.Peer{}
		w.peerList.Refresh()
		// Отправляем широковещательный запрос через discovery (можно вызвать метод)
		go w.discovery.SendDiscoveryRequest()
		time.Sleep(2 * time.Second)
		w.statusLabel.SetText("🔍 Поиск завершён, найдено устройств: " + fmt.Sprint(len(w.peers)))
	})

	infoLabel := widget.NewLabel("📁 Папка с файлами: " + w.shareDir)
	infoLabel.Wrapping = fyne.TextWrapBreak

	content := container.NewBorder(
		container.NewVBox(
			title,
			widget.NewSeparator(),
			w.statusLabel,
			w.peerLabel,
			container.NewVScroll(w.peerList),
			widget.NewSeparator(),
			widget.NewLabel("📂 Ваши файлы:"),
			container.NewVScroll(w.fileList),
			container.NewHBox(sendBtn, refreshBtn, layout.NewSpacer()),
			infoLabel,
		),
		nil, nil, nil,
	)

	w.window.SetContent(content)
}

func (w *Window) onPeerAdd(peer discovery.Peer) {
	// Проверяем, нет ли уже такого IP
	for _, p := range w.peers {
		if p.IP == peer.IP {
			return
		}
	}
	w.peers = append(w.peers, peer)
	w.peerList.Refresh()
	w.statusLabel.SetText("🔍 Найдено устройств: " + fmt.Sprint(len(w.peers)))
}

func (w *Window) onPeerRemove(peer discovery.Peer) {
	for i, p := range w.peers {
		if p.IP == peer.IP {
			w.peers = append(w.peers[:i], w.peers[i+1:]...)
			break
		}
	}
	w.peerList.Refresh()
	w.statusLabel.SetText("🔍 Найдено устройств: " + fmt.Sprint(len(w.peers)))
}

func (w *Window) onFileReceived(filename string, data []byte) {
	w.refreshFileList()
	w.statusLabel.SetText("📥 Получен файл: " + filename)
}

func (w *Window) refreshFileList() {
	files, err := os.ReadDir(w.shareDir)
	if err != nil {
		return
	}
	w.files = []string{}
	for _, f := range files {
		if !f.IsDir() {
			info, _ := f.Info()
			size := ""
			if info != nil {
				size = fmt.Sprintf(" (%.1f KB)", float64(info.Size())/1024)
			}
			w.files = append(w.files, f.Name()+size)
		}
	}
	w.fileList.Refresh()
}

func (w *Window) Run() {
	w.window.ShowAndRun()
}
