package api

import (
	"fmt"
	"encoding/json"
	"net/http"
	"math/rand"
	"log"
	"time"
	"taxi-backend/internal/models"
	"taxi-backend/internal/storage"
	"github.com/google/uuid" 
)

type Handler struct {
	Store storage.Storage
}

func NewHandler(store storage.Storage) *Handler {
	return &Handler{
		Store: store,
	}
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid JSON received", http.StatusBadRequest)
		return
	}

	newID := uuid.NewString()

	finalPrice := 300.0
	if req.Tariff == "comfort" {
		finalPrice += 200.0
	}
	finalPrice += float64(len(req.SelectedServices)) * 50.0

	newOrder := models.Order{
		ID:            newID,
		PassengerID:   req.PassengerID,
		PassengerType: req.PassengerType, 
		Status:        models.StatusSearching, 
		AddressFrom:   req.AddressFrom,
		AddressTo:     req.AddressTo,
		Tariff:        req.Tariff,
		Price:         finalPrice, 
		CreatedAt:     time.Now(),
	}

	if err := h.Store.SaveOrder(newOrder); err != nil {
		log.Printf("Error saving order: %v", err)
		http.Error(w, "Failed to save order", http.StatusInternalServerError)
		return
	}

	log.Printf("Order created successfully. ID = %s, Passenger = %s", newOrder.ID, newOrder.PassengerID)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"order_id": newOrder.ID,
		"status": newOrder.Status,
		"price": newOrder.Price,
	})

}

func(h *Handler) GetOrder(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet {
			http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
			return
		}

		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "Order ID is required", http.StatusBadRequest)
			return
		}

		order,err := h.Store.GetOrder(id)
		if err != nil {
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
	
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(order)
}

func(h *Handler) AcceptOrder(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	var req models.AcceptOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	order, err := h.Store.GetOrder(id)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	if order.Status != models.StatusPending && order.Status != models.StatusSearching {
		http.Error(w, "Order is not already taken or canceled", http.StatusConflict)
		return
	}

	order.Status = models.StatusAssigned
	order.DriverID = &req.DriverID
	if err := h.Store.UpdateOrder(order); err != nil {
		log.Printf("Error updating order: %v", err)
		http.Error(w, "Failed to update order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Order accepted successfully",
		"status": order.Status,
	})
}

func (h *Handler) DriverArrived(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}
	
	order, err := h.Store.GetOrder(id)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
		if order.Status != models.StatusAssigned {
		http.Error(w, "Order status is invalid for this action", http.StatusBadRequest)
		return
	}

	code := fmt.Sprintf("%05d", rand.Intn(100000))
	order.ConfirmationCode = code
	order.Status = models.StatusArrived
	if err := h.Store.UpdateOrder(order); err != nil {
		http.Error(w, "Failed to update order", http.StatusInternalServerError)
		return
	}
    log.Printf("Driver arrived for order %s. Code: %s", order.ID, code)
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":           "Driver arrived",
		"status":            order.Status,
		"confirmation_code": code,
	})
}

func (h *Handler) ChangeOrderStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	
	var req models.UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	order, err := h.Store.GetOrder(id)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	order.Status = req.Status

	if err := h.Store.UpdateOrder(order); err != nil {
		http.Error(w, "Failed to update", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": order.Status,
	})
}

func (h *Handler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	order, err := h.Store.GetOrder(id)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	if order.Status == models.StatusInProgress || order.Status == models.StatusCompleted {
		http.Error(w, "Cannot cancel order in progress or completed", http.StatusBadRequest)
		return
	}

	order.Status = models.StatusCancelled
	if err := h.Store.UpdateOrder(order); err != nil {
		http.Error(w, "Failed to update order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Order canceled successfully",
		"status":  order.Status,
	})
}

func (h *Handler) UpdateDriverStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method allowed", http.StatusMethodNotAllowed)
		return
	}

	driverID := r.Header.Get("X-Driver-ID")
	if driverID == "" {
		driverID = "unknown_driver"
	}

	var req models.UpdateDriverStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	statusStr := "unavailable"
	if req.IsAvailable {
		statusStr = "available"
	}

	log.Printf("Driver %s is %s at [%f, %f]", driverID, statusStr, req.CurrentLocation.Lat, req.CurrentLocation.Lng)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "updated",
	})
}

func (h *Handler) GetAvailableOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method allowed", http.StatusMethodNotAllowed)
		return
	}

	orders, err := h.Store.GetPendingOrders()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}