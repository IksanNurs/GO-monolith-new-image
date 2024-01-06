// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameSubmission = "submission"

// Submission mapped from table <submission>
type Submission struct {
	UserID          int32   `gorm:"column:user_id;type:int(11);not null" json:"user_id"`
	Name                 *string `json:"name"`
	Phone                *string `json:"phone"`
	Email                *string `json:"email"`
	Formation       *string `gorm:"column:formation;type:varchar(255)" json:"formation"`
	FormationQuota  *int32  `gorm:"column:formation_quota;type:int(11)" json:"formation_quota"`
	InstitutionID  *int32  `gorm:"column:institution_id;type:int(11)" json:"institution_id"`
	InstitutionName string `json:"institution_name"`
	ProvinceID      *int32  `gorm:"column:province_id;type:int(11)" json:"province_id"`
	ProvinceName     string `json:"province_name"`
	DistrictID      *int32  `gorm:"column:district_id;type:int(11)" json:"district_id"`
	DistrictName     string `json:"district_name"`
	BasicScore      *int32  `gorm:"column:basic_score;type:int(11)" json:"basic_score"`
	BasicRank       *int32  `gorm:"column:basic_rank;type:int(11)" json:"basic_rank"`
	AdvancedScore   *int32  `gorm:"column:advanced_score;type:int(11)" json:"advanced_score"`
	AdvancedRank    *int32  `gorm:"column:advanced_rank;type:int(11)" json:"advanced_rank"`
	IsPublic       int32   `gorm:"column:is_public;type:int(11);not null" json:"is_public"` // 0=private 1=public
	CreatedAt       *int32  `gorm:"column:created_at;type:int(11)" json:"created_at"`
	UpdatedAt       *int32  `gorm:"column:updated_at;type:int(11)" json:"updated_at"`
	CreatedBy       *int32  `gorm:"column:created_by;type:int(11)" json:"created_by"`
	UpdatedBy       *int32  `gorm:"column:updated_by;type:int(11)" json:"updated_by"`
	UserFromID       SubUser `gorm:"foreignkey:UserID" json:"UserFromID"`
	Institution      Institution `gorm:"foreignkey:InstitutionID" json:"Institution"`
	Province       Province `gorm:"foreignkey:ProvinceID" json:"Province"`
	District       District `gorm:"foreignkey:DistrictID" json:"District"`
}

type UpdateSubmission struct {
	UserID          int32   `gorm:"column:user_id;type:int(11);not null" json:"user_id"`
	Formation       *string `gorm:"column:formation;type:varchar(255)" json:"formation"`
	FormationQuota  *int32  `gorm:"column:formation_quota;type:int(11)" json:"formation_quota"`
	InstitutionID  *int32  `gorm:"column:institution_id;type:int(11)" json:"institution_id"`
	ProvinceID      *int32  `gorm:"column:province_id;type:int(11)" json:"province_id"`
	DistrictID      *int32  `gorm:"column:district_id;type:int(11)" json:"district_id"`
	BasicScore      *int32  `gorm:"column:basic_score;type:int(11)" json:"basic_score"`
	BasicRank       *int32  `gorm:"column:basic_rank;type:int(11)" json:"basic_rank"`
	AdvancedScore   *int32  `gorm:"column:advanced_score;type:int(11)" json:"advanced_score"`
	AdvancedRank    *int32  `gorm:"column:advanced_rank;type:int(11)" json:"advanced_rank"`
	IsPublic       int32   `gorm:"column:is_public;type:int(11);not null" json:"is_public"` 
	CreatedAt       *int32  `gorm:"column:created_at;type:int(11)" json:"created_at"`
	UpdatedAt       *int32  `gorm:"column:updated_at;type:int(11)" json:"updated_at"`
	CreatedBy       *int32  `gorm:"column:created_by;type:int(11)" json:"created_by"`
	UpdatedBy       *int32  `gorm:"column:updated_by;type:int(11)" json:"updated_by"`
}

// TableName Submission's table name
func (*Submission) TableName() string {
	return TableNameSubmission
}

func (*UpdateSubmission) TableName() string {
	return TableNameSubmission
}
