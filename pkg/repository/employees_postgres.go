package repository

import (
	"fmt"

	"github.com/Manifoldz/EmployeesRESTAPI/internal/entities"
	"github.com/jmoiron/sqlx"
)

type EmployeesPostgres struct {
	db *sqlx.DB
}

func NewEmployeesPostgres(db *sqlx.DB) *EmployeesPostgres {
	return &EmployeesPostgres{db: db}
}

func (r *EmployeesPostgres) Create(input entities.CreateEmployeeInput) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	// добавим компанию
	createCompanyQuery := fmt.Sprintf("INSERT INTO %s (id) VALUES ($1)", companiesTable)
	if _, err = tx.Exec(createCompanyQuery, input.CompanyId); err != nil {
		tx.Rollback()
		return 0, err
	}

	// добавим департмент
	var departmentId int
	createDepartmentQuery := fmt.Sprintf("INSERT INTO %s (company_id, name, phone) VALUES ($1, $2, $3) RETURNING id", departmentsTable)
	row1 := tx.QueryRow(createDepartmentQuery, input.CompanyId, input.Department.Name, input.Department.Phone)
	if err := row1.Scan(&departmentId); err != nil {
		tx.Rollback()
		return 0, err
	}

	// добавим сотрудника
	var employeeId int
	createEmployeeQuery := fmt.Sprintf("INSERT INTO %s (name, surname, phone, company_id, department_id) VALUES ($1, $2, $3, $4, $5) RETURNING id", employeesTable)
	row2 := tx.QueryRow(createEmployeeQuery, input.Name, input.Surname, input.Phone, input.CompanyId, departmentId)
	if err := row2.Scan(&employeeId); err != nil {
		tx.Rollback()
		return 0, err
	}

	// добавим паспорт
	createPassportQuery := fmt.Sprintf("INSERT INTO %s (employee_id, type, number) VALUES ($1, $2, $3)", passportsTable)
	if _, err = tx.Exec(createPassportQuery, employeeId, input.Passport.Type, input.Passport.Number); err != nil {
		tx.Rollback()
		return 0, err
	}

	return employeeId, tx.Commit()
}
