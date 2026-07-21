package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ProgressDialog struct {
	window   fyne.Window
	progress *widget.ProgressBar
	label    *widget.Label
}

func NewProgressDialog(parent fyne.Window, title string) *ProgressDialog {
	progress := widget.NewProgressBar()
	label := widget.NewLabel("Preparing...")
	label.Alignment = fyne.TextAlignCenter
	
	content := container.NewVBox(
		label,
		progress,
	)
	
	dialog := widget.NewModalPopUp(content, parent.Canvas())
	dialog.Resize(fyne.NewSize(300, 100))
	
	return &ProgressDialog{
		window:   parent,
		progress: progress,
		label:    label,
	}
}

func (p *ProgressDialog) Show() {
	// Показываем как модальное окно
	p.progress.SetValue(0)
	p.label.SetText("Starting...")
}

func (p *ProgressDialog) SetProgress(value float64) {
	p.progress.SetValue(value)
}

func (p *ProgressDialog) SetLabel(text string) {
	p.label.SetText(text)
}

func (p *ProgressDialog) Close() {
	// Закрываем
}
