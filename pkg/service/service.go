package service

import (
	"github.com/Manifoldz/EmployeesRESTAPI/internal/entities"
	"github.com/Manifoldz/EmployeesRESTAPI/pkg/repository"
)

type Employees interface {
	Create(input entities.EmployeeInputAndResponse) (int, error)
	GetAll(companyId *int, departmentName *string, offset, limit int) ([]entities.EmployeeInputAndResponse, error)
	UpdateById(id int, input entities.UpdateEmployeeInput) error
	DeleteById(id int) error
}

type Service struct {
	Employees
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Employees: NewEmployeesService(repos.Employees),
	}
}
