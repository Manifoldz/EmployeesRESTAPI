package repository

import (
	"github.com/Manifoldz/EmployeesRESTAPI/internal/entities"
	"github.com/jmoiron/sqlx"
)

type Employees interface {
	Create(input entities.EmployeeInputAndResponse) (int, error)
	GetAll(companyId, departmentId *int, offset, limit int) ([]entities.EmployeeInputAndResponse, error)
	// GetById(id int) (entities.ToDoList, error)
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
