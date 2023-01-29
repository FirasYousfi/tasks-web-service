package handlers

import (
	"bytes"
	"github.com/FirasYousfi/tasks-web-servcie/domain/entity"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

// INFO, since the DBs are being modified by different handlers, when running all tests together they fail when they share the same resource.
// For this purpose I created another DB for each handler test.
var updateTaskDB = []*entity.Task{{
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

var testUpdateTask entity.Task = entity.Task{
	ID: "1",
	TaskDescription: entity.TaskDescription{
		Title:       "update",
		Description: "update",
		Priority:    2,
		Status:      "New",
	},
}

func TestUpdate_ServeHTTP(t *testing.T) {
	var (
		taskService         = newMockTaskService(updateTaskDB)
		validPutReq         = httptest.NewRequest("PUT", "http://localhost:8080/v1/api/tasks/1", bytes.NewReader(reqBodyToJson(testUpdateTask.TaskDescription, t)))
		validPatchReq       = httptest.NewRequest("PATCH", "http://localhost:8080/v1/api/tasks/1", bytes.NewReader(reqBodyToJson(testUpdateTask.TaskDescription, t)))
		methodNotAllowedReq = httptest.NewRequest("GET", "http://localhost:8080/v1/api/tasks/1", bytes.NewReader(reqBodyToJson(testUpdateTask.TaskDescription, t)))
		invalidBodyReq      = httptest.NewRequest("PUT", "http://localhost:8080/v1/api/tasks/1", strings.NewReader("no-json"))
		invalidPathParamReq = httptest.NewRequest("PUT", "http://localhost:8080/v1/api/tasks", strings.NewReader("no-json"))
	)

	type want struct {
		body   entity.Task
		status int
	}
	tests := []struct {
		name    string
		fields  UpdateTask
		request *http.Request
		want    want
	}{
		{
			name: "should update item fully successfully",
			fields: UpdateTask{ // no need to fill it with req and res because those get populated inside the handler itself
				TaskService: taskService,
			},
			request: mux.SetURLVars(validPutReq, map[string]string{"id": testUpdateTask.ID}),
			want: want{
				body:   testUpdateTask,
				status: http.StatusOK,
			},
		},
		{
			name: "should update item partially successfully",
			fields: UpdateTask{ // no need to fill it with req and res because those get populated inside the handler itself
				TaskService: taskService,
			},
			request: mux.SetURLVars(validPatchReq, map[string]string{"id": testUpdateTask.ID}),
			want: want{
				body:   testUpdateTask,
				status: http.StatusOK,
			},
		},
		{
			name: "should fail to create tasks with StatusMethodNotAllowed",
			fields: UpdateTask{
				TaskService: taskService,
			},
			request: mux.SetURLVars(methodNotAllowedReq, map[string]string{"id": testUpdateTask.ID}),
			want: want{
				body:   entity.Task{},
				status: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "should fail to create because body is not a json",
			fields: UpdateTask{
				TaskService: taskService,
			},
			request: mux.SetURLVars(invalidBodyReq, map[string]string{"id": testUpdateTask.ID}),
			want: want{
				body:   entity.Task{},
				status: http.StatusBadRequest,
			},
		},
		{
			name: "should fail because no id provided in path params",
			fields: UpdateTask{
				TaskService: taskService,
			},
			request: mux.SetURLVars(invalidPathParamReq, map[string]string{"id": ""}),
			want: want{
				body:   entity.Task{},
				status: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We set the ID path parameter here for testing purposes, because it is not being read from the url
			//I moved it to the definition of the request in the testTable.
			//tt.request = mux.SetURLVars(tt.request, map[string]string{"id": testUpdateTask.ID})

			// Having a recorder for each separate request is safer!
			response := httptest.NewRecorder()

			u := UpdateTask{
				req:         tt.fields.req,
				res:         tt.fields.res,
				TaskService: tt.fields.TaskService,
			}
			u.ServeHTTP(response, tt.request)

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
