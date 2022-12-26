package interfaces

import "github.com/FirasYousfi/tasks-web-servcie/domain/entity"

// ITaskService defines the functions needed for the use-cases, they should contain all the business logic needed to fulfill the services required from the user.
type ITaskService interface {
	Create(task *entity.TaskDescription) (*entity.Task, error)
	Get() ([]*entity.Task, error)
	GetByID(id string) (*entity.Task, error)
	DeleteByID(id string) error
	UpdatePartial(task *entity.TaskDescription, id string) (*entity.Task, error)
	UpdateFully(task *entity.TaskDescription, id string) (*entity.Task, error)
}
