package locations

import (
	"errors"
	"time"

	"github.com/uptrace/bun"
)

var (
	ErrMissingFields      = errors.New("missing required fields")
	ErrInvalidLatitude    = errors.New("invalid latitude")
	ErrInvalidLongitude   = errors.New("invalid longitude")
	ErrInvalidCountryCode = errors.New("invalid country code")
)

// Location is a struct that serves as both a bun model (DB columns) and a JSON payload for create and patch requests.
// All data fields are pointers to be able to distinguish "field not provided" (nil) from "field provided" — necessary for a partial PATCH.
type Location struct {
	bun.BaseModel `bun:"table:locations,alias:l" json:"-"`

	ID   int64    `bun:"location_id,pk,autoincrement" json:"location_id,omitempty"`
	MID  *int64   `bun:"mapillary_id,unique,notnull" json:"mapillary_id,omitempty"`
	Lat  *float64 `bun:"latitude,notnull" json:"latitude,omitempty"`
	Long *float64 `bun:"longitude,notnull" json:"longitude,omitempty"`
	CC   *string  `bun:"country_code,type:varchar(2),notnull" json:"country_code,omitempty"`
	DisabledAt *time.Time `bun:"disabled_at" json:"disabled_at,omitempty"`
}

// ValidateForCreate : used for POST — all fields must be present and valid
func (l *Location) ValidateForCreate() error {
	if l.MID == nil || l.Lat == nil || l.Long == nil || l.CC == nil {
		return ErrMissingFields
	}
	return l.validateRanges()
}

// ValidateForPatch : used for PATCH — only the provided fields (non-nil) are validated
func (l *Location) ValidateForPatch() error {
	return l.validateRanges()
}

// validateRanges only checks the non-nil fields — shared behavior between create and patch
func (l *Location) validateRanges() error {
	if l.Lat != nil && (*l.Lat < -90 || *l.Lat > 90) {
		return ErrInvalidLatitude
	}
	if l.Long != nil && (*l.Long < -180 || *l.Long > 180) {
		return ErrInvalidLongitude
	}
	if l.CC != nil && len(*l.CC) != 2 {
		return ErrInvalidCountryCode
	}
	return nil
}