package main

import (
	"time"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

// TaskRequest represents request for a task
type TaskRequest struct {
	ID      string            `json:"id,omitempty"`
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
}

// Validate is a TaskRequest validator
func (t TaskRequest) Validate() error {
	return validation.ValidateStruct(&t,
		validation.Field(&t.Method, validation.Required, validation.In(
			"GET",
			"POST",
			"PUT",
			"HEAD",
			"DELETE",
			"PATCH",
		)),
		validation.Field(&t.URL, validation.Required, is.URL),
	)
}

// CompletedTask represents a completed task
type CompletedTask struct {
	ID      string              `json:"id,omitempty"`
	Status  string              `json:"status,omitempty"`
	Headers map[string][]string `json:"headers,omitempty"`
	Size    int64               `json:"size,omitempty"`
	Errors  []string            `json:"error,omitempty"`
	Date    time.Time           `json:"-"`
}
