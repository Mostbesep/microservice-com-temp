package order

import (
	"context"
	pb "github.com/Mostbesep/microservice-com-temp/order/pb/microservice-com-temp.order.pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.OrderServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := pb.NewOrderServiceClient(conn)
	return &Client{conn, c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostOrder(
	ctx context.Context,
	accountID string,
	products []OrderedProduct,
) (*Order, error) {
	protoProducts := []*pb.PostOrderRequest_OrderProduct{}
	for _, p := range products {
		protoProducts = append(protoProducts, &pb.PostOrderRequest_OrderProduct{
			ProductId: p.Id,
			Quantity:  p.Quantity,
		})
	}
	r, err := c.service.PostOrder(
		ctx,
		&pb.PostOrderRequest{
			AccountId: accountID,
			Products:  protoProducts,
		},
	)
	if err != nil {
		return nil, err
	}

	// Create response order
	newOrder := r.Order
	newOrderCreatedAt := time.Time{}
	err = newOrderCreatedAt.UnmarshalBinary(newOrder.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &Order{
		Id:         newOrder.Id,
		CreatedAt:  newOrderCreatedAt,
		TotalPrice: newOrder.TotalPrice,
		AccountId:  newOrder.AccountId,
		Products:   products,
	}, nil
}

func (c *Client) GetAccountOrders(ctx context.Context, accountID string) ([]Order, error) {
	r, err := c.service.GetAccountOrders(ctx, &pb.GetAccountOrdersRequest{
		AccountId: accountID,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Create response orders
	orders := []Order{}
	for _, orderProto := range r.Orders {
		newOrder := Order{
			Id:         orderProto.Id,
			TotalPrice: orderProto.TotalPrice,
			AccountId:  orderProto.AccountId,
		}
		newOrder.CreatedAt = time.Time{}
		err := newOrder.CreatedAt.UnmarshalBinary(orderProto.CreatedAt)
		if err != nil {
			return nil, err
		}

		products := []OrderedProduct{}
		for _, p := range orderProto.Products {
			products = append(products, OrderedProduct{
				Id:          p.Id,
				Quantity:    p.Quantity,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
		newOrder.Products = products

		orders = append(orders, newOrder)
	}
	return orders, nil
}

func (c *Client) GetOrder(ctx context.Context, id string) (Order, error) {
	r, err := c.service.GetOrder(ctx, &pb.GetOrderRequest{
		Id: id,
	})
	if err != nil {
		log.Println(err)
		return Order{}, err
	}

	// Create response order
	order := Order{
		Id:        id,
		AccountId: r.Order.AccountId,
	}

	for _, product := range r.Order.Products {
		order.Products = append(order.Products, OrderedProduct{
			Id:          product.Id,
			Quantity:    product.Quantity,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}

	return order, nil
}
