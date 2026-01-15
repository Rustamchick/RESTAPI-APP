package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"restapi-app"
	"strings"

	"github.com/jmoiron/sqlx"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db}
}

func (p *TodoItemPostgres) CreateItem(listid int, tdItem restapi.TodoItem) (int, error) {
	tr, err := p.db.Begin() // start transaction
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tr.Rollback()
		}
	}()

	var itemid int
	queryTodoItems := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoItemsTable)
	row := tr.QueryRow(queryTodoItems, tdItem.Title, tdItem.Description)
	if err = row.Scan(&itemid); err != nil {
		// tr.Rollback()
		return 0, errors.New("error scanning list id")
	}

	queryListItems := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2)", listItemsTable)
	_, err = tr.Exec(queryListItems, listid, itemid)
	if err != nil {
		// tr.Rollback()
		return 0, errors.New("error listItem. item_postgres")
	}

	return itemid, tr.Commit()
}

func (p *TodoItemPostgres) GetAllItems(listid int) ([]restapi.TodoItem, error) {
	var lists []restapi.TodoItem

	queryTodoItems := fmt.Sprintf("SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li ON ti.id=li.item_id WHERE li.list_id=$1", todoItemsTable, listItemsTable)
	err := p.db.Select(&lists, queryTodoItems, listid)

	return lists, err
}

func (p *TodoItemPostgres) GetItemByID(userid, itemid int) (restapi.TodoItem, error) {
	var item restapi.TodoItem

	query := fmt.Sprintf(`SELECT ti.id, ti.title, ti.description, ti.done FROM %s ti INNER JOIN %s li ON ti.id=li.item_id 
																					 INNER JOIN %s ul ON ul.list_id=li.list_id WHERE ti.id=$1 AND ul.user_id=$2`, todoItemsTable, listItemsTable, usersListsTable)

	err := p.db.Get(&item, query, itemid, userid)

	return item, err
}

func (p *TodoItemPostgres) UpdateItem(userid, itemid int, input restapi.UpdateItemInput) error {
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

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done = $%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ",") // title = $1, description = $2, done = $3

	updateItemQuery := fmt.Sprintf("UPDATE %s ti SET %s FROM %s li, %s ul WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d AND ti.id = $%d", todoItemsTable, setQuery, listItemsTable, usersListsTable, argId, argId+1)
	args = append(args, userid, itemid)

	fmt.Printf("updateItemQuery: %s", updateItemQuery)
	fmt.Printf("args: %s", args)

	_, err := p.db.Exec(updateItemQuery, args...)

	if err == sql.ErrNoRows { // хаха, dont use QueryRow() for updating sql shit, postgres dont return anything, so it will be ErrNoRows everytime
		fmt.Print("\n!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!Error No Rows!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n")
		return err
	}

	return err
}

// func (p *TodoItemPostgres) DeleteItem(userid, itemid int) error {
// 	tr, err := p.db.Begin()
// 	if err != nil {
// 		return err
// 	}

// 	defer func() {
// 		if err != nil {
// 			tr.Rollback()
// 		}
// 	}()

// 	deleteListItemQuery := fmt.Sprintf("DELETE FROM %s li USING %s ul WHERE ul.list_id=li.list_id AND ul.user_id=$1 and li.item_id=$2", listItemsTable, usersListsTable)
// 	_, err = tr.Exec(deleteListItemQuery, userid, itemid)
// 	if err != nil {
// 		return err // tr.Rollback()
// 	}

// 	deleteItemQuery := fmt.Sprintf("DELETE FROM %s WHERE id=$1", todoItemsTable)
// 	_, err = tr.Exec(deleteItemQuery, itemid)
// 	if err != nil {
// 		return err // tr.Rollback()
// 	}

// 	return tr.Commit()
// }

func (p *TodoItemPostgres) DeleteItem(userid, itemid int) error {

	query := fmt.Sprintf("DELETE FROM %s ti USING %s li, %s ul WHERE ti.id=li.item_id AND li.list_id=ul.list_id AND ul.user_id=$1 AND ti.id=$2", todoItemsTable, listItemsTable, usersListsTable)

	_, err := p.db.Exec(query, userid, itemid)

	return err
}
