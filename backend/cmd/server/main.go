package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

// ====== Глобальные переменные ======
// count — текущее значение счётчика
// mu — мьютекс для безопасного доступа из разных запросов
var (
	count int
	mu    sync.Mutex
)

func main() {
	// ====== Обработка фронтенда ======
	// Берём путь до папки с index.html (лежит в ../frontend относительно backend/cmd/server)
	frontendPath := filepath.Join("..", "..", "frontend")

	// FileServer умеет отдавать статические файлы (HTML, JS, CSS)
	fs := http.FileServer(http.Dir(frontendPath))
	// Все запросы по "/" будут искать файлы в frontendPath
	http.Handle("/", fs)

	// ====== API: получить текущее значение ======
	http.HandleFunc("/api/count", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		// Ответ в JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"count": count})
	})

	// ====== API: увеличить счётчик ======
	http.HandleFunc("/api/increment", func(w http.ResponseWriter, r *http.Request) {
		// Только POST-запрос
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Увеличиваем счётчик с блокировкой
		mu.Lock()
		count++
		newValue := count
		mu.Unlock()

		// Возвращаем новое значение в JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"count": newValue})
	})

	// ====== Запуск сервера ======
	log.Println("➡️  Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}