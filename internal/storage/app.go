package storage

import (
	"context"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"github.com/sergeysynergy/gok/internal/consts"
	gokModel "github.com/sergeysynergy/gok/internal/data/model"
	recRepo "github.com/sergeysynergy/gok/internal/data/repository/sql/record"
	"github.com/sergeysynergy/gok/internal/storage/data/model"
	brnRepo "github.com/sergeysynergy/gok/internal/storage/data/repository/psql/branch"
	storageClient "github.com/sergeysynergy/gok/internal/storage/delivery/client"
	ServerGRPC "github.com/sergeysynergy/gok/internal/storage/delivery/server"
	brnUC "github.com/sergeysynergy/gok/internal/storage/useCase/branch"
	recUC "github.com/sergeysynergy/gok/internal/storage/useCase/record"
	pb "github.com/sergeysynergy/gok/proto"
	"github.com/sergeysynergy/gok/tool/conf/service"
)

type App struct {
	cfg *service.Conf
	lg  *zap.Logger

	dbOnce     *sync.Once
	db         *gorm.DB
	grpcServer *grpc.Server

	authClient     pb.AuthClient
	authClientConn *grpc.ClientConn

	branch brnUC.UseCase
	record recUC.UseCase
}

func New(cfg *service.Conf, lg *zap.Logger) *App {
	s := &App{
		dbOnce: &sync.Once{},
		cfg:    cfg,
		lg:     lg,
	}
	s.init()

	return s
}

func (a *App) init() {
	a.dbConnect()
	a.initAuthClient()
	a.initUseCases()
	a.initGRPCServer()
}

func (a *App) dbConnect() {
	a.dbOnce.Do(func() {
		db, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  a.cfg.DSN,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &gorm.Config{})
		if err != nil {
			a.lg.Fatal(fmt.Sprintf("Connection to Postgres is failed: %s", err))
		}

		// Create and migrate database tables.
		err = db.AutoMigrate(&model.Branch{}, &gokModel.Record{})
		if err != nil {
			a.lg.Fatal(fmt.Sprintf("Auto migration has failed: %s", err))
		}

		a.db = db
		a.lg.Info("Established connection with DB")
	})
}

func (a *App) initAuthClient() {
	conn, err := grpc.Dial(a.cfg.AuthAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		a.lg.Fatal(err.Error())
	}

	a.authClient = pb.NewAuthClient(conn)
	a.authClientConn = conn
}

func (a *App) initUseCases() {
	recordRepo := recRepo.New(a.db)
	a.record = recUC.New(a.lg, recordRepo)

	branchRepo := brnRepo.New(a.db)
	client := storageClient.New(a.authClient)
	a.branch = brnUC.New(a.lg, branchRepo, client, a.record)
}

func (a *App) initGRPCServer() {
	// Create gRPC service server with interceptors.
	a.grpcServer = grpc.NewServer(
		//grpc.UnaryInterceptor(ServerGRPC.UnaryEncrypt),
	)

	// Register our service with realization for protobuf methods.
	srv := ServerGRPC.New(a.lg, a.branch, a.record)
	pb.RegisterStorageServer(a.grpcServer, srv)
}

// runGraceDown Gracefully shutdown service on signals syscall.SIGTERM, syscall.SIGINT and syscall.SIGQUIT.
func (a *App) runGraceDown() {
	// Properly finish work with `zap` logger.
	defer a.lg.Sync()

	// Properly close connection with Auth service.
	defer a.authClientConn.Close()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-sig

	shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), consts.ServerGraceTimeout)
	defer shutdownCtxCancel()
	// Force shutdown after grace timeout exceeded.
	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			a.lg.Fatal("Graceful shutdown timed out! Forcing exit.")
		}
	}()

	// Gracefully shutdown gRPC service server.
	a.grpcServer.GracefulStop()
	a.lg.Info("Gracefully shutdown gRPC-service")
}

func (a *App) start() {
	go func() {
		listen, err := net.Listen("tcp", a.cfg.StorageAddr)
		if err != nil {
			a.lg.Fatal(err.Error())
		}

		a.lg.Info(fmt.Sprintf("gRPC service server started at: %s", a.cfg.StorageAddr))
		if err = a.grpcServer.Serve(listen); err != nil {
			a.lg.Fatal(err.Error())
		}
	}()
}

// Run start gRPC service and then run graceful shutdown listener:
// gracefully shutdown on signals syscall.SIGTERM, syscall.SIGINT and syscall.SIGQUIT.
func (a *App) Run() {
	a.start()
	a.runGraceDown()
}
