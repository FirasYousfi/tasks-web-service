package handlers

import (
	"github.com/FirasYousfi/tasks-web-servcie/application/interfaces"
	"github.com/FirasYousfi/tasks-web-servcie/domain/entity"
)

// CreateTask is the struct used for that will implement the Handler interface for the creation
type CreateTask struct {
	req     entity.TaskDescription
	res     entity.Task
	Service interfaces.IService
}

type DeleteTask struct {
	Service interfaces.IService
}

// GetTask In case some response type or sth similar is needed in the future
type GetTask struct {
	res     entity.Task
	Service interfaces.IService
}

// ListTasks In case some response type or sth similar is needed in the future
type ListTasks struct {
	res     []*entity.Task
	Service interfaces.IService
}

// UpdateTask represents the struct that implement the handler for update
type UpdateTask struct {
	req     entity.TaskDescription
	res     entity.Task
	Service interfaces.IService
}
