package interfaces

import "github.com/FirasYousfi/tasks-web-servcie/domain/entity"

// IRepository defines the CRUD operations that are done to the database
type IRepository interface {
	TaskReadWriter
	CollectionReadWriter
}

type TaskReadWriter interface {
	TaskWriter
	TaskReader
}

type CollectionReadWriter interface {
	CollectionWriter
	CollectionReader
}

type TaskWriter interface {
	CreateTask(task *entity.Task) error
	DeleteTaskByID(id string) error
	UpdateTask(fields map[string]interface{}, id string) error
}

type TaskReader interface {
	FindAllTasks() ([]*entity.Task, error)
	FindTaskByID(id string) (*entity.Task, error)
}

type CollectionWriter interface {
	CreateCollection(task *entity.Collection) error
	DeleteCollectionByID(id string) error
}

type CollectionReader interface {
	FindAllCollections() ([]*entity.Collection, error)
	FindCollectionByID(id string) (*entity.Collection, error)
}
