package discovery

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	broadcastPort = 9999
	discoveryMsg  = "SENDDROP_DISCOVERY"
	responsePrefix = "SENDDROP_RESPONSE|"
)

type Peer struct {
	IP   string
	Name string
}

type Discovery struct {
	peers      map[string]Peer
	mu         sync.RWMutex
	onPeerAdd  func(Peer)
	onPeerRemove func(Peer)
	localIP    string
	peerName   string
}

func NewDiscovery(peerName string) *Discovery {
	d := &Discovery{
		peers:    make(map[string]Peer),
		localIP:  getLocalIP(),
		peerName: peerName,
	}
	return d
}

func (d *Discovery) SetCallbacks(onAdd func(Peer), onRemove func(Peer)) {
	d.onPeerAdd = onAdd
	d.onPeerRemove = onRemove
}

func (d *Discovery) Start() {
	go d.broadcastLoop()
	go d.listenLoop()
}

// broadcastLoop отправляет UDP broadcast каждые 3 секунды
func (d *Discovery) broadcastLoop() {
	conn, err := net.Dial("udp", "255.255.255.255:"+fmt.Sprint(broadcastPort))
	if err != nil {
		fmt.Println("Broadcast error:", err)
		return
	}
	defer conn.Close()

	msg := []byte(discoveryMsg)
	for {
		conn.Write(msg)
		time.Sleep(3 * time.Second)
	}
}

// listenLoop слушает UDP ответы
func (d *Discovery) listenLoop() {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", broadcastPort))
	if err != nil {
		fmt.Println("Listen error:", err)
		return
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("ListenUDP error:", err)
		return
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		msg := string(buf[:n])
		if msg == discoveryMsg {
			// Кто-то ищет устройства — отвечаем
			response := responsePrefix + d.localIP + "|" + d.peerName
			conn.WriteToUDP([]byte(response), remoteAddr)
		} else if strings.HasPrefix(msg, responsePrefix) {
			// Получили ответ от другого устройства
			parts := strings.Split(strings.TrimPrefix(msg, responsePrefix), "|")
			if len(parts) >= 2 {
				ip := parts[0]
				name := parts[1]
				if ip != d.localIP {
					d.addPeer(ip, name)
				}
			}
		}
	}
}

func (d *Discovery) addPeer(ip, name string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, exists := d.peers[ip]; !exists {
		peer := Peer{IP: ip, Name: name}
		d.peers[ip] = peer
		if d.onPeerAdd != nil {
			go d.onPeerAdd(peer)
		}
	}
}

func (d *Discovery) GetPeers() []Peer {
	d.mu.RLock()
	defer d.mu.RUnlock()
	peers := make([]Peer, 0, len(d.peers))
	for _, p := range d.peers {
		peers = append(peers, p)
	}
	return peers
}

func getLocalIP() string {
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
