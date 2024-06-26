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

func (s *EmployeesService) GetAll(companyId *int, departmentName *string, offset, limit int) ([]entities.EmployeeInputAndResponse, error) {
	return s.employeesRepo.GetAll(companyId, departmentName, offset, limit)
}

func (s *EmployeesService) UpdateById(id int, input entities.UpdateEmployeeInput) error {
	return s.employeesRepo.UpdateById(id, input)
}

func (s *EmployeesService) DeleteById(id int) error {
	return s.employeesRepo.DeleteById(id)
}
