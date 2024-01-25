package model

const TableNameUser = "user"

// User mapped from table <_user>

type SubUser struct {
	ID                     *int32      `gorm:"column:id;type:int(11)" json:"id"`
	Name                   *string     `gorm:"column:name;type:varchar(255)" json:"name"`
	Phone                  *string     `gorm:"column:phone;type:varchar(255)" json:"phone"`
	Email                  *string     `gorm:"column:email;type:varchar(255)" json:"email"`
	EducationInstitutionID *int32      `gorm:"column:education_institution_id;type:int(11)" json:"education_institution_id"`
	
}

type UserSelect struct {
	ID    *int32  `gorm:"column:id;type:int(15)" json:"id"`
    Name  *string `gorm:"column:name;type:varchar(255)" json:"text"`
}

type InputUser struct {
	ID    *int32  `gorm:"column:id;type:int(15)" json:"id" form:"id"`
	Name  *string `gorm:"column:name;type:varchar(255)" json:"name" form:"name"`
	Email *string `gorm:"column:email;type:varchar(255)" json:"email" form:"email"`
	Alamat   string `gorm:"column:alamat;type:varchar(255);not null" json:"alamat" form:"alamat"`
	State    string `gorm:"column:state;type:varchar(255);not null" json:"state" form:"state"`
	JoinDate string `gorm:"column:join_date;type:varchar(255);not null" json:"join_date" form:"join_date"`
	Phone    string `gorm:"column:phone;type:varchar(255);not null" json:"phone" form:"phone"`
	/*
		1=admin, 2=member, 3=nonmember

	*/
	CategoriUser int32 `gorm:"column:categori_user;type:int(11);not null" json:"categori_user" form:"categori_user"`
}

type User struct {
	ID        int32 `gorm:"column:id;type:int(11);primaryKey" json:"id"`
	Email    string `gorm:"column:email;type:varchar(255);primaryKey;uniqueIndex:email,priority:1" json:"email"`
	Name     string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Alamat   string `gorm:"column:alamat;type:varchar(255);not null" json:"alamat"`
	State    string `gorm:"column:state;type:varchar(255);not null" json:"state"`
	JoinDate string `gorm:"column:join_date;type:varchar(255);not null" json:"join_date"`
	Phone    string `gorm:"column:phone;type:varchar(255);not null" json:"phone"`
	/*
		1=admin, 2=member, 3=nonmember

	*/
	CategoriUser int32 `gorm:"column:categori_user;type:int(11);not null" json:"categori_user"`
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
