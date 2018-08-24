package main

import (
	"log"
	"net/http"
)

var taskQueue = make(chan TaskRequest, 100)
var exitCh = make(chan struct{})

func main() {
	port := ":8080"
	repo := initRepo()
	numWorkers := 5

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

	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(port, new(&repo)))
}
