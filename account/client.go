package account

import (
	"context"
	pb "github.com/Mostbesep/microservice-com-temp/account/pb/microservice-com-temp.account.pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.AccountServiceClient
}

func NewClient(url string) (*Client, error) {

	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := pb.NewAccountServiceClient(conn)
	return &Client{conn: conn, service: c}, nil
}

func (c *Client) Close() error {
	err := c.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) PostAccount(ctx context.Context, name string) (Account, error) {
	response, err := c.service.PostAccount(ctx, &pb.PostAccountRequest{Name: name})
	if err != nil {
		return Account{}, err
	}
	return Account{
		ID:   response.Account.Id,
		Name: response.Account.Name,
	}, nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (Account, error) {
	response, err := c.service.GetAccount(ctx, &pb.GetAccountRequest{Id: id})
	if err != nil {
		return Account{}, err
	}
	return Account{
		ID:   response.Account.Id,
		Name: response.Account.Name,
	}, nil
}

func (c *Client) GetAccounts(ctx context.Context, skip uint64, take uint64) (*[]Account, error) {
	response, err := c.service.GetAccounts(ctx, &pb.GetAccountsRequest{Skip: skip, Take: take})
	if err != nil {
		return nil, err
	}
	accounts := make([]Account, len(response.Accounts))
	for i, account := range response.Accounts {
		accounts[i] = Account{
			ID:   account.Id,
			Name: account.Name,
		}
	}
	return &accounts, nil
}
