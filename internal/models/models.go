package models

import "time"

const (
StatusPending = "pending"
StatusSearching = "searching"
StatusAssigned = "driver_assigned"
StatusArrived = "waiting_for_confirmation"
StatusInProgress = "in_progress"
StatusCompleted = "completed"
StatusCancelled = "cancelled"
)

type CreateOrderRequest struct {
	PassengerID string `json:"passenger_id"`
	PassengerType string `json:"passenger_type"`
	AddressFrom string `json:"address_from"`
	AddressTo string `json:"address_to"`
	Tariff string `json:"tariff"`
	SelectedServices []string `json:"selected_services"`
	Comments string `json:"comments"`
}

type Order struct {
	ID string `json:"id"`
	PassengerID string `json:"passenger_id"`
	PassengerType string `json:"passenger_type"`
	DriverID *string `json:"driver_id"`
	Status string `json:"status"`	
	AddressFrom string `json:"address_from"`	
	AddressTo string `json:"address_to"`
	Tariff string `json:"tariff"`
	Price float64 `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	ConfirmationCode string `json:"-"`
}
	
type AcceptOrderRequest struct {
	DriverID string `json:"driver_id"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

type DriverLocation struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type UpdateDriverStatusRequest struct {
	IsAvailable     bool           `json:"is_available"`
	CurrentLocation DriverLocation `json:"current_location"`
}
	
	