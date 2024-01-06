// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNamePackageQuestion = "package_question"

// PackageQuestion mapped from table <package_question>
type PackageQuestion struct {
	ID         int32  `gorm:"column:id;type:int(11);primaryKey;autoIncrement:true" json:"id"`
	PackageID  int32  `gorm:"column:package_id;type:int(11);not null;uniqueIndex:package_id_question_id,priority:1" json:"package_id"`
	QuestionID int32  `gorm:"column:question_id;type:int(11);not null;uniqueIndex:package_id_question_id,priority:2;index:question_id,priority:1" json:"question_id"`
	SectionID   int32   `gorm:"column:section_id;type:int(11);not null;index:section_id,priority:1" json:"section_id"`
	Sequence   *int32 `gorm:"column:sequence;type:int(11)" json:"sequence"`
	IsDelete   int32 `gorm:"-" json:"is_delete"`
	CreatedAt  *int32 `gorm:"column:created_at;type:int(11)" json:"created_at"`
	UpdatedAt  *int32 `gorm:"column:updated_at;type:int(11)" json:"updated_at"`
	CreatedBy  *int32 `gorm:"column:created_by;type:int(11);index:created_by,priority:1" json:"created_by"`
	UpdatedBy  *int32 `gorm:"column:updated_by;type:int(11);index:updated_by,priority:1" json:"updated_by"`
    Question     Question `gorm:"foreignkey:QuestionID" json:"Question"`
    Section     Section `gorm:"foreignkey:SectionID" json:"Section"`
}

type PackageQuestion1 struct {
	ID         int32  `gorm:"column:id;type:int(11);primaryKey;autoIncrement:true" json:"id"`
	PackageID  int32  `gorm:"column:package_id;type:int(11);not null;uniqueIndex:package_id_question_id,priority:1" json:"package_id"`
	QuestionID int32  `gorm:"column:question_id;type:int(11);not null;uniqueIndex:package_id_question_id,priority:2;index:question_id,priority:1" json:"question_id"`
	SectionID   int32   `gorm:"column:section_id;type:int(11);not null;index:section_id,priority:1" json:"section_id"`
	Sequence   *int32 `gorm:"column:sequence;type:int(11)" json:"sequence"`
	CreatedAt  int32 `gorm:"column:created_at;type:int(11)autoCreateTime;loc:UTC" json:"created_at"`
	UpdatedAt  int32 `gorm:"column:updated_at;type:int(11)autoUpdateTime;loc:UTC" json:"updated_at"`
	CreatedBy  *int32 `gorm:"column:created_by;type:int(11);index:created_by,priority:1" json:"created_by"`
	UpdatedBy  *int32 `gorm:"column:updated_by;type:int(11);index:updated_by,priority:1" json:"updated_by"`
}

type InputPackageQuestion struct {
	PackageID  int32  `gorm:"column:package_id;type:int(11);not null;uniqueIndex:package_id_question_id,priority:1" json:"package_id"`
	QuestionID int32  `gorm:"column:question_id;type:int(11);not null;uniqueIndex:package_id_question_id,priority:2;index:question_id,priority:1" json:"question_id"`
	SectionID   int32   `gorm:"column:section_id;type:int(11);not null;index:section_id,priority:1" json:"section_id"`
	Sequence   *int32 `gorm:"column:sequence;type:int(11)" json:"sequence"`
	CreatedAt  *int32 `gorm:"column:created_at;type:int(11)" json:"created_at"`
	CreatedBy  int32 `gorm:"column:created_by;type:int(11);index:created_by,priority:1" json:"created_by"`
}


type PackageQuestionCount struct {
	ID           int32  `gorm:"column:id;type:int(11);primaryKey;autoIncrement:true" json:"id"`
    SectionID   int32   `gorm:"column:section_id;type:int(11);not null;index:section_id,priority:1" json:"section_id"`  
	Count        int32  `gorm:"column:count;type:int(11);" json:"count"`
	Section     Section `gorm:"foreignkey:SectionID" json:"Section"`
}

type PackageQuestion3 struct {
	ID          int32   `gorm:"column:id;type:int(11);primaryKey;autoIncrement:true" json:"id"`
	PackageID  int32  `gorm:"column:package_id;type:int(11);not null;uniqueIndex:package_id_question_id,priority:1" json:"package_id"`
	QuestionID int32  `gorm:"column:question_id;type:int(11);not null;uniqueIndex:package_id_question_id,priority:2;index:question_id,priority:1" json:"question_id"`
	Sequence   int32 `gorm:"column:sequence;type:int(11)" json:"sequence"`
	SectionID   int32   `gorm:"column:section_id;type:int(11);not null;index:section_id,priority:1" json:"section_id"`
	Package     Package `gorm:"foreignkey:PackageID" json:"Package"`

}

func (i *InputPackageQuestion) BeforeCreate(scope *gorm.DB) error {
    now := int32(time.Now().UTC().Unix())
    i.CreatedAt = &now
    return nil
}
// TableName PackageQuestion's table name
func (*PackageQuestion) TableName() string {
	return TableNamePackageQuestion
}

func (*PackageQuestion1) TableName() string {
	return TableNamePackageQuestion
}

func (*PackageQuestion3) TableName() string {
	return TableNamePackageQuestion
}

func (*InputPackageQuestion) TableName() string {
	return TableNamePackageQuestion
}


func (*PackageQuestionCount) TableName() string {
	return TableNamePackageQuestion
}
