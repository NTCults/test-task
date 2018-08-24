package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"test-task/utils"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

const (
	defaultPageSize = 10
	defaultPageNum  = 1
)

type taskService struct {
	repo      Repository
	mux       *mux.Router
	taskQueue chan<- TaskRequest
}

func (s *taskService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func newService(repo Repository, queue chan<- TaskRequest) http.Handler {
	mux := mux.NewRouter()
	service := &taskService{repo, mux, queue}

	mux.HandleFunc("/task/", service.createTask).Methods("POST")
	mux.HandleFunc("/task/", service.listTask).Methods("GET")
	mux.HandleFunc("/task/{ID}", service.getTask).Methods("GET")
	mux.HandleFunc("/task/{ID}", service.deleteTask).Methods("DELETE")

	return service
}

func (s *taskService) createTask(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var task TaskRequest
	err := decoder.Decode(&task)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := task.Validate(); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	taskUUID := uuid.Must(uuid.NewV4())
	task.ID = taskUUID.String()

	s.taskQueue <- task
	log.Printf("Task %s queued", task.ID)

	response := make(map[string]string)
	response["taskID"] = task.ID
	utils.ResponseJSON(w, http.StatusCreated, response)
}

func (s *taskService) getTask(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["ID"]

	result, err := s.repo.Get(id)
	if err != nil {
		utils.ResponseError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.ResponseJSON(w, http.StatusOK, result)
}

func (s *taskService) deleteTask(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["ID"]

	err := s.repo.Delete(id)
	if err != nil {
		utils.ResponseError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.ResponseJSON(w, http.StatusNoContent, nil)
}

func (s *taskService) listTask(w http.ResponseWriter, req *http.Request) {
	page, err := strconv.Atoi(req.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = defaultPageNum
	}

	size, err := strconv.Atoi(req.URL.Query().Get("size"))
	if err != nil || size < 1 {
		size = defaultPageSize
	}

	sortedTasks := s.repo.List()
	utils.ResponseJSON(w, http.StatusOK, paginateTasks(sortedTasks, page-1, size))
}

func paginateTasks(tasks []CompletedTask, skip int, size int) []CompletedTask {
	lenTasks := len(tasks)
	if skip > lenTasks {
		skip = len(tasks)
	}
	end := skip + size
	if end > lenTasks {
		end = lenTasks
	}
	return tasks[skip:end]
}
