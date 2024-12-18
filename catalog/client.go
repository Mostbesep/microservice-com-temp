package catalog

import (
	"context"
	pb "github.com/Mostbesep/microservice-com-temp/catalog/pb/microservice-com-temp.catalog.pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	Service pb.CatalogServiceClient
}

func NewClient(url string) (*Client, error) {

	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := pb.NewCatalogServiceClient(conn)
	return &Client{conn: conn, Service: c}, nil
}

func (c *Client) Close() error {
	err := c.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) PostProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	r, err := c.Service.PostProduct(
		ctx,
		&pb.PostProductRequest{
			Name:        name,
			Description: description,
			Price:       price,
		},
	)
	if err != nil {
		return nil, err
	}
	return &Product{
		Id:          r.Product.Id,
		Name:        r.Product.Name,
		Description: r.Product.Description,
		Price:       r.Product.Price,
	}, nil
}

func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
	r, err := c.Service.GetProduct(
		ctx,
		&pb.GetProductRequest{
			Id: id,
		},
	)
	if err != nil {
		return nil, err
	}

	return &Product{
		Id:          r.Product.Id,
		Name:        r.Product.Name,
		Description: r.Product.Description,
		Price:       r.Product.Price,
	}, nil
}

func (c *Client) GetProducts(ctx context.Context, skip uint64, take uint64, ids []string, query string) ([]Product, error) {
	r, err := c.Service.GetProducts(
		ctx,
		&pb.GetProductsRequest{
			Ids:   ids,
			Skip:  skip,
			Take:  take,
			Query: query,
		},
	)
	if err != nil {
		return nil, err
	}
	products := []Product{}
	for _, p := range r.Products {
		products = append(products, Product{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		})
	}
	return products, nil
}
