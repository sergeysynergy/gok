// Package location Пакет реализует базовый тип Местоположение и методы работы с ним.
package server

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcZap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcCtxTags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	gormMysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	client "back-mfsb/api/pb_go"
	generated "back-mfsb/pb/location"
	locRepo "back-mfsb/service/location/internal/data/repository/memory"
	myGRPC "back-mfsb/service/location/internal/delivery/server/grpc"
	"back-mfsb/service/location/internal/domain/useCase"
	"back-mfsb/tool/simpleLocationConfig"
)

type Server struct {
	cfg *simpleLocationConfig.Config

	lg         *zap.Logger
	dbOnce     sync.Once
	db         *gorm.DB
	clientConn *grpc.ClientConn
	cc         client.AuthorizationClient
	httpServer *http.Server

	gRPCLogger  *zap.Logger
	gRPCServer  *grpc.Server
	gRPCService *myGRPC.LocationService

	repo useCase.Repo
	uc   useCase.UseCase
}

func NewServer(
	cfg *simpleLocationConfig.Config,
	gRPCLogger *zap.Logger,
	ucLogger *zap.Logger,
) *Server {
	s := &Server{
		cfg:        cfg,
		gRPCLogger: gRPCLogger,
		lg:         ucLogger,
	}

	s.init()

	s.gRPCService = myGRPC.New(
		gRPCLogger,
		s.cc,
		s.repo,
		s.uc,
	)

	return s
}

func (s *Server) init() {
	s.newDatabase()
	s.newGRPCServer()
	s.newAuthorizationClient()
	s.newHTTPServer()
	s.newRepos()
	s.newUseCases()
}

func (s *Server) newDatabase() {
	//db, err := gorm.Open(postgres.New(postgres.Config{
	//	DSN:                  s.dsnPG,
	//	PreferSimpleProtocol: true, // disables implicit prepared statement usage
	//}), &gorm.Config{})
	//if err != nil {
	//  s.lg.Fatal("Connection to Postgres is failed: %w", err)
	//}

	s.dbOnce.Do(func() {
		db, err := gorm.Open(gormMysqlDriver.Open(s.cfg.Database.DSN), &gorm.Config{
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

func (s *Server) newGRPCServer() {
	s.gRPCServer = grpc.NewServer(
		grpcMiddleware.WithUnaryServerChain(
			grpcCtxTags.UnaryServerInterceptor(grpcCtxTags.WithFieldExtractor(grpcCtxTags.CodeGenRequestFieldExtractor)),
			grpcZap.UnaryServerInterceptor(s.lg),
			// grpc_recovery.UnaryServerInterceptor(),
			//watcher.UnaryServerInterceptor(lg),
		),
		grpcMiddleware.WithStreamServerChain(
			grpcCtxTags.StreamServerInterceptor(grpcCtxTags.WithFieldExtractor(grpcCtxTags.CodeGenRequestFieldExtractor)),
			grpcZap.StreamServerInterceptor(s.lg),
			// grpc_recovery.StreamServerInterceptor(),
			//watcher.StreamInterceptor(),
		),
	)
}

func (s *Server) newAuthorizationClient() {
	// Connection to client-service
	var err error
	s.clientConn, err = grpc.Dial(
		s.cfg.Service.IP+":"+strconv.Itoa(s.cfg.Service.GRPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		/*grpc.WithPerRPCCredentials(&server.LoginCreds{
			Username: "",
			Password: "",
		}),*/
	)
	if err != nil {
		s.lg.Fatal(fmt.Sprintf("Не удалось создать клиент авторизации: %s", err))
	}

	s.cc = client.NewAuthorizationClient(s.clientConn)
}

func (s *Server) newHTTPServer() {
	grpcWebServer := grpcweb.WrapServer(s.gRPCServer)

	s.httpServer = &http.Server{
		Addr: ":" + strconv.Itoa(s.cfg.Service.GRPCPort),
		Handler: h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.lg.Info(fmt.Sprintf("Версия HTTP: %v", r.Proto))
			if r.ProtoMajor == 2 {
				grpcWebServer.ServeHTTP(w, r)
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-Agent, X-Grpc-Web")
				w.Header().Set("grpc-status", "")
				w.Header().Set("grpc-message", "")
				if grpcWebServer.IsGrpcWebRequest(r) {
					grpcWebServer.ServeHTTP(w, r)
				}
			}
		}), &http2.Server{}),
	}
}

func (s *Server) newRepos() {
	s.repo = locRepo.New()
}

func (s *Server) newUseCases() {
	s.uc = useCase.New(s.lg, s.repo)
}

func (s *Server) Start() {
	// Объявим необходимые закрытия по завершению работы сервиса.
	defer s.clientConn.Close()
	defer s.lg.Sync()

	generated.RegisterLocationsServer(s.gRPCServer, s.gRPCService)
	reflection.Register(s.gRPCServer)

	s.lg.Info(fmt.Sprint("Location service started at :", strconv.Itoa(s.cfg.Service.GRPCPort)))
	if err := s.httpServer.ListenAndServe(); err != nil {
		s.lg.Fatal("Failed to start Location service")
	}
}
