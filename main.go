package main

import (
	"fmt"
	"log"
	"net/http"
	"test-task/utils"
	"time"
)

func main() {
	port := utils.Getenv("PORT", 8080)
	numWorkers := utils.Getenv("WORKERS", 5)
	timeout := utils.Getenv("TIMEOUT", 15)
	queueSize := utils.Getenv("QUEUE_SIZE", 100)

	var taskQueue = make(chan TaskRequest, queueSize)
	var exitCh = make(chan struct{})

	repo := initRepo()
	for i := 0; i < numWorkers; i++ {
		workerID := i
		wr := worker{
			workerID:  workerID,
			repo:      &repo,
			taskQueue: taskQueue,
			closed:    exitCh,
		}
		wr.start()
	}
	srv := &http.Server{
		Handler:      newService(&repo, taskQueue),
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: time.Duration(timeout) * time.Second,
		ReadTimeout:  time.Duration(timeout) * time.Second,
	}

	log.Println("Listening on port", port)
	log.Fatal(srv.ListenAndServe())
}
