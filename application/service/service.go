package service

// INFO Important you can see it does not depend on the repository but on the interface that the repo implements
import (
	"github.com/FirasYousfi/tasks-web-servcie/application/interfaces"
	"github.com/FirasYousfi/tasks-web-servcie/domain/entity"
	"github.com/FirasYousfi/tasks-web-servcie/domain/validation"
	"github.com/google/uuid"
	"log"
)

// TaskService The attributes should be the dependencies needed from the outer layer's stuff, those will be injected
// DIP happens here,
type TaskService struct {
	TaskRepository interfaces.ITaskRepository
}

// NewTaskService Dependency Inversion Principle. DIP suggests that we should depend on abstractions (interfaces), not concrete classes.
// => also that way we respect the Dependency Rule. This rule says that source code dependencies can only point inwards.
// Inner circles never mention a name in an outer circle. Repository impl is in outer circle, but the interfaces are in the app layer.
func NewTaskService(repo interfaces.ITaskRepository) *TaskService {
	if repo == nil {
		log.Fatalf("nil repo provided")
	}
	return &TaskService{TaskRepository: repo}
}

func (t *TaskService) Create(req *entity.TaskDescription) (*entity.Task, error) {
	description, err := validation.ValidateParams(req)
	if err != nil {
		return nil, err
	}

	task := entity.Task{ID: uuid.NewString(), TaskDescription: *description}
	log.Printf("creating task with ID '%s' ...", task.ID)

	err = t.TaskRepository.Create(&task)
	if err != nil {
		return nil, err
	}
	return &task, err
}

func (t *TaskService) Get() ([]*entity.Task, error) {
	log.Printf("listing all tasks ...")
	return t.TaskRepository.FindAll()
}

func (t *TaskService) DeleteByID(id string) error {
	log.Printf("deleting task with id '%s' ...", id)
	return t.TaskRepository.DeleteByID(id)
}

func (t *TaskService) GetByID(id string) (*entity.Task, error) {
	log.Printf("getting task with id '%s' ...", id)
	return t.TaskRepository.FindByID(id)
}

func (t *TaskService) UpdateFully(req *entity.TaskDescription, id string) (*entity.Task, error) {
	log.Printf("updating task with id '%s' ...", id)
	_, err := t.TaskRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	request, err := validation.ValidateParams(req)
	if err != nil {
		return nil, err
	}

	values := map[string]interface{}{"title": request.Title, "description": request.Description, "priority": request.Priority, "status": request.Status}
	err = t.TaskRepository.Update(values, id)
	if err != nil {
		return nil, err
	}

	task, err := t.TaskRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (t *TaskService) UpdatePartial(req *entity.TaskDescription, id string) (*entity.Task, error) {
	log.Printf("updating task with id '%s' ...", id)
	_, err := t.TaskRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	values := make(map[string]interface{})
	if req.Title != "" {
		err := validation.ValidateTitle(req.Title)
		if err != nil {
			return nil, err
		}
		values["title"] = req.Title
	}
	if req.Description != "" {
		err := validation.ValidateDescription(req.Description)
		if err != nil {
			return nil, err
		}
		values["description"] = req.Description
	}
	if req.Priority != 0 {
		err := validation.ValidatePriority(req.Priority)
		if err != nil {
			return nil, err
		}
		values["Priority"] = req.Priority
	}

	if req.Status != "" {
		values["status"] = req.Status
	}
	err = t.TaskRepository.Update(values, id)
	if err != nil {
		return nil, err
	}

	task, err := t.TaskRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return task, nil
}
