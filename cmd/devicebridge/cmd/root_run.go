package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/maxiiot/devicebridge/backend/server"
	"github.com/maxiiot/devicebridge/config"
	"github.com/maxiiot/devicebridge/controllers"
	"github.com/maxiiot/devicebridge/routers"
	"github.com/maxiiot/devicebridge/storage"

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

var run = func(cmd *cobra.Command, args []string) (err error) {
	setVersion()
	setLogLevel()
	if err := connectPostgres(); err != nil {
		return err
	}

	server.Serv, err = server.NewServer(config.Cfg)
	if err != nil {
		return err
	}

	serv := newHttpServer()

	tasks := []func() error{
		startWebServer(serv),
		startBackendServer(server.Serv),
	}
	for _, t := range tasks {
		err := t()
		if err != nil {
			return err
		}
	}

	quit := make(chan os.Signal)
	exitChan := make(chan struct{})
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	log.WithField("signal", <-quit).Info("quit signal received: ")

	go func() {
		log.Println("Shutdown web Server ...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := serv.Shutdown(ctx); err != nil {
			log.Error("Shutdown web Server error: ", err)
		}
		// catching ctx.Done(). timeout of 5 seconds.
		// select {
		// case <-ctx.Done():
		// 	log.Println("timeout of 5 seconds.")
		// }

		log.Println("Shutdown backend server...")
		if err := server.Serv.Stop(); err != nil {
			log.Error("backend server shutdown error:", err)
		}
		log.Println("gracefull shutdown complete")
		exitChan <- struct{}{}
	}()

	select {
	case <-exitChan:
	case s := <-quit:
		log.WithField("signal", s).Info("signal received, stopping immediately")
	}

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

func startBackendServer(serv *server.Server) func() error {
	return func() error {
		serv.Start()
		return nil
	}
}

func startWebServer(serv *http.Server) func() error {
	return func() error {
		go func() {
			log.WithField("port", serv.Addr).Info("start web server.")
			if err := serv.ListenAndServe(); err != nil {
				if err != http.ErrServerClosed {
					log.WithError(err).Fatal("web server error.")
				}
			}
		}()
		return nil
	}
}
