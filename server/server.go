package server

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rdmnl/nexora/services"
	"github.com/rdmnl/nexora/shared"
)

var (
	clients     = make(map[*websocket.Conn]bool)
	broadcast   = make(chan []shared.ServerInfo)
	upgrader    = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	servers     []shared.ServerInfo
	mu          sync.Mutex
	configNodes []shared.ServerInfo
)

func Init(config *shared.Config) {
	configNodes = config.Nodes
	servers = services.DetectAndMerge([]shared.ServerInfo{}, configNodes)

	fs := http.FileServer(http.Dir("./templates"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", streamServerData)
}

func StartHTTPServer(addr string) *http.Server {
	srv := &http.Server{Addr: addr}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()
	return srv
}

func streamServerData(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	if err := conn.WriteJSON(servers); err != nil {
		log.Printf("Error sending initial data: %v", err)
		removeClient(conn)
		return
	}

	for data := range broadcast {
		if err := conn.WriteJSON(data); err != nil {
			log.Printf("Error sending data: %v", err)
			removeClient(conn)
			break
		}
	}
}

func StartBroadcast() {
	for {
		servers = services.DetectAndMerge(servers, configNodes)
		services.UpdateServerUsage(servers)

		broadcast <- servers
		time.Sleep(5 * time.Second)
	}
}

func removeClient(conn *websocket.Conn) {
	mu.Lock()
	delete(clients, conn)
	mu.Unlock()
}
