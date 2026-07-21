package gui
type TrayManager struct{}
func NewTrayManager(*Window) *TrayManager { return &TrayManager{} }
func (t *TrayManager) Run() {}
func (t *TrayManager) getIcon() []byte { return nil }
