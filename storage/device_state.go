package storage

import (
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// DeviceState define device state
type DeviceState struct {
	DeviceEUI  EUI64           `db:"device_eui"`
	LastSeenAt *time.Time      `db:"last_seen_at"`
	Location   *GPSPoint       `db:"location"`
	Detail     json.RawMessage `db:"detail"`
}

// GetDeviceState get device state.
func GetDeviceState(eui string) (DeviceState, error) {
	var devEUI EUI64
	var devState DeviceState
	err := devEUI.UnmarshalText([]byte(eui))
	if err != nil {
		return devState, err
	}

	err = sqlx.Get(db, &devState, `
		select device_eui,
		last_seen_at,
		location,
		detail
		from device_state
		where device_eui=$1`,
		devEUI,
	)
	if err != nil {
		return devState, err
	}

	return devState, nil
}

// CreateAndUpdateState create and update device state.
func CreateAndUpdateState(devState DeviceState) error {
	_, err := db.Exec(`
		insert into device_state(
			device_eui,
			last_seen_at,
			location,
			detail
		)values($1,$2,$3,$4)
		on conflict(device_eui)
		do update set 
		    last_seen_at=$2,
			location=$3,
			detail=$4`,
		devState.DeviceEUI,
		devState.LastSeenAt,
		devState.Location,
		devState.Detail,
	)

	if err != nil {
		return errors.Wrap(err, "insert or update device state error.")
	}
	return nil
}
