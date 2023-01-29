package repository

import (
	"fmt"
	"github.com/FirasYousfi/tasks-web-servcie/domain/entity"
	"gorm.io/gorm"
	"log"
)

// Repository The attributes should be the dependencies needed from the outer layer's stuff, those will be injected
type Repository struct {
	db *gorm.DB
}

// NewRepository is the constructor of a Repository with the database dependency injected
func NewRepository(db *gorm.DB) *Repository {
	if db == nil {
		log.Fatalf("nil db provided")
	}
	return &Repository{db: db}
}

// CreateTask creates a new task in the database
func (t *Repository) CreateTask(task *entity.Task) error {
	tx := t.db.Create(task)
	return tx.Error
}

// FindAllTasks returns all the tasks in the database
func (t *Repository) FindAllTasks() ([]*entity.Task, error) {
	var tasks []*entity.Task
	// SELECT * FROM users;
	tx := t.db.Find(&tasks) // pointer to our array because it needs to be modified
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tasks, nil
}

// DeleteTaskByID Deletes a task identified by its uuid given as parameter
func (t *Repository) DeleteTaskByID(id string) error {
	tx := t.db.Where("id = ?", id).Delete(&entity.Task{})

	return tx.Error
}

// FindTaskByID Finds a task identified by its uuid given as parameter
func (t *Repository) FindTaskByID(id string) (*entity.Task, error) {
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

// UpdateTask updates a task by the new values passed as parameters. The ID of the task to update would be part of the task given as argument.
// When update with struct, GORM will only update non-zero fields. So better use map to make sure.
// Will be used for both patch and PUT, checking for empty values will be done in the Service function.
func (t *Repository) UpdateTask(fields map[string]interface{}, id string) error {
	tx := t.db.Model(entity.Task{}).Where("id = ?", id).Updates(fields)
	return tx.Error
}

// CreateCollection creates a new collection in the database
func (t *Repository) CreateCollection(collection *entity.Collection) error {
	tx := t.db.Create(collection)
	return tx.Error
}

// FindAllCollections returns all the collections in the database
func (t *Repository) FindAllCollections() ([]*entity.Collection, error) {
	var collections []*entity.Collection
	// SELECT * FROM users;
	tx := t.db.Model(entity.Collection{}).Preload("Tasks").Find(&collections) // pointer to our array because it needs to be modified
	fmt.Println("here are the collections", collections)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return collections, nil
}

// DeleteCollectionByID Deletes a collection identified by its uuid given as parameter
func (t *Repository) DeleteCollectionByID(id string) error {
	tx := t.db.Where("id = ?", id).Delete(&entity.Collection{})

	return tx.Error
}

// FindCollectionByID Finds a collections identified by its uuid given as parameter
func (t *Repository) FindCollectionByID(id string) (*entity.Collection, error) {
	var collection entity.Collection //This is necessary, should not create pointer and pass it directly
	tx := t.db.Where("id = ?", id).First(&collection)
	if &collection == nil {
		return nil, fmt.Errorf("could not find collection")
	}
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &collection, nil
}

// UpdateCollection updates a collection with new fields
func (t *Repository) UpdateCollection(fields map[string]interface{}, id string) error {
	tx := t.db.Model(entity.Task{}).Where("id = ?", id).Updates(fields)
	return tx.Error
}
