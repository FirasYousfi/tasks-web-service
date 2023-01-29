package interfaces

import "github.com/FirasYousfi/tasks-web-servcie/domain/entity"

type IService interface {
	TaskService
	CollectionService
}

// TaskService defines the functions needed for the use-cases, they should contain all the business logic needed to fulfill the services required from the user.
type TaskService interface {
	CreateTask(task *entity.TaskDescription) (*entity.Task, error)
	GetTasks() ([]*entity.Task, error)
	GetTaskByID(id string) (*entity.Task, error)
	DeleteTaskByID(id string) error
	UpdateTaskPartial(task *entity.TaskDescription, id string) (*entity.Task, error)
	UpdateTaskFully(task *entity.TaskDescription, id string) (*entity.Task, error)
}

// CollectionService defines the functions needed for the use-cases, they should contain all the business logic needed to fulfill the services required from the user.
type CollectionService interface {
	CreateCollection(collection *entity.CollectionDescription) (*entity.Collection, error)
	GetCollections() ([]*entity.Collection, error)
	GetCollectionsByID(id string) (*entity.Collection, error)
	DeleteCollectionByID(id string) error
}
