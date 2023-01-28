package repository

import (
	"fmt"
	"github.com/FirasYousfi/tasks-web-servcie/domain/entity"
	"gorm.io/gorm"
	"log"
)

// TaskRepository The attributes should be the dependencies needed from the outer layer's stuff, those will be injected
type TaskRepository struct {
	db *gorm.DB
}

// NewTaskRepository is the constructor of a TaskRepository with the database dependency injected
func NewTaskRepository(db *gorm.DB) *TaskRepository {
	if db == nil {
		log.Fatalf("nil db provided")
	}
	return &TaskRepository{db: db}
}

// Create creates a new task in the database
func (t *TaskRepository) Create(task *entity.Task) error {
	tx := t.db.Create(task)
	return tx.Error
}

// FindAll returns all the tasks in the database
func (t *TaskRepository) FindAll() ([]*entity.Task, error) {
	var tasks []*entity.Task
	// SELECT * FROM users;
	tx := t.db.Find(&tasks) // pointer to our array because it needs to be modified
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tasks, nil
}

// DeleteByID Deletes a task identified by its uuid given as parameter
func (t *TaskRepository) DeleteByID(id string) error {
	tx := t.db.Where("id = ?", id).Delete(&entity.Task{})

	return tx.Error
}

// FindByID Finds a task identified by its uuid given as parameter
func (t *TaskRepository) FindByID(id string) (*entity.Task, error) {
	var task entity.Task //This is necessary, should not create pointer and pass it directly
	tx := t.db.Where("id = ?", id).First(&task)
	if &task == nil {
		return nil, fmt.Errorf("could not find task")
	}
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &task, nil
}

// Update updates a task by the new values passed as parameters. The ID of the task to update would be part of the task given as argument.
// When update with struct, GORM will only update non-zero fields. So better use map to make sure.
// Will be used for both patch and PUT, checking for empty values will be done in the Service function.
func (t *TaskRepository) Update(fields map[string]interface{}, id string) error {
	tx := t.db.Model(entity.Task{}).Where("id = ?", id).Updates(fields)
	return tx.Error
}
