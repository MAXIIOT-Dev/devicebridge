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
			log.Infof("publish %s %.1f", topictemp, hum.Temperature)
			if token := conn.Publish(topictemp, 0, false, fmt.Sprintf("%.1f", hum.Temperature)); token.Wait() && token.Error() != nil {
				log.Error(token.Error())
			}
			topichum := fmt.Sprintf(topic, data.DevEUI, "hum")
			log.Infof("publish %s %.1f", topichum, hum.Humidity)
			conn.Publish(topichum, 0, false, fmt.Sprintf("%.1f", hum.Humidity))
			topicele := fmt.Sprintf(topic, data.DevEUI, "ele")
			log.Infof("publish %s %.1f", topicele, hum.Electricity)
			conn.Publish(topicele, 0, false, fmt.Sprintf("%.1f", hum.Electricity))

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
			log.Infof("publish %s %s", topicsmoke, "heartbeat")
			conn.Publish(topicsmoke, 0, false, "heartbeat")
		}
		if smoke.Alarm != nil {
			topicsmoke := fmt.Sprintf(topic, data.DevEUI, "smoke")
			log.Infof("publish %s %s", topicsmoke, smoke.Alarm.String())
			conn.Publish(topicsmoke, 0, false, smoke.Alarm.String())
		}

	}

	return nil
}
