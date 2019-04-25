package backend

import (
	"errors"
	"fmt"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/maxiiot/devicebridge/backend/protocol"
	"github.com/maxiiot/devicebridge/storage"
	log "github.com/sirupsen/logrus"
)

var topic = "device/%s/%s"

// HandleUplink handle uplink data
func HandleUplink(conn paho.Client, data DataUpPayloadChan) error {
	dev, err := storage.GetDeviceByEUI(data.DevEUI.String())
	if err != nil {
		return err
	}

	switch dev.ProtocolType {
	case storage.ProtocolHumiture:
		var hums protocol.Humitures
		err := hums.Unmarshal(data.Data)
		if err != nil {
			return err
		}
		for _, hum := range hums.Hums {

			topictemp := fmt.Sprintf(topic, data.DevEUI, "temp")
			if token := conn.Publish(topictemp, 0, false, fmt.Sprintf("%.1f", hum.Temperature)); token.Wait() && token.Error() != nil {
				log.WithError(token.Error()).Errorf("publish %s %.1f", topictemp, hum.Temperature)
			} else {
				log.Infof("publish success,topic:%s msg:%.1f", topictemp, hum.Temperature)
			}

			topichum := fmt.Sprintf(topic, data.DevEUI, "hum")
			if token := conn.Publish(topichum, 0, false, fmt.Sprintf("%.1f", hum.Humidity)); token.Wait() && token.Error() != nil {
				log.WithError(token.Error()).Errorf("publish %s %.1f", topichum, hum.Humidity)
			} else {
				log.Infof("publish success,topic:%s msg:%.1f", topichum, hum.Humidity)
			}

			topicele := fmt.Sprintf(topic, data.DevEUI, "ele")
			if token := conn.Publish(topicele, 0, false, fmt.Sprintf("%.1f", hum.Electricity)); token.Wait() && token.Error() != nil {
				log.WithError(token.Error()).Errorf("publish %s %.1f", topicele, hum.Electricity)
			} else {
				log.Infof("publish success,topic: %s msg: %.1f", topicele, hum.Electricity)
			}

		}
	case storage.ProtocolSmoke:
		if len(data.Data) < 5 {
			return errors.New("data format error.")
		}
		var smoke protocol.Smoke
		err := smoke.Unmarshal(data.Data[4:])
		if err != nil {
			return err
		}
		if smoke.IsHeartBeat {
			topicsmoke := fmt.Sprintf(topic, data.DevEUI, "smoke")
			if token := conn.Publish(topicsmoke, 0, false, "heartbeat"); token.Wait() && token.Error() != nil {
				log.WithError(token.Error()).Errorf("publish %s %s", topicsmoke, "heartbeat")
			} else {
				log.Infof("publish success,topic: %s msg: %s", topicsmoke, "heartbeat")
			}
		}
		if smoke.Alarm != nil {
			topicsmoke := fmt.Sprintf(topic, data.DevEUI, "smoke")
			if token := conn.Publish(topicsmoke, 0, false, smoke.Alarm.String()); token.Wait() && token.Error() != nil {
				log.WithError(token.Error()).Errorf("publish %s %s", topicsmoke, smoke.Alarm.String())
			} else {
				log.Infof("publish success,topic: %s msg: %s", topicsmoke, smoke.Alarm.String())
			}
		}

	}

	return nil
}
