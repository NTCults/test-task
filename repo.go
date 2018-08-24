package main

import (
	"errors"
	"sort"
	"sync"
)

// Repository interface
type Repository interface {
	Set(CompletedTask)
	Get(string) (*CompletedTask, error)
	Delete(string) error
	List() []CompletedTask
}

type inMemoryRepo struct {
	mu   *sync.RWMutex
	data map[string]CompletedTask
}

func initRepo() inMemoryRepo {
	return inMemoryRepo{mu: &sync.RWMutex{}, data: make(map[string]CompletedTask)}
}

func (repo *inMemoryRepo) Set(t CompletedTask) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	repo.data[t.ID] = t
}

func (repo *inMemoryRepo) Get(id string) (*CompletedTask, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	res, ok := repo.data[id]
	if !ok {
		return nil, errors.New("No such entry")
	}
	return &res, nil
}

func (repo *inMemoryRepo) Delete(id string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	_, ok := repo.data[id]
	if !ok {
		return errors.New("No such entry")
	}
	delete(repo.data, id)
	return nil
}

func (repo *inMemoryRepo) List() []CompletedTask {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	sortedTasks := make([]CompletedTask, 0, len(repo.data))

	for _, t := range repo.data {
		sortedTasks = append(sortedTasks, t)
	}
	sort.Slice(sortedTasks, func(i, j int) bool { return sortedTasks[i].Date.Before(sortedTasks[j].Date) })
	return sortedTasks
}
