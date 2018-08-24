package main

import (
	"strconv"
	"testing"
)

func TestInMemoryRepo(t *testing.T) {
	repo := initRepo()
	var testTask = CompletedTask{
		ID:     "testID",
		Status: "test",
	}

	repo.Set(testTask)
	task, err := repo.Get(testTask.ID)

	if err != nil {
		t.Error(err)
	}

	if task.Status != testTask.Status {
		t.Error("Wrong task data")
	}

	repo.Delete(testTask.ID)
	_, err = repo.Get(testTask.ID)
	if err.Error() != "No such entry" {
		t.Error("Delete operation do not work correctly")
	}

	for i := 0; i < 10; i++ {
		id := i
		testTask.ID = strconv.Itoa(id)
		repo.Set(testTask)
	}

	tasks := repo.List()
	if len(tasks) != 10 {
		t.Errorf("List should return 10 enties, returned %d", len(tasks))
	}
}
