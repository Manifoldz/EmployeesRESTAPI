package service

import (
	"github.com/Manifoldz/EmployeesRESTAPI/pkg/repository"
)

type Employees interface {
}

type Service struct {
	Employees
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Employees: NewEmployeesService(repos.Employees),
	}
}
