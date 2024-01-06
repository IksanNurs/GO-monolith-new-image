package model

const TableNameInstitution = "_institution"

// Institution mapped from table <_institution>
type Institution struct {
	ID            *int32  `gorm:"column:id;type:int(11)" json:"id"`
	Name          *string `gorm:"column:name;type:varchar(255)" json:"name"`
	Type          *int32  `gorm:"column:type;type:int(11)" json:"type"`
	Code          *string `gorm:"column:code;type:varchar(255)" json:"code"`
	ProvinceID    *string `gorm:"column:province_id;type:char(2)" json:"province_id"`
	DistrictID    *string `gorm:"column:district_id;type:char(4)" json:"district_id"`
	SubdistrictID *string `gorm:"column:subdistrict_id;type:char(6)" json:"subdistrict_id"`
	VillageID     *string `gorm:"column:village_id;type:char(10)" json:"village_id"`
	Address       *string `gorm:"column:address;type:varchar(255)" json:"address"`
	MergedTo      *int32  `gorm:"column:merged_to;type:int(11)" json:"merged_to"`
	ProvinceName  *string `gorm:"column:province_name;type:varchar(255)" json:"province_name"`
	CreatedAt     *int32  `gorm:"column:created_at;type:int(11)" json:"created_at"`
	UpdatedAt     *int32  `gorm:"column:updated_at;type:int(11)" json:"updated_at"`
	CreatedBy     *int32  `gorm:"column:created_by;type:int(11)" json:"created_by"`
	UpdatedBy     *int32  `gorm:"column:updated_by;type:int(11)" json:"updated_by"`
}

// TableName Institution's table name
func (*Institution) TableName() string {
	return TableNameInstitution
}
