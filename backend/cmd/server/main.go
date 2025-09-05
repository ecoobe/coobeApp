package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Thought struct {
	ID        int64     `json:"id"`
	Text      string    `json:"text"`
	Lat       float64   `json:"lat"`
	Lng       float64   `json:"lng"`
	CreatedAt time.Time `json:"created_at"`
}

type Hub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte, 128),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (h *Hub) run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c] = true
			h.mu.Unlock()
		case c := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				c.Close()
			}
			h.mu.Unlock()
		case msg := <-h.broadcast:
			h.mu.Lock()
			for c := range h.clients {
				// best-effort write (no per-client goroutine to keep example compact)
				_ = c.WriteMessage(websocket.TextMessage, msg)
			}
			h.mu.Unlock()
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // dev-only: разрешаем все origin
}

func wsHandler(h *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("ws upgrade:", err)
			return
		}
		h.register <- conn

		// читаем и игнорируем входящие сообщения (поддерживаем соединение)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				h.unregister <- conn
				return
			}
		}
	}
}

// простое in-memory хранилище (MVP)
var (
	storeMu  sync.Mutex
	thoughts []Thought
	nextID   int64 = 1
)

func postThoughtHandler(h *Hub) http.HandlerFunc {
	type Req struct {
		Text string  `json:"text"`
		Lat  float64 `json:"lat"`
		Lng  float64 `json:"lng"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// простая CORS + content-type
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		var req Req
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		storeMu.Lock()
		t := Thought{
			ID:        nextID,
			Text:      req.Text,
			Lat:       req.Lat,
			Lng:       req.Lng,
			CreatedAt: time.Now().UTC(),
		}
		thoughts = append(thoughts, t)
		nextID++
		storeMu.Unlock()

		// пушим всем WS-клиентам (тип + данные)
		msg := struct {
			Type string  `json:"type"`
			Data Thought `json:"data"`
		}{
			Type: "new_thought",
			Data: t,
		}
		b, _ := json.Marshal(msg)
		select {
		case h.broadcast <- b:
		default:
			// если очередь полна — пропускаем (best-effort)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(t)
	}
}

func getThoughtsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	storeMu.Lock()
	defer storeMu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(thoughts)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	hub := newHub()
	go hub.run()

	http.HandleFunc("/api/health", healthHandler)
	http.HandleFunc("/api/thoughts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			postThoughtHandler(hub)(w, r)
			return
		}
		if r.Method == http.MethodGet {
			getThoughtsHandler(w, r)
			return
		}
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	http.HandleFunc("/ws", wsHandler(hub))

	addr := ":" + port
	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	log.Printf("server listening on %s", addr)
	log.Fatal(srv.ListenAndServe())
}
