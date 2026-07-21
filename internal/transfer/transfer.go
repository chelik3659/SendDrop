package transfer

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const transferPort = 8082

type Transfer struct {
	onFileReceive func(filename string, data []byte)
}

func NewTransfer() *Transfer {
	return &Transfer{}
}

func (t *Transfer) SetOnFileReceive(cb func(string, []byte)) {
	t.onFileReceive = cb
}

// SendFile отправляет файл на указанный IP
func (t *Transfer) SendFile(ip string, filename string, data []byte) error {
	conn, err := net.DialTimeout("tcp", ip+":"+fmt.Sprint(transferPort), 10*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Отправляем заголовок: SEND_FILE|имя_файла|размер\n
	header := fmt.Sprintf("SEND_FILE|%s|%d\n", filename, len(data))
	if _, err := conn.Write([]byte(header)); err != nil {
		return err
	}
	// Отправляем данные
	if _, err := conn.Write(data); err != nil {
		return err
	}
	return nil
}

// StartServer запускает TCP сервер для приёма файлов
func (t *Transfer) StartServer(shareDir string) {
	listener, err := net.Listen("tcp", ":"+fmt.Sprint(transferPort))
	if err != nil {
		fmt.Println("Transfer server error:", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go t.handleConnection(conn, shareDir)
	}
}

func (t *Transfer) handleConnection(conn net.Conn, shareDir string) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// Читаем заголовок до \n
	header, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	header = strings.TrimSpace(header)
	if !strings.HasPrefix(header, "SEND_FILE|") {
		return
	}

	parts := strings.Split(strings.TrimPrefix(header, "SEND_FILE|"), "|")
	if len(parts) != 2 {
		return
	}
	filename := parts[0]
	size, err := strconv.Atoi(parts[1])
	if err != nil {
		return
	}

	// Читаем ровно size байт
	data := make([]byte, size)
	_, err = io.ReadFull(reader, data)
	if err != nil {
		return
	}

	// Сохраняем файл
	fullPath := filepath.Join(shareDir, filename)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err == nil {
		os.WriteFile(fullPath, data, 0644)
	}

	if t.onFileReceive != nil {
		t.onFileReceive(filename, data)
	}
}
