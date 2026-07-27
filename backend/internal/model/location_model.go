package model

import (
	"errors"

	"github.com/uptrace/bun"
)

var (
	ErrMissingFields      = errors.New("missing required fields")
	ErrInvalidLatitude    = errors.New("invalid latitude")
	ErrInvalidLongitude   = errors.New("invalid longitude")
	ErrInvalidCountryCode = errors.New("invalid country code")
)

// Location sert à la fois de modèle bun (colonnes DB) et de payload JSON pour les requêtes de création et de patch.
// Tous les champs de données sont des pointeurs pour pouvoir distinguer
// "champ non fourni" (nil) de "champ fourni" — nécessaire pour un PATCH partiel.
type Location struct {
	bun.BaseModel `bun:"table:locations,alias:l" json:"-"`

	ID   int64    `bun:"id,pk,autoincrement" json:"id,omitempty"`
	MID  *int64   `bun:"id_mapillary,unique,notnull" json:"mapillary_id,omitempty"`
	Lat  *float64 `bun:"latitude,notnull" json:"latitude,omitempty"`
	Long *float64 `bun:"longitude,notnull" json:"longitude,omitempty"`
	CC   *string  `bun:"country_code,type:varchar(2),notnull" json:"country_code,omitempty"`
}

// ValidateForCreate : utilisée sur POST — tous les champs doivent être présents et valides
func (l *Location) ValidateForCreate() error {
	if l.MID == nil || l.Lat == nil || l.Long == nil || l.CC == nil {
		return ErrMissingFields
	}
	return l.validateRanges()
}

// ValidateForPatch : utilisée sur PATCH — seuls les champs fournis (non-nil) sont validés
func (l *Location) ValidateForPatch() error {
	return l.validateRanges()
}

// validateRanges ne vérifie QUE les champs non-nil — comportement partagé create/patch
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