package repository

import (
	"fmt"
	"strings"

	"github.com/Manifoldz/EmployeesRESTAPI/internal/entities"
	"github.com/jmoiron/sqlx"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) Create(list entities.ToDoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", toDoListsTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TodoListPostgres) GetAll() ([]entities.ToDoList, error) {
	var list []entities.ToDoList
	query := fmt.Sprintf("SELECT * FROM %s", toDoListsTable)
	err := r.db.Select(&list, query)

	return list, err
}

func (r *TodoListPostgres) GetById(id int) (entities.ToDoList, error) {
	var list entities.ToDoList
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", toDoListsTable)
	err := r.db.Get(&list, query, id)

	return list, err
}

func (r *TodoListPostgres) DeleteById(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", toDoListsTable)
	_, err := r.db.Exec(query, id)

	return err
}

func (r *TodoListPostgres) UpdateById(id int, input entities.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id  =  $%d", toDoListsTable, setQuery, argId)
	args = append(args, id)
	_, err := r.db.Exec(query, args...)

	return err
}
