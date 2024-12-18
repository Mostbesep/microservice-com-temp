//go:generate protoc --go_out=plugins=grpc=./pb catalog.proto
package catalog

import (
	"context"
	"fmt"
	pb "github.com/Mostbesep/microservice-com-temp/catalog/pb/microservice-com-temp.catalog.pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type grpcServer struct {
	service Service
	pb.UnimplementedCatalogServiceServer
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	pb.RegisterCatalogServiceServer(server, &grpcServer{
		UnimplementedCatalogServiceServer: pb.UnimplementedCatalogServiceServer{},
		service:                           s,
	})
	reflection.Register(server)
	return server.Serve(lis)
}

func (s *grpcServer) PostProduct(ctx context.Context, r *pb.PostProductRequest) (*pb.PostProductResponse, error) {
	p, err := s.service.PostProduct(ctx, r.Name, r.Description, r.Price)
	if err != nil {
		return nil, err
	}
	return &pb.PostProductResponse{
		Product: &pb.Product{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		},
	}, nil
}

func (s *grpcServer) GetProduct(ctx context.Context, r *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	p, err := s.service.GetProduct(ctx, r.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetProductResponse{
		Product: &pb.Product{
			Id:          p.Id,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		},
	}, nil
}

func (s *grpcServer) GetProducts(ctx context.Context, r *pb.GetProductsRequest) (*pb.GetProductsResponse, error) {
	var res *[]Product
	var err error
	if r.Query != "" {
		res, err = s.service.SearchProducts(ctx, r.Query, r.Skip, r.Take)
	} else if len(r.Ids) != 0 {
		res, err = s.service.ListProductsByIDs(ctx, r.Ids)
	} else {
		res, err = s.service.ListProducts(ctx, r.Skip, r.Take)
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var products []*pb.Product
	for _, p := range *res {
		products = append(
			products,
			&pb.Product{
				Id:          p.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			},
		)
	}
	return &pb.GetProductsResponse{Products: products}, nil
}
