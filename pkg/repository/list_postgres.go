package repository

import (
	"errors"
	"fmt"
	"restapi-app"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db}
}

func (p *TodoListPostgres) Create(userId int, tdlist restapi.TodoList) (int, error) {
	tr, err := p.db.Begin() // creating a transaction
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tr.Rollback()
		}
	}()

	var listId int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListsTable)
	row := tr.QueryRow(createListQuery, tdlist.Title, tdlist.Description)

	if err = row.Scan(&listId); err != nil {
		return 0, errors.New("error scanning list id")
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	_, err = tr.Exec(createUsersListQuery, userId, listId)
	if err != nil {
		return 0, errors.New("error while executing")
	}

	return listId, tr.Commit()
}

func (p *TodoListPostgres) GetAllLists(userid int) ([]restapi.TodoList, error) { // посмотрел
	var tdlists []restapi.TodoList

	query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul ON ul.list_id = tl.id WHERE ul.user_id=$1`, todoListsTable, usersListsTable)
	err := p.db.Select(&tdlists, query, userid)

	return tdlists, err
}

func (p *TodoListPostgres) GetListById(userid, listid int) (restapi.TodoList, error) { // надо проверить что он выдает
	var tdlist restapi.TodoList

	query := fmt.Sprintf(`SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul ON ul.list_id = tl.id WHERE ul.list_id=$1 AND  ul.user_id=$2`, todoListsTable, usersListsTable)
	err := p.db.Get(&tdlist, query, listid, userid)

	return tdlist, err
}

func (p *TodoListPostgres) Delete(userid, listid int) error {
	tr, err := p.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tr.Rollback()
		}
	}()

	deleteUsersListsQuery := fmt.Sprintf("DELETE FROM %s WHERE list_id=$1 and user_id=$2", usersListsTable)
	_, err = tr.Exec(deleteUsersListsQuery, listid, userid)
	if err != nil {
		return err // tr.Rollback()
	}

	deleteListsQuery := fmt.Sprintf("DELETE FROM %s WHERE id=$1", todoListsTable)
	_, err = tr.Exec(deleteListsQuery, listid)
	if err != nil {
		return err // tr.Rollback()
	}

	return tr.Commit()
}

// func (p *TodoListPostgres) Delete(userid, listid int) error {
// 	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = ul.list_id AND ul.user_id=$1 AND ul.list_id=$2", todoListsTable, usersListsTable)
// 	_, err := p.db.Exec(query, userid, listid)
// 	if err != nil { // под вопросом
// 		return err
// 	}
// 	return err
// }

func (p *TodoListPostgres) UpdateList(userid, listid int, input restapi.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title = $%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description = $%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	setQuery := strings.Join(setValues, ",") // title = $1, description = $2

	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.user_id = $%d AND ul.list_id = $%d", todoListsTable, setQuery, usersListsTable, argId, argId+1)
	args = append(args, userid, listid)

	logrus.Debugf("updateListQuery: %s /n", query)
	logrus.Debugf("args: %s /n", args)

	_, err := p.db.Exec(query, args...)

	return err
}
