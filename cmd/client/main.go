package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	todov1 "todo/gen/todo/v1"
	"todo/gen/todo/v1/todov1connect"

	"connectrpc.com/connect"
)

func main() {
	client := todov1connect.NewTodoServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
	)

	if len(os.Args) >= 2 {
		value := os.Args[1]

		if value == "show" {
			ShowGetAllTasks(client)
		} else if value == "new" && len(os.Args) >= 3 {
			joinTaskName := strings.Join(os.Args[2:], " ")
			CreateNewTask(client, joinTaskName)
		} else if value == "delete" && len(os.Args) >= 3 {
			DeleteTask(client, os.Args[2])
		} else if value == "update" && len(os.Args) >= 4 {
			status := os.Args[3]
			if isValidTodoStatus(status) {
				UpdateTask(client, os.Args[2], os.Args[3])
			} else {
				log.Println("update:", "invalid status")
				log.Println("\t", "need todo or done")
			}
		} else {
			log.Println("invalid parameter")
		}
	} else {
		log.Println("need parameters")
	}
}

func ShowGetAllTasks(client todov1connect.TodoServiceClient) {
	var req = &todov1.GetAllTasksRequest{}
	res, err := client.GetAllTasks(
		context.Background(),
		connect.NewRequest(req),
	)

	if err != nil {
		log.Println(err)
		return
	}

	for _, task := range res.Msg.Items {
		log.Println(task)
	}
}

func CreateNewTask(client todov1connect.TodoServiceClient, taskName string) {
	var newTask = &todov1.CreateTaskRequest{
		Name:   taskName,
		Status: todov1.Status_STATUS_TODO,
	}

	res, err := client.CreateTask(
		context.Background(),
		connect.NewRequest(newTask),
	)

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Add New Task {", res.Msg, "}")
}

func DeleteTask(client todov1connect.TodoServiceClient, taskId string) {
	delId, err := strconv.ParseUint(taskId, 10, 64)
	if err != nil {
		log.Printf("DeleteTask: failed to convert taskId to uint64: %v", err)
		return
	}

	var delTask = &todov1.DeleteTaskRequest{
		Id: delId,
	}

	res, err := client.DeleteTask(
		context.Background(),
		connect.NewRequest(delTask),
	)

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Delete Task {", res.Msg, "}")
}

func UpdateTask(client todov1connect.TodoServiceClient, taskId string, status string) {
	updateId, err := strconv.ParseUint(taskId, 10, 64)
	if err != nil {
		log.Printf("UpdateTask: failed to convert taskId to uint64: %v", err)
		return
	}
	var updateTask = &todov1.UpdateTaskStatusRequest{
		Id:     updateId,
		Status: todov1.Status_STATUS_UNKNOWN_UNSPECIFIED,
	}

	switch status {
	case "todo":
		updateTask.Status = todov1.Status_STATUS_TODO
	case "done":
		updateTask.Status = todov1.Status_STATUS_DONE
	default:
		updateTask.Status = todov1.Status_STATUS_UNKNOWN_UNSPECIFIED
	}

	res, err := client.UpdateTaskStatus(
		context.Background(),
		connect.NewRequest(updateTask),
	)

	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Update Task {", res.Msg, "}")
}

func isValidTodoStatus(s string) bool {
	validStatus := []string{"todo", "done"}
	isValid := false
	for _, status := range validStatus {
		if s == status {
			isValid = true
			break
		}
	}
	return isValid
}
