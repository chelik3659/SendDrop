package gui

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (w *Window) showSettingsDialog() {
	portEntry := widget.NewEntry()
	portEntry.SetText("8081")
	
	content := container.NewVBox(
		widget.NewLabel("Settings"),
		widget.NewLabel("Port:"),
		portEntry,
		widget.NewButton("Close", func() {}),
	)
	
	dialog.ShowCustom("Settings", "Close", content, w.window)
}
