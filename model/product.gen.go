// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameProduct = "product"

// Product mapped from table <product>
type Product struct {
	ID             int32  `gorm:"column:id;type:int(11);primaryKey;autoIncrement:true" json:"id"`
	Name           string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	PriceNonmember int32  `gorm:"column:price_nonmember;type:int(11);not null" json:"price_nonmember"`
	PriceMember    int32  `gorm:"column:price_member;type:int(11);not null" json:"price_member"`
	Pv             int32  `gorm:"column:pv;type:int(11);not null" json:"pv"`
	Stock          int32  `gorm:"column:stock;type:int(11);not null" json:"stock"`
	CreatedAt int32 `gorm:"column:created_at;type:int(11);" json:"created_at"`
	CreatedAt_t string `gorm:"column:created_at;type:int(11);" json:"created_at_t"`
}


type InputProduct struct {
	ID             int32  `gorm:"column:id;type:int(11);primaryKey;autoIncrement:true" json:"id" form:"id"`
	Name           string `gorm:"column:name;type:varchar(255);not null" json:"name" form:"name"`
	PriceNonmember int32  `gorm:"column:price_nonmember;type:int(11);not null" json:"price_nonmember" form:"price_nonmember"`
	PriceMember    int32  `gorm:"column:price_member;type:int(11);not null" json:"price_member" form:"price_member"`
	Pv             int32  `gorm:"column:pv;type:int(11);not null" json:"pv" form:"pv"`
	Stock          *int32  `gorm:"column:stock;type:int(11);not null" json:"stock" form:"stock"`
	CreatedAt int32 `gorm:"column:created_at;type:int(11);autoCreateTime;loc:UTC" json:"created_at"`
	UpdatedAt int32 `gorm:"column:updated_at;type:int(11);autoUpdateTime;loc:UTC" json:"updated_at"`
}

type ProductSelect struct {
	ID            int32   `gorm:"column:id;type:int(11);primaryKey;autoIncrement:true" json:"id"`
	Name          string  `gorm:"column:name;type:varchar(255);not null" json:"text"`
}

// TableName Product's table name
func (*Product) TableName() string {
	return TableNameProduct
}

func (*InputProduct) TableName() string {
	return TableNameProduct
}

func (*ProductSelect) TableName() string {
	return TableNameProduct
}

