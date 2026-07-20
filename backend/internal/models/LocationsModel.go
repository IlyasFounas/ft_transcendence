package locations

import (
	"github.com/uptrace/bun"
)

type Location struct {
	bun.BaseModel `bun:"table:locations,alias:l"`

	ID int64 `bun:"id,pk,autoincrement"`
	ID_mplry int64 `bun:"id_mapillary,unique,notnull"`
	Lat float64 `bun:"latitude,notnull"`
	Long float64 `bun:"longitude,notnull"`
	CC string `bun:"country_code,type:varchar(2),notnull"`
}