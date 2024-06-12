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

func (s *EmployeesService) Create(input entities.EmployeeInputAndResponse) (int, error) {
	return s.employeesRepo.Create(input)
}

func (s *EmployeesService) GetAll(companyId, departmentId *int, offset, limit int) ([]entities.EmployeeInputAndResponse, error) {
	return s.employeesRepo.GetAll(companyId, departmentId, offset, limit)
}
