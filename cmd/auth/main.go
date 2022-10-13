package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"

	"github.com/sergeysynergy/gok/internal/auth"
	"github.com/sergeysynergy/gok/pkg/utils"
	"github.com/sergeysynergy/gok/tool/conf/service"
	"github.com/sergeysynergy/gok/tool/zapLogger"
)

var (
	buildVersion string
	buildDate    string
)

func main() {
	var err error

	cfg := service.New(
		service.WithDSN("user=gok password=Passw0rd33 host=localhost port=45432 dbname=auth"),
	)

	flag.BoolVar(&cfg.Debug, "debug", cfg.Debug, "run service in debug mode")
	flag.StringVar(&cfg.Addr, "a", cfg.Addr, "address to listen on")
	flag.StringVar(&cfg.DSN, "d", cfg.DSN, "Postgres DSN")
	flag.Parse()

	// Rewriting config with environment variables: highest priority.
	err = env.Parse(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	lg := zapLogger.NewServerLogger(cfg.Debug)

	// Warnings for using DEBUG mode.
	if cfg.Debug {
		lg.Warn("ATTENTION: service is running in debug mode")
		// Выведем переменные окружения.
		lg.Debug(fmt.Sprintf("%#v", cfg))
	}

	// Print build variables that has been set on linking stage, for example:
	// go run -ldflags "-X main.Version=v1.0.1" main.go
	lg.Info(fmt.Sprintf("GoK service version: %s", utils.CheckNA(buildVersion)))
	lg.Info(fmt.Sprintf("GoK build date: %s", utils.CheckNA(buildDate)))

	server := auth.New(cfg, lg)
	server.Run()

	// Проверка на выполнение контракта интерфейса.
	//var _ service.Repo = new(pgsql.Storage)

	// Получим реализацию репозитория для работы с БД.
	//repo := service.Repo(memory.New())
	//repoDB := pgsql.New(cfg.DatabaseDSN)
	//if repoDB != nil {
	//	repo = service.Repo(repoDB)
	//}

	// Проверка на выполнение контракта интерфейса.
	//var _ service.FileRepo = new(filestore.FileStore)
	//// Создадим файловое хранилище на базе Storage
	//fileStorer := filestore.New(
	//	filestore.WithStorer(repo),
	//	filestore.WithStoreFile(cfg.StoreFile),
	//	filestore.WithRestore(cfg.Restore),
	//	filestore.WithStoreInterval(cfg.StoreInterval),
	//)
	//
	//uc := service.New(
	//	service.WithDBStorer(repo),
	//	service.WithFileStorer(fileStorer),
	//)
	//
	//// Подключим обработчики запросов.
	//privateKey, err := crypter.OpenPrivate(cfg.CryptoKey)
	//if err != nil {
	//	log.Println("[WARNING] Failed to get private key - ", err)
	//}
	//h := handlers.New(uc,
	//	//handlers.WithFileStorer(fileStorer),
	//	//handlers.WithDBStorer(dbStorer),
	//	handlers.WithKey(cfg.Key),
	//	handlers.WithPrivateKey(privateKey),
	//	handlers.WithTrustedSubnet(cfg.TrustedSubnet),
	//)
	//
	//// Проинициализируем сервер с использованием ранее объявленных обработчиков и файлового хранилища.
	//s := httpserver.New(uc, h.GetRouter(),
	//	httpserver.WithAddress(cfg.Addr),
	//	//httpserver.WithFileStorer(fileStorer),
	//	//httpserver.WithDBStorer(dbStorer),
	//)
	//
	////go http.ListenAndServe(":8090", nil) // запускаем сервер для нужд профилирования
	//
	//s.Serve() // запускаем основной http-сервер
}
