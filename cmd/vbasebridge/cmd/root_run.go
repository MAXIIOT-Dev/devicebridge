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
	setVersion()
	setLogLevel()
	if err := connectPostgres(); err != nil {
		return err
	}

	devs, err := storage.GetDevicesEUI()
	if err != nil {
		return err
	}

	mqtt.MQTTBackend = mqtt.NewBackend(config.Cfg.Mqtt, devs)
	if err != nil {
		return err
	}
	backendServ, err := lorahandler.NewServer(mqtt.MQTTBackend)
	if err != nil {
		return err
	}

	errs := make(chan error, 1)
	serv := newHttpServer()
	
	tasks := []func() error{
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

func newHttpServer() *http.Server {
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
	return serv
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
