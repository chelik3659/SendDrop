package gui

import (
	"fmt"
	"image/color"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type StatusBar struct {
	widget.BaseWidget
	statusLabel *widget.Label
	ipLabel     *widget.Label
	portLabel   *widget.Label
	statusIcon  *canvas.Circle
	container   *fyne.Container
}

func NewStatusBar() *StatusBar {
	s := &StatusBar{}
	
	s.statusIcon = canvas.NewCircle(color.NRGBA{R: 255, G: 0, B: 0, A: 255})
	s.statusIcon.Resize(fyne.NewSize(12, 12))
	
	s.statusLabel = widget.NewLabel("Server: OFFLINE")
	s.statusLabel.TextStyle = fyne.TextStyle{Bold: true}
	
	s.ipLabel = widget.NewLabel("IP: -")
	s.portLabel = widget.NewLabel("Port: -")
	
	statusContainer := container.NewHBox(s.statusIcon, s.statusLabel)
	infoContainer := container.NewHBox(
		widget.NewSeparator(),
		s.ipLabel,
		widget.NewSeparator(),
		s.portLabel,
	)
	
	s.container = container.NewBorder(
		nil, nil,
		statusContainer,
		infoContainer,
	)
	
	return s
}

func (s *StatusBar) SetOnline(online bool) {
	if online {
		s.statusIcon.FillColor = color.NRGBA{R: 0, G: 255, B: 0, A: 255}
		s.statusLabel.Text = "Server: ONLINE"
	} else {
		s.statusIcon.FillColor = color.NRGBA{R: 255, G: 0, B: 0, A: 255}
		s.statusLabel.Text = "Server: OFFLINE"
	}
	s.statusLabel.Refresh()
	s.statusIcon.Refresh()
}

func (s *StatusBar) SetIP(ip string) {
	s.ipLabel.Text = fmt.Sprintf("IP: %s", ip)
	s.ipLabel.Refresh()
}

func (s *StatusBar) SetPort(port int) {
	s.portLabel.Text = fmt.Sprintf("Port: %d", port)
	s.portLabel.Refresh()
}

func (s *StatusBar) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(s.container)
}

type QRWidget struct {
	widget.BaseWidget
	image    *canvas.Image
	container *fyne.Container
}

func NewQRWidget() *QRWidget {
	q := &QRWidget{
		image: canvas.NewImageFromResource(nil),
	}
	q.image.SetMinSize(fyne.NewSize(200, 200))
	q.image.FillMode = canvas.ImageFillContain
	
	placeholder := widget.NewLabel("Start server to generate QR code")
	placeholder.Alignment = fyne.TextAlignCenter
	
	q.container = container.NewCenter(
		container.NewVBox(
			q.image,
			placeholder,
		),
	)
	return q
}

func (q *QRWidget) SetQR(data []byte) {
	if len(data) == 0 {
		q.image.Resource = nil
	} else {
		resource := fyne.NewStaticResource("qr.png", data)
		q.image.Resource = resource
	}
	q.image.Refresh()
}

func (q *QRWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(q.container)
}

type CustomTheme struct {
	darkMode bool
}

func NewCustomTheme(dark bool) *CustomTheme {
	return &CustomTheme{darkMode: dark}
}

func (t *CustomTheme) Color(c fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	if t.darkMode {
		switch c {
		case fyne.ThemeColorNameBackground:
			return color.NRGBA{R: 30, G: 30, B: 30, A: 255}
		case fyne.ThemeColorNameForeground:
			return color.White
		default:
			return fyne.CurrentApp().Settings().Theme().Color(c, v)
		}
	}
	return fyne.CurrentApp().Settings().Theme().Color(c, v)
}

func (t *CustomTheme) Font(s fyne.TextStyle) fyne.Resource {
	return fyne.CurrentApp().Settings().Theme().Font(s)
}

func (t *CustomTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return fyne.CurrentApp().Settings().Theme().Icon(n)
}

func (t *CustomTheme) Size(n fyne.ThemeSizeName) float32 {
	return fyne.CurrentApp().Settings().Theme().Size(n)
}
