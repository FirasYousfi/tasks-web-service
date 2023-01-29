package service

// INFO Important you can see it does not depend on the repository but on the interface that the repo implements
import (
	"github.com/FirasYousfi/tasks-web-servcie/application/interfaces"
	"github.com/FirasYousfi/tasks-web-servcie/domain/entity"
	"github.com/FirasYousfi/tasks-web-servcie/domain/validation"
	"github.com/google/uuid"
	"log"
)

// Service The attributes should be the dependencies needed from the outer layer's stuff, those will be injected
// DIP happens here,
type Service struct {
	Repository interfaces.IRepository
}

// NewTaskService Dependency Inversion Principle. DIP suggests that we should depend on abstractions (interfaces), not concrete classes.
// also that way we respect the Dependency Rule. This rule says that source code dependencies can only point inwards.
// Inner circles never mention a name in an outer circle. Repository impl is in outer circle, but the interfaces are in the app layer.
func NewTaskService(repo interfaces.IRepository) *Service {
	if repo == nil {
		log.Fatalf("nil repo provided")
	}
	return &Service{Repository: repo}
}

func (t *Service) CreateTask(req *entity.TaskDescription) (*entity.Task, error) {
	description, err := validation.ValidateParams(req)
	if err != nil {
		return nil, err
	}

	task := entity.Task{TaskDescription: *description}
	log.Printf("creating task with ID '%s' ...", task.ID)

	err = t.Repository.CreateTask(&task)
	if err != nil {
		return nil, err
	}
	return &task, err
}

func (t *Service) GetTasks() ([]*entity.Task, error) {
	log.Printf("listing all tasks ...")
	return t.Repository.FindAllTasks()
}

func (t *Service) DeleteTaskByID(id string) error {
	log.Printf("deleting task with id '%s' ...", id)
	return t.Repository.DeleteTaskByID(id)
}

func (t *Service) GetTaskByID(id string) (*entity.Task, error) {
	log.Printf("getting task with id '%s' ...", id)
	return t.Repository.FindTaskByID(id)
}

func (t *Service) UpdateTaskFully(req *entity.TaskDescription, id string) (*entity.Task, error) {
	log.Printf("updating task with id '%s' ...", id)
	_, err := t.Repository.FindTaskByID(id)
	if err != nil {
		return nil, err
	}
	request, err := validation.ValidateParams(req)
	if err != nil {
		return nil, err
	}

	values := map[string]interface{}{"title": request.Title, "description": request.Description, "priority": request.Priority, "status": request.Status}
	err = t.Repository.UpdateTask(values, id)
	if err != nil {
		return nil, err
	}

	task, err := t.Repository.FindTaskByID(id)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (t *Service) UpdateTaskPartial(req *entity.TaskDescription, id string) (*entity.Task, error) {
	log.Printf("updating task with id '%s' ...", id)
	_, err := t.Repository.FindTaskByID(id)
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
	err = t.Repository.UpdateTask(values, id)
	if err != nil {
		return nil, err
	}

	task, err := t.Repository.FindTaskByID(id)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return task, nil
}

func (t *Service) CreateCollection(description *entity.CollectionDescription) (*entity.Collection, error) {
	collection := entity.Collection{ID: uuid.NewString(), CollectionDescription: *description}
	log.Printf("creating collection with ID '%s' ...", collection.ID)

	err := t.Repository.CreateCollection(&collection)
	if err != nil {
		return nil, err
	}
	return &collection, err
}

func (t *Service) GetCollections() ([]*entity.Collection, error) {
	log.Printf("listing collections ...")
	return t.Repository.FindAllCollections()
}

func (t *Service) GetCollectionsByID(id string) (*entity.Collection, error) {
	log.Printf("getting collection with id '%s' ...", id)
	return t.Repository.FindCollectionByID(id)
}

func (t *Service) DeleteCollectionByID(id string) error {
	log.Printf("deleting collection with id '%s' ...", id)
	return t.Repository.DeleteCollectionByID(id)
}
