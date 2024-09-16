package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"reimirno.com/golinks/pkg/config"
	"reimirno.com/golinks/pkg/logging"
	"reimirno.com/golinks/pkg/mapper"
	"reimirno.com/golinks/pkg/version"
	"reimirno.com/golinks/svr/crud"
	"reimirno.com/golinks/svr/redirector"
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

	configurators := make([]mapper.MapperConfigurer, len(cfg.Mapper.Mappers))
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

	svrErrChan := make(chan error, 2)
	sigTermChan := make(chan os.Signal, 1)

	redirectorServer.Start(svrErrChan)
	crudServer.Start(svrErrChan)

	signal.Notify(sigTermChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigTermChan:
		logger.Infof("Received shutdown signal, shutting down...")
		redirectorServer.Stop()
		crudServer.Stop()
		mapperManager.Teardown()
		os.Exit(0)
	case err = <-svrErrChan:
		if err != nil {
			logger.Errorf("Server error: %v", err)
		}
	}

	logger.Info("Application stopped. Bye!")
}
