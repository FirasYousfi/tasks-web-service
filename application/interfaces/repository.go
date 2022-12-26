package interfaces

import "github.com/FirasYousfi/tasks-web-servcie/domain/entity"

// ITaskRepository defines the CRUD operations that are done to the database
type ITaskRepository interface {
	WriterRepository
	ReaderRepository
}

type WriterRepository interface {
	Create(task *entity.Task) error
	DeleteByID(id string) error
	Update(fields map[string]interface{}, id string) error
}

type ReaderRepository interface {
	FindAll() ([]*entity.Task, error)
	FindByID(id string) (*entity.Task, error)
}
