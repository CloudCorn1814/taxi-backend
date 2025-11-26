package main

import (
	"log"
	"net/http"
	"taxi-backend/internal/config"
	"taxi-backend/internal/storage"
	"taxi-backend/internal/api"
)

func main() {
	cfg := config.LoadConfig()
	log.Println("Config loaded, server starts on port:", cfg.Port)
	log.Printf("DB Config: %s", cfg.DatabaseURL)

	memoryStore := storage.NewMemoryStorage()

	handler := api.NewHandler(memoryStore)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("POST /api/order", handler.CreateOrder)
	mux.HandleFunc("GET /api/order/{id}", handler.GetOrder)
	mux.HandleFunc("POST /api/order/{id}/accept", handler.AcceptOrder)
	mux.HandleFunc("POST /api/order/{id}/arrived", handler.DriverArrived)
	mux.HandleFunc("POST /api/order/{id}/status", handler.ChangeOrderStatus)
	mux.HandleFunc("POST /api/order/{id}/cancel", handler.CancelOrder)
	mux.HandleFunc("POST /api/driver/status", handler.UpdateDriverStatus)
	mux.HandleFunc("GET /api/orders/available", handler.GetAvailableOrders)
	
	log.Println("Starting server...")
	
	if err := http.ListenAndServe(cfg.Port, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
