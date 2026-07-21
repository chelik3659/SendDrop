package gui

type StatusBar struct{}
func NewStatusBar() *StatusBar { return &StatusBar{} }
func (s *StatusBar) SetOnline(bool) {}
func (s *StatusBar) SetIP(string) {}
func (s *StatusBar) SetPort(int) {}

type QRWidget struct{}
func NewQRWidget() *QRWidget { return &QRWidget{} }
func (q *QRWidget) SetQR([]byte) {}
