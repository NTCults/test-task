package main

import (
	"log"
	"net/http"
	"strings"
	"time"
)

type worker struct {
	workerID    int
	repo        Repository
	taskQueue   <-chan TaskRequest
	closed      chan struct{}
	queueClosed bool
}

func (w *worker) start() {
	go func() {
		for {
			task, ok := <-w.taskQueue
			if !ok {
				w.closed <- struct{}{}
				return
			}
			var completedTask CompletedTask
			completedTask.ID = task.ID

			client := &http.Client{}

			req, err := http.NewRequest(task.Method, task.URL, strings.NewReader(task.Body))
			if err != nil {
				log.Printf("[Worker %d] Task %s: %s", w.workerID, task.ID, err)
				completedTask.Errors = append(completedTask.Errors, err.Error())
			}

			for k, v := range task.Headers {
				req.Header.Add(k, v)
			}

			resp, err := client.Do(req)
			if err != nil {
				log.Printf("[Worker %d] Task %s: %s", w.workerID, task.ID, err)
				completedTask.Errors = append(completedTask.Errors, err.Error())
			}

			if completedTask.Errors == nil {
				completedTask = parseResponse(completedTask, resp)
			}
			completedTask.Date = time.Now()
			w.repo.Set(completedTask)
		}
	}()
}

func parseResponse(task CompletedTask, r *http.Response) CompletedTask {
	task.Status = r.Status
	task.Size = r.ContentLength
	task.Headers = make(map[string][]string)
	for k, v := range r.Header {
		task.Headers[k] = v
	}
	return task
}
