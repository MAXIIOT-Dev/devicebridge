package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

const (
	// ProtocolHumiture humiture protocol
	ProtocolHumiture = "humiture"
	// ProtocolSmoke smoke protocol
	ProtocolSmoke = "smoke"
	// ProtocolDefault default protocol
	ProtocolDefault = "digital"
)

// Device define device model
type Device struct {
	DeviceEUI    EUI64     `db:"device_eui" json:"device_eui"`
	ProtocolType string    `db:"protocol_type" json:"protocol_type"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

// CreateDevice create device on database
func CreateDevice(dev Device) error {
	now := time.Now()
	_, err := db.Exec(`
		insert into device (
			device_eui,
			protocol_type,
			created_at,
			updated_at
		)values($1,$2,$3,$3)`,
		dev.DeviceEUI,
		dev.ProtocolType,
		now,
	)
	return err
}

// GetDeviceByEUI get device by device eui
func GetDeviceByEUI(eui string) (Device, error) {
	var devEUI EUI64
	var dev Device
	err := devEUI.UnmarshalText([]byte(eui))
	if err != nil {
		return dev, err
	}

	err = sqlx.Get(db, &dev, `
		select device_eui,
		protocol_type,
		created_at
		from device
		where device_eui=$1`,
		devEUI,
	)
	if err != nil {
		return dev, err
	}

	return dev, nil
}

// UpdateDevice update device
func UpdateDevice(dev Device) error {
	_, err := db.Exec(`
		update device set
		protocol_type=$2
		where device_eui=$1`,
		dev.DeviceEUI,
		dev.ProtocolType,
	)

	return err
}

// DeleteDevice delete device
func DeleteDevice(eui string) error {
	var devEUI EUI64
	err := devEUI.UnmarshalText([]byte(eui))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = delteDeviceInfo(tx, devEUI)
	if err != nil {
		return tx.Rollback()
	}

	err = tx.Commit()
	return err
}

func delteDeviceInfo(tx *sql.Tx, devEUI EUI64) (err error) {
	_, err = tx.Exec(`
		delete from device 
		where device_eui=$1`,
		devEUI,
	)
	if err != nil {
		return
	}

	return nil
}

// GetDevices get devices.
func GetDevices(limit, offset int) ([]Device, error) {
	var devs []Device
	err := sqlx.Select(db, &devs, `
		select d.device_eui,
		d.protocol_type,
		d.created_at
		from device d 
		limit $1 offset $2`,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	return devs, nil
}

// GetDevicesEUI get all devices eui
func GetDevicesEUI() ([]string, error) {
	var euis []EUI64
	err := sqlx.Select(db, &euis, `
		select device_eui
		from device`,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	devEUIS := make([]string, 0, len(euis))
	for _, eui := range euis {
		devEUIS = append(devEUIS, eui.String())
	}
	return devEUIS, nil
}

// GetDevicesCount get count of device.
func GetDevicesCount() (int, error) {
	var count int
	err := sqlx.Get(db, &count, `
		select count(device_eui) as cnt
		from device`,
	)

	if err != nil {
		return 0, err
	}

	return count, nil
}
