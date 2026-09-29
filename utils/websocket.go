package utils

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WsMessage mendefinisikan struktur pesan yang akan dikirim ke klien (Admin)
type WsMessage struct {
	Type    string      `json:"type"`    // Contoh: "NEW_BOOKING", "NEW_CONTACT"
	Message string      `json:"message"` // Pesan singkat untuk notifikasi
	Data    interface{} `json:"data"`    // Data tambahan jika diperlukan
}

// Konfigurasi Upgrader untuk mengubah koneksi HTTP biasa menjadi WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Izinkan semua origin (CORS)
	},
}

// WsHub mengelola semua koneksi klien yang aktif
type WsHub struct {
	Clients    map[*websocket.Conn]bool
	Broadcast  chan WsMessage
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
	mu         sync.Mutex // Mencegah race condition saat map diakses bersamaan
}

// Global Hub Instance
var Hub = WsHub{
	Clients:    make(map[*websocket.Conn]bool),
	Broadcast:  make(chan WsMessage),
	Register:   make(chan *websocket.Conn),
	Unregister: make(chan *websocket.Conn),
}

// Run menjalankan loop utama Hub di dalam Goroutine
func (h *WsHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.mu.Unlock()
			log.Println("Admin Terkoneksi ke WebSocket. Total klien:", len(h.Clients))

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				client.Close()
				log.Println("Admin Terputus dari WebSocket. Sisa klien:", len(h.Clients))
			}
			h.mu.Unlock()

		case message := <-h.Broadcast:
			h.mu.Lock()
			// Kirim pesan ke semua klien yang terhubung
			for client := range h.Clients {
				err := client.WriteJSON(message)
				if err != nil {
					log.Println("Error mengirim pesan WebSocket:", err)
					client.Close()
					delete(h.Clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

// ServeWS adalah Handler Gin untuk menerima koneksi WebSocket baru
func ServeWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Gagal Upgrade WebSocket:", err)
		return
	}

	Hub.Register <- conn

	// Tunggu pesan dari klien (bisa dipakai jika Admin mau mengirim balik, tapi di sini kita biarkan saja menunggu sampai koneksi putus)
	go func() {
		defer func() {
			Hub.Unregister <- conn
		}()
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}