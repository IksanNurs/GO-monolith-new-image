package model

import (
	"time"
)

const TableNameUser = "user"

// User mapped from table <_user>
type User struct {
	ID                     *int32     `gorm:"column:id;type:int(11)" json:"id"`
	UUID                   *string    `gorm:"column:uuid;type:varchar(36)" json:"uuid"`
	Name                   *string    `gorm:"column:name;type:varchar(255)" json:"name"`
	Phone                  *string    `gorm:"column:phone;type:varchar(255)" json:"phone"`
	Email                  *string    `gorm:"column:email;type:varchar(255)" json:"email"`
	Username               *string    `gorm:"column:username;type:varchar(255)" json:"username"`
	AuthKey                *string    `gorm:"column:auth_key;type:text" json:"auth_key"`
	PasswordHash           *string    `gorm:"column:password_hash;type:varchar(255)" json:"password_hash"`
	PasswordResetToken     *string    `gorm:"column:password_reset_token;type:varchar(255)" json:"password_reset_token"`
	RefreshToken           *string    `gorm:"column:refresh_token;type:varchar(255)" json:"refresh_token"`
	FcmToken               *string    `gorm:"column:fcm_token;type:varchar(255)" json:"fcm_token"`
	DeviceToken            *string    `gorm:"column:device_token;type:varchar(255)" json:"device_token"`
	VerificationToken      *string    `gorm:"column:verification_token;type:varchar(255)" json:"verification_token"`
	OneTimePassword        *string    `gorm:"column:one_time_password;type:varchar(255)" json:"one_time_password"`
	OtpExpiredAt           *int32     `gorm:"column:otp_expired_at;type:int(11)" json:"otp_expired_at"`
	MustChangePassword     *int32     `gorm:"column:must_change_password;type:int(11)" json:"must_change_password"`
	Status                 *int32     `gorm:"column:status;type:int(11)" json:"status"`
	Nickname               *string    `gorm:"column:nickname;type:varchar(255)" json:"nickname"`
	Birthplace             *string    `gorm:"column:birthplace;type:varchar(255)" json:"birthplace"`
	Birthdate              *time.Time `gorm:"column:birthdate;type:date" json:"birthdate"`
	Sex                    *int32     `gorm:"column:sex;type:int(11)" json:"sex"`
	AddressProvinceID      *int32     `gorm:"column:address_province_id;type:int(11)" json:"address_province_id"`
	AddressDistrictID      *int32     `gorm:"column:address_district_id;type:int(11)" json:"address_district_id"`
	AddressSubdistrictID   *int32     `gorm:"column:address_subdistrict_id;type:int(11)" json:"address_subdistrict_id"`
	AddressVillageID       *int32     `gorm:"column:address_village_id;type:int(11)" json:"address_village_id"`
	AddressStreet          *string    `gorm:"column:address_street;type:text" json:"address_street"`
	EducationLevel         *int32     `gorm:"column:education_level;type:int(11)" json:"education_level"`
	EducationInstitutionID *int32     `gorm:"column:education_institution_id;type:int(11)" json:"education_institution_id"`
	Occupation             *int32     `gorm:"column:occupation;type:int(11)" json:"occupation"`
	OccupationDescription  *string    `gorm:"column:occupation_description;type:varchar(255)" json:"occupation_description"`
	WorkInstitutionName    *string    `gorm:"column:work_institution_name;type:varchar(255)" json:"work_institution_name"`
	ClientOrigin           *string    `gorm:"column:client_origin;type:text" json:"client_origin"`
	ConfirmedAt            *int32     `gorm:"column:confirmed_at;type:int(11)" json:"confirmed_at"`
	CreatedAt              *int32     `gorm:"column:created_at;type:int(11)" json:"created_at"`
	UpdatedAt              *int32     `gorm:"column:updated_at;type:int(11)" json:"updated_at"`
	CreatedBy              *int32     `gorm:"column:created_by;type:int(11)" json:"created_by"`
	UpdatedBy              *int32     `gorm:"column:updated_by;type:int(11)" json:"updated_by"`
}

type SubUser struct {
	ID                     *int32      `gorm:"column:id;type:int(11)" json:"id"`
	Name                   *string     `gorm:"column:name;type:varchar(255)" json:"name"`
	Phone                  *string     `gorm:"column:phone;type:varchar(255)" json:"phone"`
	Email                  *string     `gorm:"column:email;type:varchar(255)" json:"email"`
	EducationInstitutionID *int32      `gorm:"column:education_institution_id;type:int(11)" json:"education_institution_id"`
	Institution            Institution `gorm:"foreignkey:EducationInstitutionID" json:"Institution"`
}

type UserSelect struct {
	ID    int32  `gorm:"column:id;type:int(11);primaryKey;autoIncrement:true" json:"id"`
	Email string `gorm:"column:email;type:varchar(255);not null" json:"text"`
}

type InputUser struct {
	ID        *int32  `gorm:"column:id;type:int(11)" json:"id" form:"id"`
	Name      *string `gorm:"column:name;type:varchar(255)" json:"name" form:"name"`
	Email     *string `gorm:"column:email;type:varchar(255)" json:"email" form:"email"`
}

// TableName User's table name
func (*User) TableName() string {
	return TableNameUser
}
func (*InputUser) TableName() string {
	return TableNameUser
}

func (*SubUser) TableName() string {
	return TableNameUser
}

func (*UserSelect) TableName() string {
	return TableNameUser
}
