package repository

import (
	"errors"
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
	if err := CreateCompany(tx, input.CompanyId); err != nil {
		tx.Rollback()
		return 0, err
	}

	// добавим департамент, если он не существует
	var departmentId int
	if err := CreateDepartment(tx, input.CompanyId, input.Department.Name, input.Department.Phone, &departmentId); err != nil {
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

func (r *EmployeesPostgres) UpdateById(id int, input entities.UpdateEmployeeInput) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	if is_exists, err := isResourceExists(tx, employeesTable, id); err != nil {
		tx.Rollback()
		return err
	} else if !is_exists {
		tx.Rollback()
		return errors.New("employee not found")
	}

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	// если id передан, то обновляем его вместе с увязкой паспорта
	if input.Id != nil {
		queryCopyOldEmployee := fmt.Sprintf("INSERT INTO %s (id, name, surname, phone, company_id, department_id) SELECT $1, name, surname, phone, company_id, department_id FROM %s WHERE id = $2;", employeesTable, employeesTable)
		if _, err := tx.Exec(queryCopyOldEmployee, *input.Id, id); err != nil {
			tx.Rollback()
			return err
		}
		queryUpdatePassport := fmt.Sprintf("UPDATE %s SET employee_id = $1 WHERE employee_id = $2;", passportsTable)
		if _, err := tx.Exec(queryUpdatePassport, *input.Id, id); err != nil {
			tx.Rollback()
			return err
		}
		queryDeleteOldEmployee := fmt.Sprintf("DELETE FROM %s WHERE id = $1;", employeesTable)
		if _, err := tx.Exec(queryDeleteOldEmployee, id); err != nil {
			tx.Rollback()
			return err
		}
		id = *input.Id
	}

	if input.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *input.Name)
		argId++
	}

	if input.Surname != nil {
		setValues = append(setValues, fmt.Sprintf("surname=$%d", argId))
		args = append(args, *input.Surname)
		argId++
	}

	if input.Phone != nil {
		setValues = append(setValues, fmt.Sprintf("phone=$%d", argId))
		args = append(args, *input.Phone)
		argId++
	}

	if input.CompanyId != nil {
		// добавим компанию, если она не существует
		if err := CreateCompany(tx, *input.CompanyId); err != nil {
			tx.Rollback()
			return err
		}
		setValues = append(setValues, fmt.Sprintf("company_id=$%d", argId))
		args = append(args, *input.CompanyId)
		argId++
	}

	if input.Department != nil && input.Department.Name != nil {
		if input.CompanyId == nil {
			var companyId int
			getCompanyIdQuery := fmt.Sprintf("SELECT company_id  FROM  %s WHERE id  =  $1;", employeesTable)
			row := tx.QueryRow(getCompanyIdQuery, id)
			err = row.Scan(&companyId)
			if err != nil {
				tx.Rollback()
				return err
			}
			input.CompanyId = &companyId
		}
		// добавим департамент, если он не существует, а если существует получим id департамента
		var departmentId int
		var department_phone string
		if input.Department != nil && input.Department.Phone != nil {
			department_phone = *input.Department.Phone
		}

		if err := CreateDepartment(tx, *input.CompanyId, *input.Department.Name, department_phone, &departmentId); err != nil {
			tx.Rollback()
			return err
		}
		setValues = append(setValues, fmt.Sprintf("department_id=$%d", argId))
		args = append(args, departmentId)
		argId++
	}

	if input.Department != nil && input.Department.Phone != nil && input.Department.Name == nil {
		var departmentId int
		getDepartmentIdQuery := fmt.Sprintf("SELECT department_id  FROM  %s WHERE id  =  $1;", employeesTable)
		row := tx.QueryRow(getDepartmentIdQuery, id)
		err = row.Scan(&departmentId)
		if err != nil {
			tx.Rollback()
			return err
		}
		queryUpdateDepartPhone := fmt.Sprintf("UPDATE %s SET phone = $1 WHERE id = $2;", departmentsTable)
		if _, err := tx.Exec(queryUpdateDepartPhone, *input.Department.Phone, departmentId); err != nil {
			tx.Rollback()
			return err
		}
	}

	if input.Passport != nil && input.Passport.Number != nil {
		queryUpdatePassport := fmt.Sprintf("UPDATE %s SET number = $1 WHERE employee_id = $2;", passportsTable)
		if _, err := tx.Exec(queryUpdatePassport, *input.Passport.Number, id); err != nil {
			tx.Rollback()
			return err
		}
	}

	if input.Passport != nil && input.Passport.Type != nil {
		queryUpdatePassport := fmt.Sprintf("UPDATE %s SET type = $1 WHERE employee_id = $2;", passportsTable)
		if _, err := tx.Exec(queryUpdatePassport, *input.Passport.Type, id); err != nil {
			tx.Rollback()
			return err
		}
	}

	if argId != 1 {
		setQuery := strings.Join(setValues, ", ")

		query := fmt.Sprintf("UPDATE %s t1 SET %s WHERE t1.id  = $%d;", employeesTable, setQuery, argId)
		args = append(args, id)
		if _, err := tx.Exec(query, args...); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *EmployeesPostgres) DeleteById(id int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	if is_exists, err := isResourceExists(tx, employeesTable, id); err != nil {
		tx.Rollback()
		return err
	} else if !is_exists {
		tx.Rollback()
		return errors.New("employee not found")
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1;", employeesTable)
	_, err = tx.Exec(query, id)
	if err != nil {
		tx.Rollback()
	}
	return err
}
