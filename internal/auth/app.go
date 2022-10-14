package auth

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	"github.com/sergeysynergy/gok/internal/auth/data/model"
	sesRepo "github.com/sergeysynergy/gok/internal/auth/data/repository/psql/session"
	usrRepo "github.com/sergeysynergy/gok/internal/auth/data/repository/psql/user"
	ServerGRPC "github.com/sergeysynergy/gok/internal/auth/delivery/server"
	sesUC "github.com/sergeysynergy/gok/internal/auth/useCase/session"
	usrUC "github.com/sergeysynergy/gok/internal/auth/useCase/user"
	"github.com/sergeysynergy/gok/internal/consts"
	pb "github.com/sergeysynergy/gok/proto"
	"github.com/sergeysynergy/gok/tool/conf/service"
)

type App struct {
	cfg *service.Conf
	lg  *zap.Logger

	dbOnce     *sync.Once
	db         *gorm.DB
	grpcServer *grpc.Server

	user usrUC.UseCase
}

func New(cfg *service.Conf, lg *zap.Logger) *App {
	s := &App{
		dbOnce: &sync.Once{},
		cfg:    cfg,
		lg:     lg,
	}

	//userRepo := userRepo.New(db)
	//userUC := user.New(lg, userRepo)

	s.init()

	return s
}

func (a *App) init() {
	a.dbConnect()
	a.initUseCase()
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

		err = db.AutoMigrate(&model.User{}, &model.Session{})
		if err != nil {
			a.lg.Fatal(fmt.Sprintf("Auto migration has failed: %s", err))
		}

		a.db = db
		a.lg.Info("Established connection with DB")
	})
}

func (a *App) initUseCase() {
	sessionRepo := sesRepo.New(a.db)
	session := sesUC.New(a.lg, sessionRepo)
	userRepo := usrRepo.New(a.db)
	a.user = usrUC.New(a.lg, userRepo, session)
}

func (a *App) initGRPCServer() {
	// Create gRPC service server with interceptors.
	a.grpcServer = grpc.NewServer(
	//grpc.UnaryInterceptor(ServerGRPC.UnaryEncrypt),
	)

	// Register our service with realization for protobuf methods.
	srv := ServerGRPC.New(a.lg, a.user)
	pb.RegisterAuthServer(a.grpcServer, srv)
}

// runGraceDown Gracefully shutdown service on signals syscall.SIGTERM, syscall.SIGINT and syscall.SIGQUIT.
func (a *App) runGraceDown() {
	// Properly finish work with `zap` logger.
	defer a.lg.Sync()

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
		listen, err := net.Listen("tcp", a.cfg.Addr)
		if err != nil {
			a.lg.Fatal(err.Error())
		}

		a.lg.Info(fmt.Sprintf("gRPC service server started at: %s", a.cfg.Addr))
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
