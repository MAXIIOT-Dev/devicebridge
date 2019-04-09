package cmd

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/maxiiot/vbaseBridge/backend/lorahandler"
	"github.com/maxiiot/vbaseBridge/backend/mqtt"
	"github.com/maxiiot/vbaseBridge/config"
	"github.com/maxiiot/vbaseBridge/controllers"
	"github.com/maxiiot/vbaseBridge/routers"
	"github.com/maxiiot/vbaseBridge/storage"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var levels = map[string]log.Level{
	log.FatalLevel.String(): log.FatalLevel,
	log.ErrorLevel.String(): log.ErrorLevel,
	log.WarnLevel.String():  log.WarnLevel,
	log.InfoLevel.String():  log.InfoLevel,
	log.DebugLevel.String(): log.DebugLevel,
}

var run = func(cmd *cobra.Command, args []string) error {
	errs := make(chan error, 1)
	setVersion()
	setLogLevel()

	var port string
	if config.Cfg.General.Port > 0 {
		port = fmt.Sprintf(":%d", config.Cfg.General.Port)
	} else {
		port = ":8080"
	}
	r := routers.Route(gin.ReleaseMode)
	serv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	backend, err := mqtt.NewBackend(config.Cfg.Mqtt)
	if err != nil {
		return err
	}
	backendServ := lorahandler.NewServer(backend)

	tasks := []func() error{
		connectPostgres,
		startWebServer(serv, errs),
		startBackendServer(backendServ),
	}
	for _, t := range tasks {
		err := t()
		if err != nil {
			return err
		}
	}
	log.Error(<-errs)
	return nil
}

func setVersion() {
	controllers.SetVersion(version)
}

func setLogLevel() {
	level := unmarshalLogLevel(config.Cfg.General.LogLevel)
	log.SetLevel(level)
	if runtime.GOOS == "windows" {
		log.SetFormatter(&log.TextFormatter{
			DisableColors: true,
		})
	}
}

func unmarshalLogLevel(level string) log.Level {
	if v, ok := levels[level]; ok {
		return v
	}
	return log.InfoLevel
}

func connectPostgres() error {
	storage.Connect(config.Cfg.Postgres.DSN)
	if config.Cfg.Postgres.AutoMigrate {
		if err := storage.Migrate(); err != nil {
			return err
		}
	}
	return nil
}

func startBackendServer(serv *lorahandler.Server) func() error {
	return func() error {
		serv.Start()
		return nil
	}
}

func startWebServer(serv *http.Server, errs chan error) func() error {
	return func() error {
		go func(errs chan error) {
			if err := serv.ListenAndServe(); err != nil {
				errs <- err
			}
		}(errs)
		return nil
	}
}
