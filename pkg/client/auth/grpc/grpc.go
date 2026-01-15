package grpc_auth

import (
	"context"
	"errors"

	grpcpetv1 "github.com/Rustamchick/protobuff/gen/go/pet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	api grpcpetv1.AuthClient
}

func NewClient(ctx context.Context, addr string) (*Client, error) {
	cc, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, errors.New("Error connecting to " + addr)
	}

	return &Client{grpcpetv1.NewAuthClient(cc)}, nil
}

func (c *Client) Register(ctx context.Context, username, password string) (int64, error) {
	RegReq := &grpcpetv1.RegisterRequest{
		Email:    username,
		Password: password,
	}

	resp, err := c.api.Register(ctx, RegReq)
	if err != nil {
		return 0, err
	}

	return resp.UserId, nil
}

func (c *Client) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	isAdminReq := &grpcpetv1.IsAdminRequest{
		UserId: userId,
	}

	resp, err := c.api.IsAdmin(ctx, isAdminReq)
	if err != nil {
		return false, err
	}

	return resp.IsAdmin, nil
}
