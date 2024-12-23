package order

import (
	"context"
	"github.com/segmentio/ksuid"
	"time"
)

type Service interface {
	PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (Order, error)
	GetOrder(ctx context.Context, id string) (Order, error)
	GetAccountOrders(ctx context.Context, accountID string) (*[]Order, error)
}

type Order struct {
	Id         string
	CreatedAt  time.Time
	TotalPrice float64
	AccountId  string
	Products   []OrderedProduct
}

type OrderedProduct struct {
	Id          string
	Name        string
	Description string
	Price       float64
	Quantity    uint32
}

type ReceivedProduct struct {
	Id       string
	Quantity uint32
}

type orderService struct {
	repository Repository
}

func (s *orderService) GetOrder(ctx context.Context, id string) (Order, error) {
	return s.repository.GetOrder(ctx, id)
}

func (s *orderService) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (Order, error) {
	o := Order{
		Id:        ksuid.New().String(),
		CreatedAt: time.Now().UTC(),
		AccountId: accountID,
		Products:  products,
	}
	// Calculate total price
	o.TotalPrice = 0.0
	for _, p := range products {
		o.TotalPrice += p.Price * float64(p.Quantity)
	}
	err := s.repository.PutOrder(ctx, o)
	if err != nil {
		return Order{}, err
	}
	return o, nil
}

func (s *orderService) GetAccountOrders(ctx context.Context, accountID string) (*[]Order, error) {
	return s.repository.GetAccountOrders(ctx, accountID)
}

func NewService(repository Repository) Service {
	return &orderService{repository: repository}
}
