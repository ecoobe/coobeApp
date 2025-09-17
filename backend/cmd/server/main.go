package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Глобальные переменные
var (
	count int
	mu    sync.Mutex
)

// Метрики Prometheus
var (
	requestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "coobeapp_requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"endpoint"})

	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "coobeapp_request_duration_seconds",
		Help:    "Duration of HTTP requests in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"endpoint"})

	goroutinesCount = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "coobeapp_goroutines_total",
		Help: "Current number of goroutines",
	})

	tasksProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "coobeapp_tasks_processed_total",
		Help: "Total number of tasks processed",
	}, []string{"type"})
)

// Функция для имитации "тяжелой" работы
func simulateWork(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

// Функция для последовательной обработки
func processSequential(tasks int, workTime int) {
	for i := 0; i < tasks; i++ {
		simulateWork(workTime)
		mu.Lock()
		count++
		mu.Unlock()
		tasksProcessed.WithLabelValues("sequential").Inc()
	}
}

// Функция для параллельной обработки с использованием горутин
func processParallel(tasks int, workTime int) {
	var wg sync.WaitGroup
	wg.Add(tasks)

	for i := 0; i < tasks; i++ {
		go func() {
			defer wg.Done()
			simulateWork(workTime)
			mu.Lock()
			count++
			mu.Unlock()
			tasksProcessed.WithLabelValues("parallel").Inc()
		}()
	}

	wg.Wait()
}

func main() {
	// Выводим информацию о зарегистрированных обработчиках
	log.Println("Registered handlers:")
	log.Println("  - /api/count")
	log.Println("  - /api/increment")
	log.Println("  - /api/sequential")
	log.Println("  - /api/parallel")
	log.Println("  - /api/reset")
	log.Println("  - /metrics")

	// Инициализация генератора случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Статический файловый сервер
	frontendPath := filepath.Join("frontend")
	http.Handle("/", http.FileServer(http.Dir(frontendPath)))

	// API: получить текущее значение счетчика
	http.HandleFunc("/api/count", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestsTotal.WithLabelValues("/api/count").Inc()

		mu.Lock()
		defer mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"count": count})

		requestDuration.WithLabelValues("/api/count").Observe(time.Since(start).Seconds())
	})

	// API: увеличить счетчик (базовая версия)
	http.HandleFunc("/api/increment", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestsTotal.WithLabelValues("/api/increment").Inc()

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

		requestDuration.WithLabelValues("/api/increment").Observe(time.Since(start).Seconds())
	})

	// API: последовательная обработка
	http.HandleFunc("/api/sequential", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestsTotal.WithLabelValues("/api/sequential").Inc()

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Параметры по умолчанию
		tasks := 10
		workTime := 100

		processSequential(tasks, workTime)

		mu.Lock()
		currentCount := count
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"count":    currentCount,
			"tasks":    tasks,
			"workTime": workTime,
			"duration": time.Since(start).Seconds(),
		})

		requestDuration.WithLabelValues("/api/sequential").Observe(time.Since(start).Seconds())
	})

	// API: параллельная обработка
	http.HandleFunc("/api/parallel", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestsTotal.WithLabelValues("/api/parallel").Inc()

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Параметры по умолчанию
		tasks := 10
		workTime := 100

		processParallel(tasks, workTime)

		mu.Lock()
		currentCount := count
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"count":    currentCount,
			"tasks":    tasks,
			"workTime": workTime,
			"duration": time.Since(start).Seconds(),
		})

		requestDuration.WithLabelValues("/api/parallel").Observe(time.Since(start).Seconds())
	})

	// API: сброс счетчика
	http.HandleFunc("/api/reset", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestsTotal.WithLabelValues("/api/reset").Inc()

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		mu.Lock()
		count = 0
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"count": 0})

		requestDuration.WithLabelValues("/api/reset").Observe(time.Since(start).Seconds())
	})

	// Метрики для Prometheus
	http.Handle("/metrics", promhttp.Handler())

	// Запускаем горутину для мониторинга количества горутин
	go func() {
		for {
			goroutinesCount.Set(float64(runtime.NumGoroutine()))
			time.Sleep(1 * time.Second)
		}
	}()

	// Запуск сервера
	log.Println("➡️  Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
