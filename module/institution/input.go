package institution

type FormCreateInstitution struct {
	Name        string `form:"name" json:"name" binding:"required"`
	Type        int64  `form:"type" json:"type"`
	Code        string `form:"code" json:"code"`
	Province_ID string `form:"province_id" json:"province_id"`
	Province_Name string `form:"province_name" json:"province_name"`
	Error       error
}

// type FormUpdateInstitution struct {
// 	ID          int    `form:"id" json:"id" binding:"required"`
// 	Name        string `form:"name" json:"name" binding:"required"`
// 	Type        int64  `form:"type" json:"type"`
// 	Code        string `form:"code" json:"code"`
// 	Province_ID string `form:"province_id" json:"province_id"`
// }
