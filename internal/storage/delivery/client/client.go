package client

import (
	"context"
	"github.com/sergeysynergy/gok/internal/entity"
	pb "github.com/sergeysynergy/gok/proto"
)

type Client struct {
	auth pb.AuthClient
}

func New(auth pb.AuthClient) *Client {
	return &Client{
		auth: auth,
	}
}

func (c Client) GetUser(ctx context.Context, token string) (*entity.User, error) {
	req := &pb.GetUserRequest{
		Token: token,
	}

	resp, err := c.auth.GetUser(ctx, req)
	if err != nil {
		return nil, err
	}
	usrPB := resp.User

	usr := &entity.User{
		ID:    entity.UserID(usrPB.ID),
		Login: usrPB.Login,
	}

	return usr, nil
}
