package service

import (
	"github.com/Manifoldz/EmployeesRESTAPI/internal/entities"
	"github.com/Manifoldz/EmployeesRESTAPI/pkg/repository"
)

type TodoListService struct {
	repo repository.ToDoList
}

func NewTodoListService(repo repository.ToDoList) *TodoListService {
	return &TodoListService{repo: repo}
}

func (s *TodoListService) Create(list entities.ToDoList) (int, error) {
	return s.repo.Create(list)
}

func (s *TodoListService) GetAll() ([]entities.ToDoList, error) {
	return s.repo.GetAll()
}

func (s *TodoListService) GetById(id int) (entities.ToDoList, error) {
	return s.repo.GetById(id)
}

func (s *TodoListService) DeleteById(id int) error {
	return s.repo.DeleteById(id)
}

func (s *TodoListService) UpdateById(id int, input entities.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.UpdateById(id, input)
}
