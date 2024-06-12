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
