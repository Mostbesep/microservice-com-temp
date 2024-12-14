package main

import "context"

type accountResolver struct {
	server *Server
}

// Orders
func (a accountResolver) Orders(ctx context.Context, obj *Account) ([]*Order, error) {
	//TODO implement me
	panic("implement me")
}
