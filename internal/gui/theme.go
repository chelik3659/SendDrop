package gui

import (
	"image/color"
	
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type CustomTheme struct {
	darkMode bool
}

func NewCustomTheme(dark bool) *CustomTheme {
	return &CustomTheme{darkMode: dark}
}

func (t *CustomTheme) Color(c fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	if t.darkMode {
		switch c {
		case theme.ColorNameBackground:
			return color.NRGBA{R: 30, G: 30, B: 30, A: 255}
		case theme.ColorNameForeground:
			return color.White
		default:
			return theme.DarkTheme().Color(c, v)
		}
	}
	return theme.LightTheme().Color(c, v)
}

func (t *CustomTheme) Font(s fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(s)
}

func (t *CustomTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (t *CustomTheme) Size(n fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(n)
}
