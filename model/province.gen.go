// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameProvince = "_province"

// Province mapped from table <_province>
type Province struct {
	ID        *string `gorm:"column:id;type:char(2)" json:"id"`
	Name      *string `gorm:"column:name;type:tinytext" json:"name"`
	IsActive  *int32  `gorm:"column:is_active;type:int(11)" json:"is_active"`
	CreatedAt *int32  `gorm:"column:created_at;type:int(11)" json:"created_at"`
	UpdatedAt *int32  `gorm:"column:updated_at;type:int(11)" json:"updated_at"`
	CreatedBy *int32  `gorm:"column:created_by;type:int(11)" json:"created_by"`
	UpdatedBy *int32  `gorm:"column:updated_by;type:int(11)" json:"updated_by"`
}

// TableName Province's table name
func (*Province) TableName() string {
	return TableNameProvince
}
