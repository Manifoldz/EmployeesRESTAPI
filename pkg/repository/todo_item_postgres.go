package repository

import (
	"fmt"
	"strings"

	"github.com/Manifoldz/EmployeesRESTAPI/internal/entities"
	"github.com/jmoiron/sqlx"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Create(listId int, item entities.ToDoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var itemId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id;", toDoItemsTable)

	row := tx.QueryRow(createItemQuery, item.Title, item.Description)
	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createListItemsQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES  ($1, $2);", listsItemsTable)
	_, err = tx.Exec(createListItemsQuery, listId, itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

func (r *TodoItemPostgres) GetAll(listId int, status *bool, offset, limit int) ([]entities.ToDoItem, error) {
	var items []entities.ToDoItem
	var query string

	if status == nil {
		query = fmt.Sprintf("SELECT t1.id, t1.title, t1.description, t1.done FROM %s t1 INNER JOIN %s t2 ON t2.item_id  = t1.id WHERE t2.list_id  =  $1 ORDER BY t1.id LIMIT $2 OFFSET $3;", toDoItemsTable, listsItemsTable)
	} else {
		query = fmt.Sprintf("SELECT t1.id, t1.title, t1.description, t1.done FROM %s t1 INNER JOIN %s t2 ON t2.item_id  = t1.id WHERE t2.list_id  =  $1 AND t1.done =  $2  ORDER BY t1.id LIMIT $3 OFFSET $4;", toDoItemsTable, listsItemsTable)
	}

	var err error
	if status == nil {
		err = r.db.Select(&items, query, listId, limit, offset)
	} else {
		err = r.db.Select(&items, query, listId, *status, limit, offset)
	}

	if err != nil {
		return nil, err
	}

	return items, nil

}

func (r *TodoItemPostgres) GetById(itemId int) (entities.ToDoItem, error) {
	var item entities.ToDoItem
	query := fmt.Sprintf("SELECT t1.id, t1.title, t1.description, t1.done FROM %s t1 INNER JOIN %s t2 ON t2.item_id  = t1.id WHERE t1.id  =  $1;", toDoItemsTable, listsItemsTable)

	if err := r.db.Get(&item, query, itemId); err != nil {
		return item, err
	}

	return item, nil
}

func (r *TodoItemPostgres) DeleteById(itemId int) error {
	query := fmt.Sprintf("DELETE FROM %s t1 USING %s t2 WHERE t1.id = t2.item_id AND t1.id = $1;", toDoItemsTable, listsItemsTable)
	_, err := r.db.Exec(query, itemId)

	return err
}

func (r *TodoItemPostgres) UpdateById(id int, input entities.UpdateItemInput) error {
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

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s t1 SET %s FROM %s t2 WHERE t1.id  = t2.item_id AND t1.id  = $%d;", toDoItemsTable, setQuery, listsItemsTable, argId)
	args = append(args, id)
	_, err := r.db.Exec(query, args...)

	return err
}
