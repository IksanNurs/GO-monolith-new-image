package institution

import (
	"database/sql"
	"time"
)

type Institution struct {
	ID             int
	Name           sql.NullString
	Code           sql.NullString
	Type           sql.NullInt64
	Type_text      sql.NullString
	Province_name  sql.NullString
	Province_id    sql.NullString
	District_id    sql.NullString
	Subdistrict_id sql.NullString
	Village_id     sql.NullString
	Address        sql.NullString
	Created_at     sql.NullInt64
	Updated_at     sql.NullInt64
	// Province       Province
	// District       District
}

type ContentUser struct {
	ID           int64 `url:"id" json:"id"`
	IDUser       int64
	IDCategory   int64
	Title        string
	Subtitle     string
	Deksripsi    string
	Path         string
	LastModified time.Time
	Username     string
	Email        string
}

type InstitutionDatatable struct {
	ID            int            `json:"id"`
	Name          sql.NullString `json:"name"`
	Code          sql.NullString `json:"code"`
	Type          sql.NullInt64  `json:"type"`
	Province_name sql.NullString `json:"province_name"`
	// Province       Province
	// District       District
}
