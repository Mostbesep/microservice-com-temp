//go:generate protoc --go_out=plugins=grpc=./pb order.proto
package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/Mostbesep/microservice-com-temp/account"
	"github.com/Mostbesep/microservice-com-temp/catalog"
	pb "github.com/Mostbesep/microservice-com-temp/order/pb/microservice-com-temp.order.pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type grpcServer struct {
	service       Service
	accountClient *account.Client
	catalogClient *catalog.Client
	pb.UnimplementedOrderServiceServer
}

func ListenGRPC(s Service, accountURL, catalogURL string, port int) error {
	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		return err
	}

	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		accountClient.Close()
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return err
	}

	serv := grpc.NewServer()
	pb.RegisterOrderServiceServer(serv, &grpcServer{
		service:       s,
		accountClient: accountClient,
		catalogClient: catalogClient,
	})
	reflection.Register(serv)

	return serv.Serve(lis)
}

func (s *grpcServer) PostOrder(
	ctx context.Context,
	r *pb.PostOrderRequest,
) (*pb.PostOrderResponse, error) {
	// Check if account exists
	_, err := s.accountClient.GetAccount(ctx, r.AccountId)
	if err != nil {
		log.Println("Error getting account: ", err)
		return nil, errors.New("account not found")
	}

	// Get ordered products
	productIDs := []string{}
	for _, p := range r.Products {
		productIDs = append(productIDs, p.ProductId)
	}
	orderedProducts, err := s.catalogClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Println("Error getting products: ", err)
		return nil, errors.New("products not found")
	}

	// Construct products
	products := []OrderedProduct{}
	for _, p := range orderedProducts {
		product := OrderedProduct{
			Id:          p.Id,
			Quantity:    0,
			Price:       p.Price,
			Name:        p.Name,
			Description: p.Description,
		}
		for _, rp := range r.Products {
			if rp.ProductId == p.Id {
				product.Quantity = rp.Quantity
				break
			}
		}

		if product.Quantity != 0 {
			products = append(products, product)
		}
	}

	// Call service implementation
	order, err := s.service.PostOrder(ctx, r.AccountId, products)
	if err != nil {
		log.Println("Error posting order: ", err)
		return nil, errors.New("could not post order")
	}

	// Make response order
	orderProto := &pb.Order{
		Id:         order.Id,
		AccountId:  order.AccountId,
		TotalPrice: order.TotalPrice,
		Products:   []*pb.Order_OrderProduct{},
	}
	orderProto.CreatedAt, _ = order.CreatedAt.MarshalBinary()
	for _, p := range order.Products {
		orderProto.Products = append(orderProto.Products, &pb.Order_OrderProduct{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
		})
	}
	return &pb.PostOrderResponse{
		Order: orderProto,
	}, nil
}

func (s *grpcServer) GetOrder(ctx context.Context, request *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	// Fetch the order using the service layer
	order, err := s.service.GetOrder(ctx, request.Id)
	if err != nil {
		log.Println("Error fetching order:", err)
		return nil, err
	}

	productIDs := []string{}
	for _, p := range order.Products {
		productIDs = append(productIDs, p.Id)
	}

	// Fetch product details from the catalog service
	products, err := s.catalogClient.GetProducts(ctx, 0, 0, productIDs, "")
	if err != nil {
		log.Println("Error fetching products for order:", err)
		return nil, err
	}

	// Construct the response order
	responseOrder := &pb.Order{
		AccountId:  order.AccountId,
		Id:         order.Id,
		TotalPrice: order.TotalPrice,
		Products:   []*pb.Order_OrderProduct{},
	}
	responseOrder.CreatedAt, _ = order.CreatedAt.MarshalBinary()

	orderProducts := map[string]catalog.Product{}
	for _, product := range products {
		orderProducts[product.Id] = product
	}

	for _, product := range order.Products {
		responseOrder.Products = append(responseOrder.Products, &pb.Order_OrderProduct{
			Id:          product.Id,
			Name:        orderProducts[product.Id].Name,
			Description: orderProducts[product.Id].Description,
			Price:       orderProducts[product.Id].Price,
			Quantity:    product.Quantity,
		})
	}

	return &pb.GetOrderResponse{Order: responseOrder}, nil
}

func (s *grpcServer) GetAccountOrders(
	ctx context.Context,
	r *pb.GetAccountOrdersRequest,
) (*pb.GetAccountOrdersResponse, error) {

	// Get orders for account
	accountOrders, err := s.service.GetAccountOrders(ctx, r.AccountId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	productIDs := []string{}
	for _, order := range *accountOrders {
		for _, p := range order.Products {
			productIDs = append(productIDs, p.Id)
		}
	}
	products, err := s.catalogClient.GetProducts(ctx, 0, 0, productIDs, "")

	productsMap := make(map[string]catalog.Product)
	for _, product := range products {
		productsMap[product.Id] = product
	}
	var responseOrders []*pb.Order

	for _, order := range *accountOrders {
		newOrder := &pb.Order{
			AccountId:  order.AccountId,
			Id:         order.Id,
			TotalPrice: order.TotalPrice,
		}
		for _, product := range order.Products {
			orderProduct := &pb.Order_OrderProduct{
				Id:          product.Id,
				Name:        productsMap[product.Id].Name,
				Description: productsMap[product.Id].Description,
				Price:       productsMap[product.Id].Price,
				Quantity:    product.Quantity,
			}
			newOrder.Products = append(newOrder.Products, orderProduct)
		}
		responseOrders = append(responseOrders)
	}
	return &pb.GetAccountOrdersResponse{Orders: responseOrders}, nil
}
