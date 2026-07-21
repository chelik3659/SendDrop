package gui

type TrayManager struct {
	window *Window
}

func NewTrayManager(window *Window) *TrayManager {
	return &TrayManager{window: window}
}

func (t *TrayManager) Run() {}

func (t *TrayManager) getIcon() []byte {
	return nil
}
