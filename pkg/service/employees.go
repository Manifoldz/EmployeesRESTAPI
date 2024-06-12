package service

import (
	"github.com/Manifoldz/EmployeesRESTAPI/internal/entities"
	"github.com/Manifoldz/EmployeesRESTAPI/pkg/repository"
)

type EmployeesService struct {
	employeesRepo repository.Employees
}

func NewEmployeesService(employeesRepo repository.Employees) *EmployeesService {
	return &EmployeesService{
		employeesRepo: employeesRepo,
	}
}

func (s *EmployeesService) Create(input entities.CreateEmployeeInput) (int, error) {
	return s.employeesRepo.Create(input)
}
