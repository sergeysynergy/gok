package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"

	"github.com/sergeysynergy/gok/internal/entity"
	pb "github.com/sergeysynergy/gok/proto"
)

type Client struct {
	cl pb.AuthClient
}

func (c Client) dial() {
	conn, err := grpc.Dial(":7000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c.cl = pb.NewAuthClient(conn)
}

func (c Client) GetUser(ctx context.Context, token string) (*entity.User, error) {
	c.dial()

	req := &pb.GetUserRequest{
		Token: token,
	}

	fmt.Println("HERE HERE KITTY")
	resp, err := c.cl.GetUser(ctx, req)
	fmt.Println("HERE HERE KITTY")
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
