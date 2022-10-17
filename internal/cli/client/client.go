// Package client implements gRPC client for working with GoK storage API.
package client

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"

	pb "github.com/sergeysynergy/gok/proto"
)

type Client struct {
	lg   *zap.Logger
	once *sync.Once
	// gRPC Auth service address.
	authAddr string
	// gRPC Storage service address.
	storageAddr string
}

func New(logger *zap.Logger, authAddr, storageAddr string) *Client {
	c := &Client{
		lg:          logger,
		authAddr:    authAddr,
		storageAddr: storageAddr,
	}

	return c
}

func (c *Client) getAuthConnect() (pb.AuthClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(c.authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.lg.Fatal(err.Error())
	}

	return pb.NewAuthClient(conn), conn
}

func (c *Client) getStorageConnect() pb.StorageClient {
	conn, err := grpc.Dial(c.storageAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.lg.Fatal(err.Error())
	}

	return pb.NewStorageClient(conn)
}

func (c *Client) signIn(ctx context.Context) error {
	auth, conn := c.getAuthConnect()
	defer conn.Close()

	resp, err := auth.SignIn(ctx, &pb.SignInRequest{
		User: &pb.UserForAdd{
			Login: c.user,
		}

	})
	if err != nil {
		return err
	}

	return nil
}

// sendReport Отправляет значения всех метрик на сервер.
//func send(ctx context.Context, publicKey *rsa.PublicKey, c pb.MetricsClient, hm []metrics.Metrics) {
//	// Преобразуем метрики для отправки на сервер.
//	prm := metrics.NewProxyMetrics()
//	gauges := make([]*pb.Gauge, 0, len(prm.Gauges))
//	counters := make([]*pb.Counter, 0, len(prm.Counters))
//	for _, v := range hm {
//		switch v.MType {
//		case "gauge":
//			gauges = append(gauges, &pb.Gauge{
//				Id:    v.ID,
//				Value: *v.Value,
//			})
//		case "counter":
//			counters = append(counters, &pb.Counter{
//				Id:    v.ID,
//				Delta: *v.Delta,
//			})
//		}
//	}
//
//	// Зашифруем метрики, если определён ключ для шифрования
//	md := metadata.MD{}
//	if publicKey != nil {
//		md = metadata.New(map[string]string{"token": "crypted"})
//	}
//	ctx = metadata.NewOutgoingContext(ctx, md)
//
//	// отправим метрики на сервер
//	_, err := c.AddMetrics(ctx, &pb.AddMetricsRequest{
//		Gauges:   gauges,
//		Counters: counters,
//	})
//	if err != nil {
//		log.Println("[ERROR] Неудача отправки метрик -", err)
//	}
//
//	log.Println("[DEBUG] Метрики успешно отправлены на сервер по gRPC")
//}
