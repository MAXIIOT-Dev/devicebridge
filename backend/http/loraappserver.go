/*
 * @Description: 通过HTTP接收lora-app-server的数据
 * @Copyright: Maxiiot(c) 2019
 * @Author: tgq
 * @LastEditors: tgq
 * @Date: 2019-04-11 16:58:01
 * @LastEditTime: 2019-04-19 10:32:53
 */

package http

import (
	"context"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maxiiot/vbaseBridge/backend"

	log "github.com/sirupsen/logrus"
)

// HttpBackend http backend server
type HttpBackend struct {
	rxPacketChan chan backend.DataUpPayloadChan
	httpServer   *http.Server
}

// New return new HttpBackend
func New(addr string) *HttpBackend {
	dataChan := make(chan backend.DataUpPayloadChan, 10)
	r := router(dataChan)
	serv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.WithField("port", serv.Addr).Info("lora  http web server start.")
		if err := serv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.WithError(err).Fatal("lora http web server error.")
			}
		}
	}()

	return &HttpBackend{
		rxPacketChan: dataChan,
		httpServer:   serv,
	}
}

// Close close resource
func (b *HttpBackend) Close() error {
	close(b.rxPacketChan)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return b.httpServer.Shutdown(ctx)
}

// Notice implement notice.
func (b *HttpBackend) Notice(notice map[string]bool) {}

func loraAppServer(dataChan chan backend.DataUpPayloadChan) gin.HandlerFunc {
	return func(c *gin.Context) {
		var rxdata backend.DataUpPayload
		err := c.ShouldBind(&rxdata)
		if err != nil {
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}

		if b, err := hex.DecodeString(rxdata.Data); err == nil {
			data := backend.DataUpPayloadChan{
				Data:   b,
				DevEUI: rxdata.DevEUI,
			}
			dataChan <- data
		} else {
			log.WithError(err).Error("hex deocde payload data error ")
			c.Writer.WriteHeader(http.StatusBadRequest)
			return
		}

		c.Writer.WriteHeader(200)
	}
}

func router(dataChan chan backend.DataUpPayloadChan) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/lora/app", loraAppServer(dataChan))
	return r
}

// HandleUplinks 处理lora上行数据
func (b *HttpBackend) HandleUplinks(wg *sync.WaitGroup) {
	for uplink := range b.rxPacketChan {
		go func(uplink backend.DataUpPayloadChan) {
			wg.Add(1)
			defer wg.Done()
			if err := backend.HandleUplink(uplink); err != nil {
				log.WithFields(log.Fields{
					"device": uplink.DevEUI,
					"data":   hex.EncodeToString(uplink.Data),
				}).Errorf("process device uplink data error: %s", err)
			}
		}(uplink)
	}

}
