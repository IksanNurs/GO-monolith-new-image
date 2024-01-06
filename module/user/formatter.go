package usercamp

type UserFormatter struct {
	ID       int64  `json:"id" sql:"unique"`
	Username string `json:"username"`
	Email    string `json:"email" sql:"unique"`
	Role     int    `json:"role"`
	Token    string `json:"token"`
}
type UserFormatterById struct {
	ID       int64  `json:"id" sql:"unique"`
	Username string `json:"username"`
	Email    string `json:"email" sql:"unique"`
	Role     int    `json:"role"`
}
type UserFormatter1 struct {
	ID                    int64       `json:"id" sql:"unique"`
	UUID                  string      `json:"uuid" sql:"unique"`
	Username              string      `json:"username"`
	Email                 string      `json:"email" sql:"unique"`
	Phone                 string      `json:"phone" sql:"unique"`
	Name                  string      `json:"name"`
	Nickname              string      `json:"nickname"`
	Birthplace            string      `json:"birthplace"`
	Birthdate             string      `json:"birthdate"`
	Sex                   int64       `json:"sex"`
	Education_level       int         `json:"education_level"`
	Work_institution_name string      `json:"work_institution_name"`
	Institution_id        int64       `json:"education_institution_id"`
	Provinceid            int         `json:"address_province_id"`
	Districtid            int         `json:"address_district_id"`
	Subdistrictid         int         `json:"address_subdistrict_id"`
	Villageid             int         `json:"address_village_id"`
	Address               string      `json:"address_street"`
	Is_organizer          int         `json:"event_is_organizer"`
	Status                int         `json:"status"`
	Review_status         int         `json:"event_review_status"`
	Createdat             int64       `json:"created_at"`
	Updatedat             int64       `json:"updated_at"`
	MustChangePassword    int         `json:"must_change_password"`
	ConfirmedAt           int         `json:"confirmed_at"`
	Occupation            int         `json:"occupation"`
	Institution           interface{} `json:"educationInstitution"`
	Province              interface{} `json:"addressProvince"`
	District              interface{} `json:"addressDistrict"`
	Subdistrict           interface{} `json:"addressSubdistrict"`
	Village               interface{} `json:"addressVillage"`
}

func FormatUser(user User, token string) UserFormatter {
	formatter := UserFormatter{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	}

	return formatter
}

func FormatUserbyToken(user User) UserFormatter1 {

	formatter := UserFormatter1{
		ID:       int64(user.ID),
		UUID:     user.UUID.String,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone.String,
		// Fcm: user.FcmToken,
		Name:                  user.Name.String,
		Nickname:              user.Nickname.String,
		Birthplace:            user.Birthplace.String,
		Birthdate:             user.Birthdate.String,
		Sex:                   user.Sex.Int64,
		Education_level:       int(user.EducationLevel.Int64),
		Work_institution_name: user.WorkInstitutionName.String,
		Institution_id:        int64(user.EducationInstitutionID.Int64),
		Provinceid:            int(user.AddressProvinceID.Int64),
		Districtid:            int(user.AddressDistrictID.Int64),
		Subdistrictid:         int(user.AddressSubdistrictID.Int64),
		Villageid:             int(user.AddressVillageID.Int64),
		Address:               user.AddressStreet.String,
		Status:                int(user.Status.Int64),
		MustChangePassword:    int(user.MustChangePassword.Int64),
		ConfirmedAt:           int(user.ConfirmedAt.Int64),
		Occupation:            int(user.Occupation.Int64),
		Createdat:             int64(user.CreatedAt.Int64),
		Updatedat:             int64(user.CreatedAt.Int64),
	}

	return formatter
}

func FormatUserbyid(user User) UserFormatterById {
	formatter := UserFormatterById{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	return formatter
}
