package nats

import (
	"encoding/json"
	"log"
	"order-service/internal/domain"
	"order-service/internal/service"

	"github.com/nats-io/stan.go"
)

type Subscriber struct {
	stanConn      stan.Conn
	orderService  service.OrderService
	subscription  stan.Subscription
}

func NewSubscriber(clusterID, clientID, url string, orderService service.OrderService) (*Subscriber, error) {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(url))
	if err != nil {
		return nil, err
	}

	return &Subscriber{
		stanConn:     sc,
		orderService: orderService,
	}, nil
}

func (s *Subscriber) Subscribe(channel string) error {
	sub, err := s.stanConn.Subscribe(channel, func(msg *stan.Msg) {
		if err := s.processMessage(msg.Data); err != nil {
			log.Printf("Error processing message: %v", err)
		}
	})
	if err != nil {
		return err
	}

	s.subscription = sub
	return nil
}

func (s *Subscriber) processMessage(data []byte) error {
	var order domain.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return err
	}

	// Валидация обязательных полей
	if order.OrderUID == "" || order.TrackNumber == "" {
		return nil // Игнорируем некорректные сообщения
	}

	if err := s.orderService.CreateOrder(&order); err != nil {
		return err
	}

	log.Printf("Order %s processed and cached", order.OrderUID)
	return nil
}

func (s *Subscriber) Close() {
	if s.subscription != nil {
		s.subscription.Unsubscribe()
	}
	s.stanConn.Close()
}