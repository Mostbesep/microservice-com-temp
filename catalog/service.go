package catalog

import (
	"context"
	"github.com/segmentio/ksuid"
)

type Product struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type Service interface {
	PostProduct(ctx context.Context, name, description string, price float64) (Product, error)
	GetProduct(ctx context.Context, productID string) (Product, error)
	ListProducts(ctx context.Context, skip uint64, take uint64) (*[]Product, error)
	ListProductsByIDs(ctx context.Context, productIDs []string) (*[]Product, error)
	SearchProducts(ctx context.Context, query string, skip uint64, take uint64) (*[]Product, error)
}

type catalogService struct {
	repository Repository
}

func (c *catalogService) PostProduct(ctx context.Context, name, description string, price float64) (Product, error) {
	newProduct := Product{
		Id:          ksuid.New().String(),
		Name:        name,
		Description: description,
		Price:       price,
	}

	err := c.repository.PutProduct(ctx, newProduct)
	if err != nil {
		return Product{}, err
	}
	return newProduct, nil
}

func (c *catalogService) GetProduct(ctx context.Context, productID string) (Product, error) {
	return c.GetProduct(ctx, productID)
}

func (c *catalogService) ListProducts(ctx context.Context, skip uint64, take uint64) (*[]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return c.ListProducts(ctx, skip, take)
}

func (c *catalogService) ListProductsByIDs(ctx context.Context, productIDs []string) (*[]Product, error) {
	return c.ListProductsByIDs(ctx, productIDs)
}

func (c catalogService) SearchProducts(ctx context.Context, query string, skip uint64, take uint64) (*[]Product, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return c.SearchProducts(ctx, query, skip, take)
}

func NewService(repository Repository) Service {
	return &catalogService{repository: repository}
}
