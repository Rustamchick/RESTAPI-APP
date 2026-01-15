package service

import (
	"errors"
	"restapi-app"
	"restapi-app/pkg/repository"
)

type TodoItemService struct {
	repos     repository.TodoItem
	listRepos repository.TodoList
}

func NewTodoItemService(repos repository.TodoItem, listRepos repository.TodoList) *TodoItemService {
	return &TodoItemService{repos, listRepos}
}

func (s *TodoItemService) CreateItem(userid, listid int, tdItem restapi.TodoItem) (int, error) {
	_, err := s.listRepos.GetListById(userid, listid)
	if err != nil {
		return 0, errors.New("there is no list for this user. Item_service.CreateItem")
	}

	return s.repos.CreateItem(listid, tdItem)
}

func (s *TodoItemService) GetAllItems(userid, listid int) ([]restapi.TodoItem, error) {
	_, err := s.listRepos.GetListById(userid, listid)
	if err != nil {
		return nil, errors.New("there is no list for this user. Item_service.GetAllItems")
	}

	return s.repos.GetAllItems(listid)
}

func (s *TodoItemService) GetItemByID(userid, itemid int) (restapi.TodoItem, error) {

	return s.repos.GetItemByID(userid, itemid)
}

func (s *TodoItemService) UpdateItem(userid, itemid int, input restapi.UpdateItemInput) error {
	err := input.Validate()
	if err != nil {
		return err
	}

	return s.repos.UpdateItem(userid, itemid, input)
}

func (s *TodoItemService) DeleteItem(userid, itemid int) error {

	return s.repos.DeleteItem(userid, itemid)
}
