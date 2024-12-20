package account

import (
	"context"
	"github.com/segmentio/ksuid"
)

type Service interface {
	PostAccount(ctx context.Context, name string) (Account, error)
	GetAccount(ctx context.Context, id string) (Account, error)
	GetAccounts(ctx context.Context, skip uint64, take uint64) (*[]Account, error)
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type service struct {
	repository Repository
}

func (s *service) PostAccount(ctx context.Context, name string) (Account, error) {
	a := Account{
		Name: name,
		ID:   ksuid.New().String(),
	}
	err := s.repository.PutAccount(ctx, a)
	if err != nil {
		return Account{}, err
	}
	return a, nil
}

func (s *service) GetAccount(ctx context.Context, id string) (Account, error) {
	return s.repository.GetAccountByID(ctx, id)
}

func (s *service) GetAccounts(ctx context.Context, skip uint64, take uint64) (*[]Account, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.ListAccounts(ctx, skip, take)
}

func NewAccountService(repository Repository) Service {
	return &service{
		repository: repository,
	}
}
