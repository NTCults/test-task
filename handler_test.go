package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

var (
	repo    = initRepo()
	queue   = make(chan TaskRequest, 10)
	service = newService(&repo, queue)
)

func TestCreateHandler(t *testing.T) {
	var testTaskRequest = TaskRequest{
		Method: "GET",
		URL:    "http://www.test.test",
	}
	data, _ := json.Marshal(testTaskRequest)

	req, err := http.NewRequest("POST", "/task/", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	service.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	tdata := <-queue
	var testTask = CompletedTask{
		ID: tdata.ID,
	}
	repo.Set(testTask)

	expected := fmt.Sprintf(`{"taskID":"%s"}`, tdata.ID)
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetHandler(t *testing.T) {
	var testCompletedTask = CompletedTask{
		ID: "testID",
	}
	repo.Set(testCompletedTask)

	req, err := http.NewRequest("GET", "/task/testID", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	service.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := fmt.Sprintf(`{"id":"%s"}`, testCompletedTask.ID)
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestDeleteHandler(t *testing.T) {
	var testCompletedTask = CompletedTask{
		ID: "testID",
	}
	repo.Set(testCompletedTask)

	req, err := http.NewRequest("DELETE", "/task/testID", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	service.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusNoContent)
	}
}

func TestListHandler(t *testing.T) {
	for k := range repo.data {
		delete(repo.data, k)
	}
	numTasks := 20
	for i := 0; i < numTasks; i++ {

		taskID := strconv.Itoa(i)
		var testCompletedTask = CompletedTask{
			ID:   taskID,
			Date: time.Now().Add(time.Duration(i) * time.Hour),
		}
		repo.Set(testCompletedTask)
	}

	req, err := http.NewRequest("GET", "/task/?page=1&size=5", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	service.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `[{"id":"0"},{"id":"1"},{"id":"2"},{"id":"3"},{"id":"4"}]`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
