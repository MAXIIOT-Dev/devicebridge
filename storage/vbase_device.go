package storage

import (
	"github.com/jmoiron/sqlx"
)

const (
	VDSOnline  = "online"
	VDSOffline = "offline"
)

// VbaseDevice vbase device
type VbaseDevice struct {
	DeviceEUI EUI64     `db:"device_eui"`
	Name      string    `db:"device_name"`
	Location  *GPSPoint `db:"location"`
}

// GetVbaseDevices returns vbase devices.
func GetVbaseDevices(status string) ([]VbaseDevice, error) {
	var devs []VbaseDevice
	if status == VDSOnline {
		err := sqlx.Select(db, &devs, `
		select d.device_eui,
		d.device_name,
		ds.location
		from device d
		left join device_state ds
		on d.device_eui=ds.device_eui
		where ds.last_seen_at > (current_timestamp - (5*interval '1 minute'))`,
		)
		if err != nil {
			return nil, err
		}
	} else {
		err := sqlx.Select(db, &devs, `
		select d.device_eui,
		d.device_name,
		ds.location
		from device d
		left join device_state ds
		on d.device_eui=ds.device_eui
		where coalesce(ds.last_seen_at,'1970-01-01') <= (current_timestamp - (5*interval '1 minute'))`,
		)
		if err != nil {
			return nil, err
		}
	}

	return devs, nil
}

// GetVbaseDevicesCount returns vbase device count
func GetVbaseDevicesCount(status string) (int, error) {
	var count int
	if status == VDSOnline {
		err := sqlx.Get(db, &count, `
		select count(1) cnt
		from device d
		left join device_state ds
		on d.device_eui=ds.device_eui
		where ds.last_seen_at > (current_timestamp - (5*interval '1 minute'))`,
		)
		if err != nil {
			return 0, err
		}
	} else {
		err := sqlx.Get(db, &count, `
		select count(1) cnt
		from device d
		left join device_state ds
		on d.device_eui=ds.device_eui
		where coalesce(ds.last_seen_at,'1970-01-01') <= (current_timestamp - (5*interval '1 minute'))`,
		)
		if err != nil {
			return 0, err
		}
	}

	return count, nil
}
