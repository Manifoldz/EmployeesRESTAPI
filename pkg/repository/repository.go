package repository

import (
	"github.com/Manifoldz/EmployeesRESTAPI/internal/entities"
	"github.com/jmoiron/sqlx"
)

type Employees interface {
	Create(input entities.EmployeeInputAndResponse) (int, error)
	GetAll(companyId *int, departmentName *string, offset, limit int) ([]entities.EmployeeInputAndResponse, error)
	UpdateById(id int, input entities.UpdateEmployeeInput) error
	// DeleteById(id int) error
	// UpdateById(id int, input entities.UpdateListInput) error
}

type Repository struct {
	Employees
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Employees: NewEmployeesPostgres(db),
	}
}
