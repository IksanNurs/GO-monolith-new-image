// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

const TableNamePackage = "package"

// Package mapped from table <package>
type Package struct {
	ID            int32   `gorm:"column:id;type:int(11);primaryKey;autoIncrement:true" json:"id"`
	Name          string  `gorm:"column:name;type:varchar(255);not null" json:"name" form:"name"`
	Description   string `gorm:"column:description;type:text" json:"description" form:"description"`
	Duration      int32   `gorm:"column:duration;type:int(11);not null" json:"duration" form:"duration"`
	IsIndependent  int32   `gorm:"column:is_independent;type:int(11)" json:"is_independent" form:"is_independent"`
	IsComplete    int32   `gorm:"column:is_complete;type:int(11)" json:"is_complete" form:"is_complete"`
	Price         int32   `gorm:"column:price;type:int(11);not null" json:"price" form:"price"`
	VersionID     int32   `gorm:"column:version_id;type:int(11);not null;index:version_id,priority:1" json:"version_id" form:"version_id"`
	ScoringMethod int32  `gorm:"column:scoring_method;type:int(11);default:1" json:"scoring_method" form:"scoring_method"` // 1 = percentaged, 2  = summarized
	ShuffleType    int32  `gorm:"column:shuffle_type;type:int(11)" json:"shuffle_type" form:"shuffle_type"` // 1 = berurutan, 2 = acak keseluruhan, 3 = acak per section	
	MinimumScore  *int32  `gorm:"column:minimum_score;type:int(11)" json:"minimum_score" form:"minimum_score"`
	IsActive      int32  `gorm:"column:is_active;type:int(11);default:1" json:"is_active"`
	CreatedAt     *int32  `gorm:"column:created_at;type:int(11)" json:"created_at"`
	UpdatedAt     *int32  `gorm:"column:updated_at;type:int(11)" json:"updated_at"`
	CreatedBy     *int32  `gorm:"column:created_by;type:int(11);index:created_by,priority:1" json:"created_by"`
	UpdatedBy     *int32  `gorm:"column:updated_by;type:int(11);index:updated_by,priority:1" json:"updated_by"`
	Version        Version `gorm:"foreignkey:VersionID" json:"Version"`
}

type InputPackage struct {
	ID            int32   `gorm:"column:id;type:int(11);primaryKey;autoIncrement:true" json:"id"`
	Name          string  `gorm:"column:name;type:varchar(255);not null" json:"name" form:"name"`
	Description   *string `gorm:"column:description;type:text" json:"description" form:"description"`
	Duration      int32   `gorm:"column:duration;type:int(11);not null" json:"duration" form:"duration"`
	IsIndependent *int32   `gorm:"column:is_independent;type:int(11)" json:"is_independent" form:"is_independent"`
	IsComplete    *int32   `gorm:"column:is_complete;type:int(11)" json:"is_complete" form:"is_complete"`
	Price         *int32   `gorm:"column:price;type:int(11);" json:"price" form:"price"`
	VersionID     int32   `gorm:"column:version_id;type:int(11);not null;index:version_id,priority:1" json:"version_id" form:"version_id"`
	ScoringMethod int32  `gorm:"column:scoring_method;type:int(11);default:1" json:"scoring_method" form:"scoring_method"` // 1 = percentaged, 2  = summarized
	ShuffleType    int32  `gorm:"column:shuffle_type;type:int(11)" json:"shuffle_type" form:"shuffle_type"` // 1 = berurutan, 2 = acak keseluruhan, 3 = acak per section	
	MinimumScore  *int32  `gorm:"column:minimum_score;type:int(11)" json:"minimum_score" form:"minimum_score"`
	IsActive      *int32  `gorm:"column:is_active;type:int(11);default:1" json:"is_active" form:"is_active"`
	CreatedAt     *int32  `gorm:"column:created_at;type:int(11)" json:"created_at"`
	CreatedBy     int32  `gorm:"column:created_by;type:int(11);index:created_by,priority:1" json:"created_by"`
}


type PackageVersion struct {
	VersionID     int32   `gorm:"column:version_id;type:int(11);not null;index:version_id,priority:1" json:"version_id"`
}

type PackageSelect struct {
	ID            int32   `gorm:"column:id;type:int(11);primaryKey;autoIncrement:true" json:"id"`
	Name          string  `gorm:"column:name;type:varchar(255);not null" json:"text"`
}

type PackagePrice struct {
   Price         int32   `gorm:"column:price;type:int(11);not null" json:"price"`
}

func (i *Package) BeforeCreate(scope *gorm.DB) error {
	fmt.Println(i.CreatedAt)
    now := int32(time.Now().UTC().Unix())
    i.CreatedAt = &now
    return nil
}

func (i *InputPackage) BeforeCreate(scope *gorm.DB) error {
	fmt.Println(i.CreatedAt)
    now := int32(time.Now().UTC().Unix())
    i.CreatedAt = &now
    return nil
}
// TableName Package's table name
func (*Package) TableName() string {
	return TableNamePackage
}

func (*PackageVersion) TableName() string {
	return TableNamePackage
}

func (*PackageSelect) TableName() string {
	return TableNamePackage
}

func (*PackagePrice) TableName() string {
	return TableNamePackage
}

func (*InputPackage) TableName() string {
	return TableNamePackage
}




