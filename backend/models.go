package backend

import (
	"time"

	"github.com/maxiiot/devicebridge/storage"
)

// Location details.
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
}

// RXInfo contains the RX information.
type RXInfo struct {
	GatewayID storage.EUI64 `json:"gatewayID"`
	Name      string        `json:"name"`
	Time      *time.Time    `json:"time,omitempty"`
	RSSI      int           `json:"rssi"`
	LoRaSNR   float64       `json:"loRaSNR"`
	Location  *Location     `json:"location"`
}

// TXInfo contains the TX information.
type TXInfo struct {
	Frequency int `json:"frequency"`
	DR        int `json:"dr"`
}

// maxiiot api
// new DataUpPayload represents a data-up payload.
type DataUpPayload struct {
	ApplicationID   int64      `json:"applicationID,string"`
	ApplicationName string     `json:"applicationName"`
	Time            *time.Time `json:"time"`
	// DeviceAddr      lorawan.DevAddr `json:"devaddr"`
	DevEUI     storage.EUI64 `json:"deveui" binding:"required"`
	DeviceName string        `json:"devname"`
	GatewayEUI storage.EUI64 `json:"gatewayeui"`
	RSSI       int32         `json:"rssi"`
	LoRaSNR    float64       `json:"lsnr"`
	Size       int           `json:"size"`
	Data       string        `json:"data" binding:"required"`
	Base64Data string        `json:"b64_data"`
	Frequency  float64       `json:"freq"`
	DataRate   string        `json:"datr"`
	ADR        bool          `json:"adr"`
	RXInfo     []RXInfo      `json:"rxInfo"`
	TXInfo     TXInfo        `json:"txInfo"`
	FPort      uint8         `json:"port"`
	FCnt       uint32        `json:"uplink_count"`
	//GatewayList     []string        `json:"gateway_list,omitempty"`
	Object interface{} `json:"object,omitempty"`
}

// DataUpPayloadChan DataUpPayloadChan
type DataUpPayloadChan struct {
	Data   []byte
	DevEUI storage.EUI64
}

// type ACKNotification struct {
// 	ApplicationID   int64         `json:"applicationID,string"`
// 	ApplicationName string        `json:"applicationName"`
// 	DeviceName      string        `json:"deviceName"`
// 	DevEUI          storage.EUI64 `json:"devEUI"`
// 	Acknowledged    bool          `json:"acknowledged"`
// 	FCnt            uint32        `json:"fCnt"`
// }

// type ACKNotificationChan struct {
// 	DevEUI string
// 	FCnt   uint32
// }
