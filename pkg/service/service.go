package service

import (
	"github.com/Manifoldz/EmployeesRESTAPI/internal/entities"
	"github.com/Manifoldz/EmployeesRESTAPI/pkg/repository"
)

type ToDoList interface {
	Create(list entities.ToDoList) (int, error)
	GetAll() ([]entities.ToDoList, error)
	GetById(id int) (entities.ToDoList, error)
	DeleteById(id int) error
	UpdateById(id int, input entities.UpdateListInput) error
}
type ToDoItem interface {
	Create(listId int, item entities.ToDoItem) (int, error)
	GetAll(listId int, status *bool, offset, limit int) ([]entities.ToDoItem, error)
	GetById(itemId int) (entities.ToDoItem, error)
	DeleteById(itemId int) error
	UpdateById(itemId int, input entities.UpdateItemInput) error
}

type Service struct {
	ToDoList
	ToDoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		ToDoList: NewTodoListService(repos.ToDoList),
		ToDoItem: NewTodoItemService(repos.ToDoItem, repos.ToDoList),
	}
}
