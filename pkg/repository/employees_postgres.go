package repository

import (
	"fmt"
	"strings"

	"github.com/Manifoldz/EmployeesRESTAPI/internal/entities"
	"github.com/jmoiron/sqlx"
)

type EmployeesPostgres struct {
	db *sqlx.DB
}

type EmployeeQueryResult struct {
	Id              int    `db:"id"`
	Name            string `db:"name"`
	Surname         string `db:"surname"`
	Phone           string `db:"phone"`
	CompanyId       int    `db:"company_id"`
	PassportType    string `db:"passport_type"`
	PassportNumber  string `db:"passport_number"`
	DepartmentName  string `db:"department_name"`
	DepartmentPhone string `db:"department_phone"`
}

func NewEmployeesPostgres(db *sqlx.DB) *EmployeesPostgres {
	return &EmployeesPostgres{db: db}
}

func (r *EmployeesPostgres) Create(input entities.EmployeeInputAndResponse) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	// добавим компанию, если она не существует
	checkArgsCompany := map[string]interface{}{"id": input.CompanyId}
	is_exist, _, err := CheckIfExists(tx, companiesTable, checkArgsCompany)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if !is_exist {
		createCompanyQuery := fmt.Sprintf("INSERT INTO %s (id) VALUES ($1)", companiesTable)
		if _, err = tx.Exec(createCompanyQuery, input.CompanyId); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// добавим департамент, если он не существует
	var departmentId int
	checkArgsDepartment := map[string]interface{}{"company_id": input.CompanyId, "name": input.Department.Name}
	is_exist, departmentId, err = CheckIfExists(tx, departmentsTable, checkArgsDepartment)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if !is_exist {
		createDepartmentQuery := fmt.Sprintf("INSERT INTO %s (company_id, name, phone) VALUES ($1, $2, $3) RETURNING id", departmentsTable)
		row1 := tx.QueryRow(createDepartmentQuery, input.CompanyId, input.Department.Name, input.Department.Phone)
		if err := row1.Scan(&departmentId); err != nil {
			tx.Rollback()
			return 0, err
		}
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

func (r *EmployeesPostgres) GetAll(companyId *int, departmentName *string, offset, limit int) ([]entities.EmployeeInputAndResponse, error) {
	var queryBuilder strings.Builder
	var args []interface{}
	var argCounter int = 1

	// добавление выбора
	queryBuilder.WriteString(`
	SELECT
		e.id,
		e.name,
		e.surname,
		e.phone,
		e.company_id,
		p.type AS passport_type,
		p.number AS passport_number,
		d.name AS department_name,
		d.phone AS department_phone`)

	// добавление объединения
	queryBuilder.WriteString(fmt.Sprintf(" FROM %s e JOIN %s d ON e.department_id = d.id", employeesTable, departmentsTable))
	queryBuilder.WriteString(fmt.Sprintf(" JOIN %s p ON e.id = p.employee_id", passportsTable))

	// добавление фильтрации по companyId, если он передан
	if companyId != nil {
		queryBuilder.WriteString(fmt.Sprintf(" WHERE e.company_id  = $%d", argCounter))
		args = append(args, *companyId)
		argCounter++
	}

	// добавление фильтрации по departmentName, если он передан и уже есть условие WHERE
	if departmentName != nil {
		if companyId != nil {
			queryBuilder.WriteString(" AND")
		} else {
			queryBuilder.WriteString(" WHERE")
		}
		queryBuilder.WriteString(fmt.Sprintf(" d.name = $%d", argCounter))
		args = append(args, *departmentName)
		argCounter++
	}

	// добавление группировки
	queryBuilder.WriteString(` ORDER BY e.id`)

	// добавление пагинации
	queryBuilder.WriteString(fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCounter, argCounter+1))
	args = append(args, limit, offset)

	// выполнение запроса
	finalQuery := queryBuilder.String()
	var queryResults []EmployeeQueryResult
	err := r.db.Select(&queryResults, finalQuery, args...)
	if err != nil {
		return nil, err
	}

	// сборка ответа
	var employees []entities.EmployeeInputAndResponse
	for _, qr := range queryResults {
		employee := entities.EmployeeInputAndResponse{
			Id:        qr.Id,
			Name:      qr.Name,
			Surname:   qr.Surname,
			Phone:     qr.Phone,
			CompanyId: qr.CompanyId,
			Passport: entities.Passport{
				Type:   qr.PassportType,
				Number: qr.PassportNumber,
			},
			Department: entities.Department{
				Name:  qr.DepartmentName,
				Phone: qr.DepartmentPhone,
			},
		}
		employees = append(employees, employee)
	}

	return employees, nil
}
