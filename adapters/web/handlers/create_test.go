package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/FirasYousfi/tasks-web-servcie/domain/entity"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

var testCreateTask entity.Task = entity.Task{
	ID: "test-create-id",
	TaskDescription: entity.TaskDescription{
		Title:       "test-create-title",
		Description: "test-create-description",
		Priority:    5,
		Status:      entity.Active,
	},
}

var tasksDatabase = []*entity.Task{{
	ID: "1",
	TaskDescription: entity.TaskDescription{
		Title:       "test1",
		Description: "test1",
		Priority:    0,
	},
}, {
	ID: "2",
	TaskDescription: entity.TaskDescription{
		Title:       "test2",
		Description: "test2",
		Priority:    0,
	},
}}

type mockTaskService struct {
	tasks []*entity.Task
}

func newMockTaskService(tasks []*entity.Task) *mockTaskService {
	return &mockTaskService{tasks}
}

func (t mockTaskService) CreateTask(taskDescription *entity.TaskDescription) (*entity.Task, error) {
	task := entity.Task{
		ID:              testCreateTask.ID,
		TaskDescription: *taskDescription,
	}
	t.tasks = append(t.tasks, &task)
	return &task, nil
}

func (t mockTaskService) GetTasks() ([]*entity.Task, error) {
	return t.tasks, nil
}

func (t mockTaskService) GetTaskByID(id string) (*entity.Task, error) {
	for i := range t.tasks {
		if t.tasks[i].ID == id {
			return t.tasks[i], nil
		}
	}
	return nil, fmt.Errorf("element with ID %s not found", id)
}

func (t mockTaskService) DeleteTaskByID(id string) error {
	for i := range t.tasks {
		if t.tasks[i].ID == id {
			t.tasks[i] = t.tasks[len(t.tasks)-1] // put last element there since order does not matter
			t.tasks = t.tasks[:len(t.tasks)-1]   // Truncate slice by erasing last
			return nil
		}
	}
	return fmt.Errorf("element with ID %s not found", id)
}

// we tested the functionality already in the service package, so no need to put in a lot of logic in this simple mock
func (t mockTaskService) UpdateTaskPartial(taskDescription *entity.TaskDescription, id string) (*entity.Task, error) {
	for i := range t.tasks {
		if t.tasks[i].ID == id {
			t.tasks[i].TaskDescription = *taskDescription
			return t.tasks[i], nil
		}
	}
	return nil, fmt.Errorf("element with ID %s not found", id)
}

func (t mockTaskService) UpdateTaskFully(taskDescription *entity.TaskDescription, id string) (*entity.Task, error) {
	for i := range t.tasks {
		if t.tasks[i].ID == id {
			t.tasks[i].TaskDescription = *taskDescription
			return t.tasks[i], nil
		}
	}
	return nil, fmt.Errorf("element with ID %s not found", id)
}

func TestCreate_ServeHTTP(t *testing.T) {
	var (
		taskService          = newMockTaskService(tasksDatabase)
		validReq             = httptest.NewRequest("POST", "http://localhost:8080/v1/api/tasks", bytes.NewReader(reqBodyToJson(testCreateTask.TaskDescription, t)))
		methodNotAllowedReq  = httptest.NewRequest("PATCH", "http://localhost:8080/v1/api/tasks", bytes.NewReader(reqBodyToJson(testCreateTask.TaskDescription, t)))
		invalidBodyFormatReq = httptest.NewRequest("POST", "http://localhost:8080/v1/api/tasks/1", strings.NewReader("no-json"))
	)

	type want struct {
		body   entity.Task
		status int
	}
	tests := []struct {
		name    string
		fields  CreateTask
		request *http.Request
		want    want
	}{
		{
			name: "should create item successfully",
			fields: CreateTask{ // no need to fill it with req and res because those get populated inside the handler itself
				TaskService: taskService,
			},
			request: validReq,
			want: want{
				body:   testCreateTask,
				status: http.StatusCreated,
			},
		},
		{
			name: "should fail to create tasks with StatusMethodNotAllowed",
			fields: CreateTask{
				TaskService: taskService,
			},
			request: methodNotAllowedReq,
			want: want{
				body:   entity.Task{},
				status: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "should fail to create because body is not a json",
			fields: CreateTask{
				TaskService: taskService,
			},
			request: invalidBodyFormatReq,
			want: want{
				body:   entity.Task{},
				status: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Having a recorder for each separate request is safer!
			response := httptest.NewRecorder()

			c := CreateTask{
				req:         tt.fields.req,
				res:         tt.fields.res,
				TaskService: tt.fields.TaskService,
			}
			c.ServeHTTP(response, tt.request)

			//here we check the stuff that our recorder recorded after request was handled
			res := response.Result()
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					t.Fatal(err)
				}
			}(res.Body)

			//assert status code
			if res.StatusCode != tt.want.status {
				t.Errorf("invalid status code, expected: %d, got: %d", tt.want.status, res.StatusCode)
			}

			got, err := readTaskResponse(response)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.want.body) {
				t.Errorf("invalid response, expected: %v, got: %v", tt.want.body, got)
			}

		})
	}
}

func readTaskResponse(response *httptest.ResponseRecorder) (entity.Task, error) {
	if response.Body.Len() == 0 { // this way when method not allowed and we got no body in response, this function would not return an error
		return entity.Task{}, nil
	}
	var got entity.Task

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return entity.Task{}, err
	}
	err = json.Unmarshal(responseData, &got)
	if err != nil {
		return entity.Task{}, err
	}
	return got, nil
}

func reqBodyToJson(taskDescription entity.TaskDescription, t *testing.T) []byte {
	toByte, err := json.Marshal(taskDescription)
	if err != nil {
		t.Fatal(err)
	}
	return toByte
}
