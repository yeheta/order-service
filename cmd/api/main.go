package main

import (
	"log"
	"net/http"
	"order-service/internal/cache"
	"order-service/internal/config"
	"order-service/internal/handler"
	"order-service/internal/nats"
	"order-service/internal/repository"
	"order-service/internal/service"

	"github.com/gorilla/mux"
)

func main() {
	// Загрузка конфигурации
	cfg := config.Load()

	// Инициализация репозитория
	repo, err := repository.NewPostgresRepository(
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
	)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer repo.Close()

	// Инициализация кэша
	cache := cache.NewMemoryCache()

	// Инициализация сервиса
	orderService := service.NewOrderService(repo, cache)

	// Восстановление кэша из БД
	log.Println("Restoring cache from database...")
	if err := orderService.RestoreCache(); err != nil {
		log.Fatal("Failed to restore cache from DB:", err)
	}
	log.Printf("Cache restored. Total orders: %d", cache.Size())

	// Инициализация NATS подписчика
	subscriber, err := nats.NewSubscriber(
		cfg.NATS.ClusterID,
		cfg.NATS.ClientID,
		cfg.NATS.URL,
		orderService,
	)
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	defer subscriber.Close()

	// Подписка на канал
	if err := subscriber.Subscribe("orders"); err != nil {
		log.Fatal("Failed to subscribe to NATS channel:", err)
	}
	log.Println("Subscribed to NATS channel: orders")

	// Инициализация HTTP обработчиков
	handler := handler.NewHandler(orderService)

	// Настройка маршрутов
	router := mux.NewRouter()

	// Добавляем CORS middleware
	router.Use(corsMiddleware)

	// API routes
	router.HandleFunc("/api/order/{id}", handler.GetOrder).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/orders", handler.GetOrders).Methods("GET", "OPTIONS")
	router.HandleFunc("/health", handler.HealthCheck).Methods("GET")

	// Запуск HTTP сервера
	log.Printf("Server starting on port %s", cfg.HTTP.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.HTTP.Port, router))
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}