package auth

import (
	"context"
	"fmt"
	"github.com/sergeysynergy/gok/internal/auth/useCase"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	gormMysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	ServerGRPC "github.com/sergeysynergy/gok/internal/auth/delivery/server"
	"github.com/sergeysynergy/gok/internal/consts"
	pb "github.com/sergeysynergy/gok/proto"
	"github.com/sergeysynergy/gok/tool/conf/service"
)

type Storage struct {
	cfg *service.Conf
	lg  *zap.Logger

	dbOnce     sync.Once
	db         *gorm.DB
	grpcServer *grpc.Server

	user useCase.UseCase
}

func New(cfg *service.Conf, lg *zap.Logger) *Storage {
	s := &Storage{
		cfg: cfg,
		lg:  lg,
	}

	//userRepo := userRepo.New(db)
	//userUC := user.New(lg, userRepo)

	s.init()

	return s
}

func (s *Storage) init() {
	s.newDatabase()
	s.newGRPCServer()
	s.newUseCases()
}

func (s *Storage) newDatabase() {
	s.dbOnce.Do(func() {
		db, err := gorm.Open(gormMysqlDriver.Open(s.cfg.DSN), &gorm.Config{
			PrepareStmt:    true,
			NamingStrategy: schema.NamingStrategy{SingularTable: true},
		})
		if err != nil {
			s.lg.Fatal(fmt.Sprintf("Failed connect to DB: %s", err))
		}

		s.db = db
		s.lg.Info("Established connection with DB")
	})
}

func (s *Storage) newGRPCServer() {
	// Create gRPC service with interceptors.
	s.grpcServer = grpc.NewServer(
	//grpc.UnaryInterceptor(ServerGRPC.UnaryEncrypt),
	)

	// Register our service with realization for protobuf methods.
	service := ServerGRPC.New(s.lg, s.user)
	pb.RegisterUsersServer(s.grpcServer, service)
}

func (s *Storage) newUseCases() {
	//repo := rep
	//s.user = user.New(s.lg, s.repo)
}

// runGraceDown Gracefully shutdown service on signals syscall.SIGTERM, syscall.SIGINT and syscall.SIGQUIT.
func (s *Storage) runGraceDown() {
	// Properly finish work with `zap` logger.
	defer s.lg.Sync()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-sig

	shutdownCtx, shutdownCtxCancel := context.WithTimeout(context.Background(), consts.ServerGraceTimeout)
	defer shutdownCtxCancel()
	// Force shutdown after grace timeout exceeded.
	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			s.lg.Fatal("Graceful shutdown timed out! Forcing exit.")
		}
	}()

	// Gracefully shutdown gRPC service.
	s.grpcServer.GracefulStop()
	s.lg.Info("Gracefully shutdown gRPC-service")
}

func (s *Storage) start() {
	go func() {
		listen, err := net.Listen("tcp", s.cfg.Addr)
		if err != nil {
			s.lg.Fatal(err.Error())
		}

		s.lg.Info(fmt.Sprintf("gRPC service server started at: %s", s.cfg.Addr))
		if err = s.grpcServer.Serve(listen); err != nil {
			s.lg.Fatal(err.Error())
		}
	}()
}

// Run start gRPC service server and then run graceful shutdown listener:
// gracefully shutdown on signals syscall.SIGTERM, syscall.SIGINT and syscall.SIGQUIT.
func (s *Storage) Run() {
	s.start()
	s.runGraceDown()
}
