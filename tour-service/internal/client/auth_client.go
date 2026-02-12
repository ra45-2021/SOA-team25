package client

import (
	"context"
	"tour-service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	client pb.AuthServiceClient
}

func NewAuthClient(addr string) (*AuthClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &AuthClient{client: pb.NewAuthServiceClient(conn)}, nil
}

func (c *AuthClient) IsGuide(ctx context.Context, userID int64) (bool, error) {
	resp, err := c.client.CheckAuthor(ctx, &pb.CheckAuthorRequest{Id: userID})
	if err != nil {
		return false, err
	}
	return resp.Exists && resp.Role == "GUIDE", nil
}