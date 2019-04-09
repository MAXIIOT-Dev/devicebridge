package storage

import (
	"time"

	"github.com/jmoiron/sqlx"
)

// DeviceTrack device gps track
type DeviceTrack struct {
	DeviceEUI EUI64     `db:"device_eui"`
	CreatedAt time.Time `db:"created_at"`
	Location  GPSPoint  `db:"location"`
	Altitude  uint16    `db:"altitude"`
}

// GetDeviceTrack returns device tracks
func GetDeviceTrack(eui string, start, end time.Time) ([]DeviceTrack, error) {
	var devEUI EUI64
	var tracks []DeviceTrack
	if err := devEUI.UnmarshalText([]byte(eui)); err != nil {
		return nil, err
	}

	err := sqlx.Select(db, &tracks, `
		select device_eui,
		created_at,
		location,
		altitude
		from device_track
		where device_eui=$1
		and created_at between $2 and $3`,
		devEUI,
		start,
		end,
	)
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

// CreateDeviceTrack create device track
func CreateDeviceTrack(track DeviceTrack) error {
	_, err := db.Exec(`
		insert into device_track(
			device_eui,
			created_at,
			location,
			altitude
		)values($1,$2,$3,$4)`,
		track.DeviceEUI,
		track.CreatedAt,
		track.Location,
		track.Altitude,
	)

	return err
}
