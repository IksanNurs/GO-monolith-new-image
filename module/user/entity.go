package usercamp

import "database/sql"

type User struct {
	ID   int64
	UUID sql.NullString
	// FcmToken               string
	Password               string
	Name                   sql.NullString `json:"name"`
	AuthKey                sql.NullString `json:"auth_key"`
	Phone                  sql.NullString `json:"phone"`
	Email                  string         `json:"email"`
	Username               string         `json:"username"`
	Status                 sql.NullInt64  `json:"status"`
	Nickname               sql.NullString `json:"nickname"`
	Birthplace             sql.NullString `json:"birthplace"`
	Birthdate              sql.NullString `json:"birthdate"`
	Sex                    sql.NullInt64  `json:"sex"`
	AddressProvinceID      sql.NullInt64  `json:"address_province_id"`
	AddressDistrictID      sql.NullInt64  `json:"address_district_id"`
	AddressSubdistrictID   sql.NullInt64  `json:"address_subdistrict_id"`
	AddressVillageID       sql.NullInt64  `json:"address_village_id"`
	AddressStreet          sql.NullString `json:"address_street"`
	EducationLevel         sql.NullInt64  `json:"education_level"`
	EducationInstitutionID sql.NullInt64  `json:"education_institution_id"`
	MustChangePassword     sql.NullInt64  `json:"must_change_password"`
	ConfirmedAt            sql.NullInt64  `json:"confirmed_at"`
	Occupation             sql.NullInt64  `json:"occupation"`
	WorkInstitutionName    sql.NullString `json:"work_institution_name"`
	CreatedAt              sql.NullInt64  `json:"created_at"`
	UpdatedAt              sql.NullInt64  `json:"updated_at"`
	Institution_name       sql.NullString `json:"institution_name"`
	Province_name          sql.NullString `json:"province_name"`
	District_name          sql.NullString `json:"district_name"`
	Subdistrict_name       sql.NullString `json:"subdistrict_name"`
	Village_name           sql.NullString `json:"village_name"`
	Role                   sql.NullInt64
	// Province               Province `gorm:"foreignKey:AddressProvinceID"`
	// District               District `gorm:"foreignKey:AddressDistrictID"`
}

type User1 struct {
	ID               int            `json:"id"`
	Name             sql.NullString `json:"name"`
	AuthKey          sql.NullString `json:"auth_key"`
	Phone            sql.NullString `json:"phone"`
	Email            string         `json:"email"`
	Username         string         `json:"username"`
	Institution_name sql.NullString `json:"institution_name"`
	Status           string         `json:"status"`
	// ...
}
