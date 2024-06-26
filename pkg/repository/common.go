package repository

import (
	"database/sql"
	"fmt"
	"strings"
)

func CheckIfExists(tx *sql.Tx, table string, queryArgs map[string]interface{}) (bool, int, error) {
	placeholders := make([]string, 0, len(queryArgs))
	args := make([]interface{}, 0, len(queryArgs))
	i := 1
	for key, value := range queryArgs {
		placeholders = append(placeholders, fmt.Sprintf("%s = $%d", key, i))
		args = append(args, value)
		i++
	}
	queryString := fmt.Sprintf("SELECT id FROM %s WHERE %s LIMIT 1", table, strings.Join(placeholders, " AND "))

	var id int
	err := tx.QueryRow(queryString, args...).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, 0, nil // Запись не найдена
		}
		return false, 0, err // Ошибка запроса
	}
	return true, id, nil // Запись найдена
}

func CreateCompany(tx *sql.Tx, id int) error {
	checkArgsCompany := map[string]interface{}{"id": id}
	is_exist, _, err := CheckIfExists(tx, companiesTable, checkArgsCompany)
	if err != nil {
		return err
	}
	if !is_exist {
		createCompanyQuery := fmt.Sprintf("INSERT INTO %s (id) VALUES ($1)", companiesTable)
		if _, err = tx.Exec(createCompanyQuery, id); err != nil {
			return err
		}
	}
	return nil
}

func CreateDepartment(tx *sql.Tx, companyId int, departmentName, departmentPhone string, departmentId *int) error {
	checkArgsDepartment := map[string]interface{}{"company_id": companyId, "name": departmentName}
	is_exist, id, err := CheckIfExists(tx, departmentsTable, checkArgsDepartment)
	if err != nil {
		tx.Rollback()
		return err
	}
	if !is_exist {
		createDepartmentQuery := fmt.Sprintf("INSERT INTO %s (company_id, name, phone) VALUES ($1, $2, $3) RETURNING id", departmentsTable)
		row1 := tx.QueryRow(createDepartmentQuery, companyId, departmentName, departmentPhone)
		if err := row1.Scan(&id); err != nil {
			tx.Rollback()
			return err
		}
	}

	*departmentId = id
	return nil
}

func isResourceExists(tx *sql.Tx, table string, id int) (bool, error) {
	queryCheck := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE id = $1", table)
	var count int
	if err := tx.QueryRow(queryCheck, id).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}
