// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameDistrict = "_district"

// District mapped from table <_district>
type District struct {
	ID         *string `gorm:"column:id;type:char(4)" json:"id"`
	ProvinceID *string `gorm:"column:province_id;type:char(2)" json:"province_id"`
	Name       *string `gorm:"column:name;type:varchar(191)" json:"name"`
	AreaTypeID *int32  `gorm:"column:area_type_id;type:int(11)" json:"area_type_id"`
	IsActive   *int32  `gorm:"column:is_active;type:int(11)" json:"is_active"`
	CreatedAt  *int32  `gorm:"column:created_at;type:int(11)" json:"created_at"`
	UpdatedAt  *int32  `gorm:"column:updated_at;type:int(11)" json:"updated_at"`
	CreatedBy  *int32  `gorm:"column:created_by;type:int(11)" json:"created_by"`
	UpdatedBy  *int32  `gorm:"column:updated_by;type:int(11)" json:"updated_by"`
}

// TableName District's table name
func (*District) TableName() string {
	return TableNameDistrict
}
