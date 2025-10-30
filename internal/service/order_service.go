package service

import (
	"order-service/internal/cache"
	"order-service/internal/domain"
	"order-service/internal/repository"
)

type OrderService interface {
	CreateOrder(order *domain.Order) error
	GetOrderByUID(uid string) (*domain.Order, error)
	GetAllOrders() ([]domain.Order, error)
	RestoreCache() error
}

type orderService struct {
	repo  repository.OrderRepository
	cache cache.Cache
}

func NewOrderService(repo repository.OrderRepository, cache cache.Cache) OrderService {
	return &orderService{
		repo:  repo,
		cache: cache,
	}
}

func (s *orderService) CreateOrder(order *domain.Order) error {
	if err := s.repo.Save(order); err != nil {
		return err
	}
	
	s.cache.Set(order.OrderUID, *order)
	return nil
}

func (s *orderService) GetOrderByUID(uid string) (*domain.Order, error) {
	// Пробуем получить из кэша
	if order, ok := s.cache.Get(uid); ok {
		return &order, nil
	}

	// Если нет в кэше, ищем в БД
	order, err := s.repo.GetByUID(uid)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	s.cache.Set(uid, *order)
	return order, nil
}

func (s *orderService) GetAllOrders() ([]domain.Order, error) {
	return s.repo.GetAll()
}

func (s *orderService) RestoreCache() error {
	orders, err := s.repo.GetAll()
	if err != nil {
		return err
	}

	for _, order := range orders {
		s.cache.Set(order.OrderUID, order)
	}

	return nil
}