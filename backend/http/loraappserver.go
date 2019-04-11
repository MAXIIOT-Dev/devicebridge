/*
 * @Description: 通过HTTP接收lora-app-server的数据
 * @Copyright: Maxiiot(c) 2019
 * @Author: tgq
 * @LastEditors: tgq
 * @Date: 2019-04-11 16:58:01
 * @LastEditTime: 2019-04-11 17:33:47
 */

package http

import (
	"encoding/hex"
	"net/http"

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
		log.Fatal(serv.ListenAndServe())
	}()
	return &HttpBackend{
		rxPacketChan: dataChan,
		httpServer:   serv,
	}
}

// RXPacketChan return rxpacketchan
func (b *HttpBackend) RXPacketChan() chan backend.DataUpPayloadChan {
	return b.rxPacketChan
}

// Close close resource
func (b *HttpBackend) Close() error {
	close(b.rxPacketChan)
	return b.httpServer.Close()
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
