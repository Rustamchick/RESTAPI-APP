package service

import (
	"restapi-app"
	"restapi-app/pkg/repository"
)

type Authorization interface {
	CreateUser(user restapi.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface {
	Create(userId int, tdlist restapi.TodoList) (int, error)
	GetAllLists(userid int) ([]restapi.TodoList, error)
	GetListById(userid, listid int) (restapi.TodoList, error)
	Delete(userid, listid int) error
	UpdateList(userid, listid int, input restapi.UpdateListInput) error
}

type TodoItem interface {
	CreateItem(userid, listid int, tdItem restapi.TodoItem) (int, error)
	GetAllItems(userid, listid int) ([]restapi.TodoItem, error)
	GetItemByID(userid, itemid int) (restapi.TodoItem, error)
	UpdateItem(userid, itemid int, input restapi.UpdateItemInput) error
	DeleteItem(userid, itemid int) error
}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewTodoListService(repos.TodoList),
		TodoItem:      NewTodoItemService(repos.TodoItem, repos.TodoList),
	}
}
