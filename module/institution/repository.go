package institution

import (
	"akuntansi/module/district"
	"akuntansi/module/province"
	"database/sql"
	"fmt"
)

type Repository interface {
	FetchAllInstitution() ([]Institution, error)
	FetchAllProvince() ([]province.Province, error)
	Save(institution Institution) (Institution, error)
	Delete(ID int) (Institution, error)
	FindByID(ID int) (Institution, error)
	Update(institution Institution) (Institution, error)
}

type repository struct {
	db *sql.DB
}

func (r *repository) Save(institution Institution) (Institution, error) {
	var sqlStmt string = "INSERT INTO institution (name, type, code, province_name) VALUES (?, ?, ?, ?);"

	_, err := r.db.Exec(sqlStmt, institution.Name.String, institution.Type.Int64, institution.Code.String, institution.Province_name.String)

	if err != nil {
		return institution, err
	}

	return institution, nil
}

func (r *repository) Update(institution Institution) (Institution, error) {

	var sqlStmt string = "UPDATE institution SET name=?, type=?, code=?, province_name=? WHERE id=?"

	_, err := r.db.Exec(sqlStmt, institution.Name.String, institution.Type.Int64, institution.Code.String, institution.Province_name.String, institution.ID)

	if err != nil {
		return institution, err
	}

	return institution, nil
}

func (r *repository) FindByID(ID int) (Institution, error) {
	var institution Institution

	var sqlStmt string = "SELECT id, name, type, code, province_name FROM institution WHERE id= ?"

	row := r.db.QueryRow(sqlStmt, ID)
	err := row.Scan(
		&institution.ID,
		&institution.Name,
		&institution.Type,
		&institution.Code,
		&institution.Province_name,
	)
	if err != nil {
		return institution, err
	}

	return institution, nil
}

func (r *repository) Delete(ID int) (Institution, error) {
	var institution Institution

	sqlStmt := "DELETE FROM institution WHERE id=?"

	_, err := r.db.Exec(sqlStmt, ID)
	if err != nil {
		return institution, err
	}

	return institution, nil
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db}
}

func (r *repository) FetchAllProvince() ([]province.Province, error) {
	var sqlStmt string = "SELECT id, name FROM province"

	rows, err := r.db.Query(sqlStmt)
	if err != nil {
		return nil, err
	}

	var provinces []province.Province
	for rows.Next() {
		var province province.Province
		err = rows.Scan(
			&province.ID,
			&province.Name,
		)
		if err != nil {
			return nil, err
		}
		province.SelectDistrict = func(province_id int) string {
			var sqlStmt string = "SELECT id, name FROM district WHERE province_id=?"

			rows, err := r.db.Query(sqlStmt, province_id)
			if err != nil {
				fmt.Println(err)
				return ""
			}

			for rows.Next() {
				var district district.District
				err = rows.Scan(
					&district.ID,
					&district.Name,
				)
				if err != nil {
					fmt.Println(err)
					return ""
				}
				fmt.Println(district.ID)
			}

			return ""
		}

		provinces = append(provinces, province)
	}

	return provinces, nil
}

// func FetchAllDistrict(province_id string) ([]district.District) {
// 	var r *repository
// 	var sqlStmt string = "SELECT id, name FROM district WHERE province_id=?"

// 	rows, err := r.db.Query(sqlStmt, province_id)
// 	if err != nil {
// 		return nil
// 	}

// 	var districts []district.District
// 	for rows.Next() {
// 		var district district.District
// 		err = rows.Scan(
// 			&district.ID,
// 			&district.Name,
// 		)
// 		if err != nil {
// 			return nil
// 		}

// 		districts = append(districts, district)
// 	}

// 	return districts
// }

func (r *repository) FetchAllInstitution() ([]Institution, error) {
	var sqlStmt string = "SELECT id, name, type, code, province_name FROM institution"

	rows, err := r.db.Query(sqlStmt)
	if err != nil {
		return nil, err
	}

	var institutions []Institution
	for rows.Next() {
		var institution Institution
		err = rows.Scan(
			&institution.ID,
			&institution.Name,
			&institution.Type,
			&institution.Code,
			&institution.Province_name,
		)
		if err != nil {
			return nil, err
		}
		institutions = append(institutions, institution)
	}

	return institutions, nil
}
