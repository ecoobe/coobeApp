package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

var (
	count int
	mu    sync.Mutex
)

func main() {
	// Отдаём статику (фронтенд)
	http.Handle("/", http.FileServer(http.Dir("./frontend")))

	// API: получить текущее значение
	http.HandleFunc("/api/count", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"count": count})
	})

	// API: увеличить счётчик
	http.HandleFunc("/api/increment", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		mu.Lock()
		count++
		newValue := count
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"count": newValue})
	})

	log.Println("➡️  Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}