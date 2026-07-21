package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	status     *StatusManager
	port       int
	ip         string
	handler    *Handler
	ctx        context.Context
	cancel     context.CancelFunc
}

func NewServer(port int, ip string, shareDir string) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	
	s := &Server{
		status:  NewStatusManager(),
		port:    port,
		ip:      ip,
		ctx:     ctx,
		cancel:  cancel,
		handler: NewHandler(shareDir),
	}
	
	s.handler.SetServer(s)
	return s
}

func (s *Server) Start() error {
	s.status.SetStatus(StatusStarting)
	
	addr := fmt.Sprintf("%s:%d", s.ip, s.port)
	
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.handler,
		ReadTimeout:  30 * time.Minute,
		WriteTimeout: 30 * time.Minute,
	}
	
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.status.SetStatus(StatusOffline)
		}
	}()
	
	s.status.SetStatus(StatusOnline)
	return nil
}

func (s *Server) Stop() error {
	s.status.SetStatus(StatusStopping)
	
	if s.httpServer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.httpServer.Shutdown(ctx); err != nil {
			return err
		}
	}
	
	s.status.SetStatus(StatusOffline)
	return nil
}

func (s *Server) GetURL() string {
	ip := s.GetLocalIP()
	return fmt.Sprintf("http://%s:%d", ip, s.port)
}

func (s *Server) GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func (s *Server) GetStatus() ServerStatus {
	return s.status.GetStatus()
}

func (s *Server) IsOnline() bool {
	return s.status.IsOnline()
}

func (s *Server) GetFileURL(filename string) string {
	ip := s.GetLocalIP()
	return fmt.Sprintf("http://%s:%d/api/download/%s", ip, s.port, filename)
}