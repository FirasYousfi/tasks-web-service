package entity

import (
	"gorm.io/gorm"
)

// string mapping with the possible values for status
const (
	New    Status = "new"
	Active Status = "active"
	Closed Status = "closed"
	OnHold Status = "on-hold"
)

// Status represents the current state of the task
type Status string

// Task Represents the whole task that will be modeled with gorm DB
type Task struct {
	gorm.Model
	CollectionID int `json:"collectionId"`
	TaskDescription
}

// TaskDescription represents the description of the task to be created. Those are the values that the user can set.
type TaskDescription struct {
	Title       string `json:"title"`                                         // title of the task
	Description string `json:"description"`                                   // description of the task
	Priority    int    `json:"priority" minimum:"1" maximum:"10" default:"1"` // priority is represented by an int from 1 to 10
	Status      Status `json:"status"`                                        // current status of the task
}

// Collection represents a collection of tasks
type Collection struct {
	gorm.Model
	CollectionDescription
}

// CollectionDescription represents a collection of tasks
type CollectionDescription struct {
	Name  string `json:"name"`
	Tasks []Task `json:"tasks"`
}
