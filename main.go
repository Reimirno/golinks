package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/reimirno/golinks/pkg/config"
	"github.com/reimirno/golinks/pkg/logging"
	"github.com/reimirno/golinks/pkg/mapper"
	"github.com/reimirno/golinks/pkg/types"
	"github.com/reimirno/golinks/pkg/version"
	"github.com/reimirno/golinks/svr/crud"
	"github.com/reimirno/golinks/svr/crud_http"
	"github.com/reimirno/golinks/svr/redirector"
	"go.uber.org/zap"
)

var (
	configFile string
	logger     *zap.SugaredLogger

	Version   string
	Commit    string
	BuildDate string
)

func main() {
	flag.StringVar(&configFile, "config", "./files/config.yaml", "Path to the config file")
	flag.Parse()

	cfg, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := logging.Initialize(cfg.Server.Debug); err != nil {
		log.Fatalf("Failed to initialize logging: %v", err)
	}
	logger = logging.NewLogger("main")

	bld := version.BuildVersion{
		Version:   Version,
		Commit:    Commit,
		BuildDate: BuildDate,
	}
	logger.Info("Application starting...")
	logger.Info(bld)

	configurators := make([]types.MapperConfigurer, len(cfg.Mapper.Mappers))
	for i, wrapper := range cfg.Mapper.Mappers {
		configurators[i] = wrapper.MapperConfigurer
	}
	mapperManager, err := mapper.NewMapperManager(cfg.Mapper.Persistor, configurators)
	if err != nil {
		log.Fatalf("Failed to create mapper manager: %v", err)
	}

	redirectorServer, err := redirector.NewServer(mapperManager, cfg.Server.Port.Redirector)
	if err != nil {
		log.Fatalf("Failed to create redirector server: %v", err)
	}

	crudServer, err := crud.NewServer(mapperManager, cfg.Server.Port.Crud, cfg.Server.Debug)
	if err != nil {
		log.Fatalf("Failed to create grpc server: %v", err)
	}

	crudHttpServer, err := crud_http.NewServer(mapperManager, cfg.Server.Port.CrudHttp)
	if err != nil {
		log.Fatalf("Failed to create crud http server: %v", err)
	}

	svrErrChan := make(chan error, 3)
	sigTermChan := make(chan os.Signal, 1)

	redirectorServer.Start(svrErrChan)
	crudServer.Start(svrErrChan)
	crudHttpServer.Start(svrErrChan)

	signal.Notify(sigTermChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigTermChan:
		logger.Infof("Received shutdown signal, shutting down...")
		redirectorServer.Stop()
		crudServer.Stop()
		crudHttpServer.Stop()
		mapperManager.Teardown()
		os.Exit(0)
	case err = <-svrErrChan:
		if err != nil {
			logger.Errorf("Server error: %v", err)
		}
	}

	logger.Info("Application stopped. Bye!")
}
