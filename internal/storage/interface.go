package storage

import (
"taxi-backend/internal/models"
"errors"
)

type Storage interface {
	SaveOrder(order models.Order) error
	GetOrder(id string) (models.Order, error)
	UpdateOrder(order models.Order) error
	GetPendingOrders() ([]models.Order, error) 
}

var ErrNotFound = errors.New("Order not found")