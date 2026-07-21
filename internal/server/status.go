package server

import "sync"

type ServerStatus int

const (
	StatusOffline ServerStatus = iota
	StatusOnline
	StatusStarting
	StatusStopping
)

type StatusManager struct {
	mu     sync.RWMutex
	status ServerStatus
}

func NewStatusManager() *StatusManager {
	return &StatusManager{
		status: StatusOffline,
	}
}

func (s *StatusManager) SetStatus(status ServerStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status = status
}

func (s *StatusManager) GetStatus() ServerStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.status
}

func (s *StatusManager) IsOnline() bool {
	return s.GetStatus() == StatusOnline
}