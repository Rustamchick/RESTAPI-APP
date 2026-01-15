package repository

import (
	"restapi-app"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user restapi.User) (int, error)
	GetUser(username, password string) (restapi.User, error)
}

type TodoList interface {
	Create(userId int, tdlist restapi.TodoList) (int, error)
	GetAllLists(userid int) ([]restapi.TodoList, error)
	GetListById(userid, listid int) (restapi.TodoList, error)
	Delete(userid, listid int) error
	UpdateList(userid, listid int, input restapi.UpdateListInput) error
}

type TodoItem interface {
	CreateItem(listid int, tdItem restapi.TodoItem) (int, error)
	GetAllItems(listid int) ([]restapi.TodoItem, error)
	GetItemByID(userid, itemid int) (restapi.TodoItem, error)
	UpdateItem(userid, itemid int, input restapi.UpdateItemInput) error
	DeleteItem(userid, itemid int) error
}

type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
	}
}
