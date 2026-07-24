package model

import (
	"github.com/uptrace/bun"
)

type LocationInput struct {
	MID int64   `json:"mapillary_id" binding:"required"`
	Lat    float64 `json:"latitude" binding:"required"`
	Long   float64 `json:"longitude" binding:"required"`
	CC string  `json:"country_code" binding:"required,len=2"`
}

type Location struct {
	bun.BaseModel `bun:"table:locations,alias:l"`

	ID int64 `bun:"id,pk,autoincrement"`
	MID int64 `bun:"id_mapillary,unique,notnull"`
	Lat float64 `bun:"latitude,notnull"`
	Long float64 `bun:"longitude,notnull"`
	CC string `bun:"country_code,type:varchar(2),notnull"`
}