package usercamp

import "database/sql"

type Repository interface {
	Save(user User) (User, error)
	Updateuserrepo(user User) (User, error)
	FindByEmail(email string) (User, error)
	FindByID(ID int) (User, error)
	FetchAllUser() ([]User, error)
	Delete(ID int) (User, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(user User) (User, error) {
	var sqlStmt string = "INSERT INTO user (uuid, username, name, password_hash, email, event_is_organizer) VALUES (?, ?, ?, ?, ?, ?);"

	_, err := r.db.Exec(sqlStmt, user.UUID, user.Username, user.Name.String, user.Password, user.Email, user.Role)

	if err != nil {
		return user, err
	}
	// sqlStmt = "SELECT * FROM users ORDER BY id DESC LIMIT 1"

	// row := r.db.QueryRow(sqlStmt)
	// err = row.Scan(
	// 	&user.ID,
	// 	&user.Username,
	// 	&user.Password,
	// 	&user.Email,
	// 	&user.Role,
	// )
	// if err != nil {
	// 	return user, err
	// }

	return user, nil
}
func (r *repository) FindByEmail(email string) (User, error) {
	var user User
	var sqlStmt string = "SELECT * FROM users WHERE email= ?"

	row := r.db.QueryRow(sqlStmt, email)
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.Role,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}
func (r *repository) FindByID(ID int) (User, error) {
	var user User

	var sqlStmt string = "SELECT u.id, u.name, email, username, phone, status, nickname, birthplace, birthdate, sex, address_province_id, address_district_id, address_subdistrict_id, address_village_id, address_street, education_level, education_institution_id, must_change_password, confirmed_at, occupation, work_institution_name, u.created_at, u.updated_at, auth_key, i.name, p.name, d.name, s.name, v.name FROM user u LEFT JOIN village v ON u.address_village_id=v.id LEFT JOIN subdistrict s ON u.address_subdistrict_id=s.id LEFT JOIN district d ON u.address_district_id=d.id LEFT JOIN province p ON u.address_province_id=p.id LEFT JOIN institution i ON u.education_institution_id=i.id WHERE u.id= ?"

	row := r.db.QueryRow(sqlStmt, ID)
	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Username,
		&user.Phone,
		&user.Status,
		&user.Nickname,
		&user.Birthplace,
		&user.Birthdate,
		&user.Sex,
		&user.AddressDistrictID,
		&user.AddressDistrictID,
		&user.AddressSubdistrictID,
		&user.AddressVillageID,
		&user.AddressStreet,
		&user.EducationLevel,
		&user.EducationInstitutionID,
		&user.MustChangePassword,
		&user.ConfirmedAt,
		&user.Occupation,
		&user.WorkInstitutionName,
		&user.ConfirmedAt,
		&user.UpdatedAt,
		&user.AuthKey,
		&user.Institution_name,
		&user.Province_name,
		&user.District_name,
		&user.Subdistrict_name,
		&user.Village_name,
	)
	if err != nil {
		return user, err
	}

	return user, nil
}
func (r *repository) Updateuserrepo(user User) (User, error) {

	if user.Password != "" {
		var sqlStmt string = "UPDATE user SET password_hash=? WHERE id=? "

		_, err := r.db.Exec(sqlStmt, user.Password, user.ID)

		if err != nil {
			return user, err
		}
	} 
	// else {
	// 	var sqlStmt string = "UPDATE user SET name=?, email=?, event_is_organizer=? WHERE id=? "

	// 	_, err := r.db.Exec(sqlStmt, user.Name.String, user.Email, user.Role, user.ID)

	// 	if err != nil {
	// 		return user, err
	// 	}
	// }

	return user, nil
}

func (r *repository) FetchAllUser() ([]User, error) {
	var sqlStmt string = "SELECT id, name, username, email, phone FROM user ORDER BY id DESC"

	rows, err := r.db.Query(sqlStmt)
	if err != nil {
		return nil, err
	}

	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Phone)
		if err != nil {
			return nil, err
		}


		users = append(users, user)
	}

	return users, nil
}

func (r *repository) Delete(ID int) (User, error) {
	var user User

	// var sqlStmt = "SELECT * FROM user WHERE id=?"

	// row := r.db.QueryRow(sqlStmt, ID)
	// err := row.Scan(
	// 	&user.ID,
	// 	&user.Username,
	// 	&user.Password,
	// 	&user.Email,
	// 	&user.Role,
	// )

	// if err != nil {
	// 	return user, err
	// }

	sqlStmt := "DELETE FROM user WHERE id=?"

	_, err := r.db.Exec(sqlStmt, ID)
	if err != nil {
		return user, err
	}

	return user, nil
}
