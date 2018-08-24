package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"test-task/utils"
	"testing"
)

func TestWorker(t *testing.T) {

	type testJSONBody struct {
		Test string `json:"test"`
	}

	testHandler := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Test", "ok")
		if req.Method == "POST" {
			decoder := json.NewDecoder(req.Body)
			var message testJSONBody
			err := decoder.Decode(&message)
			if err != nil {
				utils.ResponseError(w, http.StatusInternalServerError, err.Error())
				return
			}
			if message.Test == "ok" {
				w.Header().Set("Testjson", message.Test)
			}
		}
	}

	testServer := httptest.NewServer(http.HandlerFunc(testHandler))
	defer testServer.Close()

	taskQueue := make(chan TaskRequest, 1)
	closeCh := make(chan struct{})
	testRepo := initRepo()

	var testRequest = TaskRequest{
		ID:     "testID",
		Method: "GET",
		URL:    testServer.URL,
	}

	worker := worker{
		workerID:  1,
		repo:      &testRepo,
		taskQueue: taskQueue,
		closed:    closeCh,
	}

	var jsbody = testJSONBody{"ok"}
	testRequest.Method = "POST"
	data, _ := json.Marshal(jsbody)
	testRequest.Body = string(data)

	worker.start()
	taskQueue <- testRequest
	close(taskQueue)
	<-closeCh

	result, err := testRepo.Get("testID")
	if err != nil {
		t.Error("No entry in the repository", err)
	}

	if result.Status != "200 OK" {
		t.Errorf("Wrong response status: %s", result.Status)
	}

	testHeader, ok := result.Headers["Test"]
	if !ok || testHeader[0] != "ok" {
		t.Error("Wrong or invalid header")
	}

	testJSONHeader, ok := result.Headers["Testjson"]
	if !ok || testJSONHeader[0] != "ok" {
		t.Errorf("Wrong or invalid header %v", result)
	}

}
