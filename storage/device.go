package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

// Device define device model
type Device struct {
	DeviceEUI EUI64  `db:"device_eui" json:"device_eui"`
	Name      string `db:"device_name" json:"device_name"`
	Icon      string `db:"icon" json:"icon"`
	Status    string `db:"status" json:"status"`
}

// CreateDevice create device on database
func CreateDevice(dev Device) error {
	now := time.Now()
	_, err := db.Exec(`
		insert into device (
			device_eui,
			device_name,
			icon,
			created_at,
			updated_at
		)values($1,$2,$3,$4,$4)`,
		dev.DeviceEUI,
		dev.Name,
		dev.Icon,
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
		device_name,
		icon
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
		device_name=$2,
		icon=$3
		where device_eui=$1`,
		dev.DeviceEUI,
		dev.Name,
		dev.Icon,
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

	_, err = tx.Exec(`
	delete from device_track
	where device_eui=$1`,
		devEUI,
	)
	if err != nil {
		return
	}

	_, err = tx.Exec(`
		delete from device_state
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
		d.device_name,
		d.icon,
		case when ds.last_seen_at > (current_timestamp - (5*interval '1 minute')) then 'online'
		     else 'offline' end as status
		from device d
		left join device_state ds 
		on d.device_eui=ds.device_eui
		limit $1 offset $2`,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	return devs, nil
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
