package institution

type InstitutionFormatter struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Code           string `json:"code"`
	Type           int64  `json:"type"`
	Type_text      string `json:"type_text"`
	Province_id    string `json:"province_id"`
	District_id    string `json:"district_id"`
	Subdistrict_id string `json:"subdistrict_id"`
	Village_id     string `json:"village_id"`
	Address        string `json:"address"`
	Created_at     int64  `json:"created_at"`
	Updated_at     int64  `json:"updated_at"`
}

// func FormatInstitution(institution Institution) Institution {
// 	formatter := Institution{}

// 	return formatter
// }
