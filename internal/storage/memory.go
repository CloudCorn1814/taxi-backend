package storage

import (
	"sync"
	"taxi-backend/internal/models"
)

type MemoryStorage struct {
	data map[string]models.Order
	mu   sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]models.Order),
	}
}

func (s *MemoryStorage) SaveOrder(order models.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[order.ID] = order
	return nil
}

func (s *MemoryStorage) GetOrder(id string) (models.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.data[id]
	if !ok {
		return models.Order{}, ErrNotFound
	}
	return order, nil
}

func (s *MemoryStorage) UpdateOrder(order models.Order) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _,exists := s.data[order.ID]; !exists {
		return ErrNotFound
	}
	
	s.data[order.ID] = order
	return nil
}

func (s *MemoryStorage) GetPendingOrders() ([]models.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []models.Order

	for _, order := range s.data {
		if order.Status == models.StatusPending || order.Status == models.StatusSearching {
			result = append(result, order)
		}
	}
	return result, nil
}

