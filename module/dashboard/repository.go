package dashboard

import (
	"database/sql"
)

type DashboardRepository interface {
	CountUser() (int, error)
	CountInstitution() (int, error)
	CountRequest() (int, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db}
}

func (r *repository) CountUser() (int, error) {
	var count int
	var sqlStmt string = "SELECT COUNT(*) FROM user"

	rows, err := r.db.Query(sqlStmt)

	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return count, err
		}
	}

	if err != nil {
		return count, err
	}

	return count, nil
}

func (r *repository) CountInstitution() (int, error) {
	var count int
	var sqlStmt string = "SELECT COUNT(*) FROM institution"

	rows, err := r.db.Query(sqlStmt)

	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return count, err
		}
	}

	if err != nil {
		return count, err
	}

	return count, nil
}

func (r *repository) CountRequest() (int, error) {
	var count int
	var sqlStmt string = "SELECT COUNT(*) FROM request"

	rows, err := r.db.Query(sqlStmt)

	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return count, err
		}
	}

	if err != nil {
		return count, err
	}

	return count, nil
}
