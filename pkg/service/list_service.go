package service

import (
	"restapi-app"
	"restapi-app/pkg/repository"
)

type TodoListService struct {
	repos repository.TodoList
}

func NewTodoListService(repos repository.TodoList) *TodoListService {
	return &TodoListService{repos}
}

func (s *TodoListService) Create(userId int, tdlist restapi.TodoList) (int, error) {
	return s.repos.Create(userId, tdlist)
}

func (s *TodoListService) GetAllLists(userid int) ([]restapi.TodoList, error) {
	return s.repos.GetAllLists(userid)
}

func (s *TodoListService) GetListById(userid, listid int) (restapi.TodoList, error) {
	return s.repos.GetListById(userid, listid)
}

func (s *TodoListService) Delete(userid, listid int) error {
	return s.repos.Delete(userid, listid)
}

func (s *TodoListService) UpdateList(userid, listid int, input restapi.UpdateListInput) error {
	err := input.Validate()
	if err != nil {
		return err
	}
	return s.repos.UpdateList(userid, listid, input)
}
