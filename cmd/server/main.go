package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	todov1 "todo/gen/todo/v1"        // generated by protoc-gen-go
	"todo/gen/todo/v1/todov1connect" // generated by protoc-gen-connect-go
)

type TodoServer struct {
	todos  sync.Map //スレッドセーフなMapらしい syncにいろいろある
	nextId int
}

func (s *TodoServer) CreateTask(
	ctx context.Context,
	req *connect.Request[todov1.CreateTaskRequest],
) (*connect.Response[todov1.CreateTaskResponse], error) {
	s.nextId++
	todoId := s.nextId

	createTodo := &todov1.TodoItem{
		Id:     uint64(todoId),
		Name:   req.Msg.Name,
		Status: req.Msg.Status,
	}

	s.todos.Store(createTodo.Id, createTodo)
	log.Println("CreateTodo")

	log.Println("Request headers: ", req.Header())

	res := connect.NewResponse(&todov1.CreateTaskResponse{
		Id:     createTodo.Id,
		Name:   createTodo.Name,
		Status: createTodo.Status,
	})

	//一覧をログに表示
	LogTodoList(s)

	res.Header().Set("CreateTask-Version", "v1")
	return res, nil
}

func (s *TodoServer) UpdateTaskStatus(
	ctx context.Context,
	req *connect.Request[todov1.UpdateTaskStatusRequest],
) (*connect.Response[todov1.UpdateTaskStatusResponse], error) {
	updateId := req.Msg.Id
	log.Println("Request headers: ", req.Header())

	v, ok := s.todos.Load(updateId)
	if !ok {
		err := fmt.Errorf("todo not found: %d", updateId)
		fmt.Printf("Todo not found: %d\n", updateId)
		return nil, err
	}

	todo := v.(*todov1.TodoItem)

	updateTodo := &todov1.TodoItem{
		Id:     updateId,
		Name:   todo.Name,
		Status: req.Msg.Status,
	}

	// UPDATE TODO
	s.todos.Store(updateTodo.Id, updateTodo)
	log.Println("TODOを更新")

	res := connect.NewResponse(&todov1.UpdateTaskStatusResponse{
		Id:     updateTodo.Id,
		Status: updateTodo.Status,
	})

	//Todo一覧をログに表示
	LogTodoList(s)

	res.Header().Set("UpdateTaskStatus-Version", "v1")
	return res, nil
}

func isEmpty(m *sync.Map) bool {
	isEmpty := true
	m.Range(func(key, value interface{}) bool {
		isEmpty = false
		return false
	})
	return isEmpty
}

func (s *TodoServer) DeleteTask(
	ctx context.Context,
	req *connect.Request[todov1.DeleteTaskRequest],
) (*connect.Response[todov1.DeleteTaskResponse], error) {
	delId := req.Msg.Id

	_, ok := s.todos.Load(delId)
	if !ok {
		err := fmt.Errorf("todo not found: %d", delId)
		fmt.Printf("Todo not found: %d\n", delId)
		return nil, err
	}

	// DELETE TODO
	s.todos.Delete(delId)
	// TODO: Implement this method.
	log.Println("Request headers: ", req.Header())
	res := connect.NewResponse(&todov1.DeleteTaskResponse{
		Id: delId,
	})
	LogTodoList(s)
	res.Header().Set("DeleteTask-Version", "v1")
	return res, nil
}

func LogTodoList(s *TodoServer) {
	log.Println("TODO一覧")
	if isEmpty(&s.todos) {
		log.Println("todoはありません")
	} else {
		s.todos.Range(func(key, value interface{}) bool {
			log.Printf("key is %d, value is %v", key, value)
			return true
		})
	}
}

func (s *TodoServer) GetAllTasks(
	ctx context.Context,
	req *connect.Request[todov1.GetAllTasksRequest],
) (*connect.Response[todov1.GetAllTasksResponse], error) {
	var items []*todov1.TodoItem
	s.todos.Range(func(key, value interface{}) bool {
		if item, ok := value.(*todov1.TodoItem); ok {
			items = append(items, item)
		}
		return true
	})

	log.Println("Request headers: ", req.Header())

	res := connect.NewResponse(&todov1.GetAllTasksResponse{
		Items: items,
	})
	res.Header().Set("GetAllTasks-Version", "v1")
	return res, nil
}

func main() {
	todoServer := &TodoServer{}
	mux := http.NewServeMux()
	path, handler := todov1connect.NewTodoServiceHandler(todoServer)

	log.Println("path: ", path, "handler: ", handler)
	mux.Handle(path, handler)

	http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
