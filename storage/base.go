package storage

import (
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
)

// EUI64 device eui
type EUI64 [8]byte

// Scan implements sql.Scanner.
func (e *EUI64) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.New("[]byte type expected")
	}
	if len(b) != len(e) {
		return fmt.Errorf("[]byte must have length %d", len(e))
	}
	copy(e[:], b)
	return nil
}

// Value implements driver.Valuer.
func (e EUI64) Value() (driver.Value, error) {
	return e[:], nil
}

// UnmarshalText unmarshal string to EUI64
func (e *EUI64) UnmarshalText(eui []byte) error {
	hexEUI, err := hex.DecodeString(string(eui))
	if err != nil {
		return err
	}
	if len(e) != len(hexEUI) {
		return fmt.Errorf("eui must have length %d", len(e)*2)
	}
	copy(e[:], hexEUI)
	return nil
}

// MarshalText implement json marshal
func (e EUI64) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(e[:])), nil
}

// String return hex format
func (e EUI64) String() string {
	return hex.EncodeToString(e[:])
}

// GPSPoint contains a GPS point.
type GPSPoint struct {
	Latitude  float64
	Longitude float64
}

// Value implements the driver.Valuer interface
func (p GPSPoint) Value() (driver.Value, error) {
	return fmt.Sprintf("(%s,%s)",
		strconv.FormatFloat(p.Latitude, 'f', -1, 64),
		strconv.FormatFloat(p.Longitude, 'f', -1, 64),
	), nil
}

// Scan implements the sql.Scanner interface.
func (p *GPSPoint) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte, got %T", src)
	}

	_, err := fmt.Sscanf(string(b), "(%f,%f)", &p.Latitude, &p.Longitude)
	return err
}
