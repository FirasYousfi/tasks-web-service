package handlers

import (
	"github.com/FirasYousfi/tasks-web-servcie/domain/entity"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var getTaskDB = []*entity.Task{{
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

var testGetTask = *getTaskDB[0]

func TestGet_ServeHTTP(t *testing.T) {
	var (
		taskService         = newMockTaskService(getTaskDB)
		validReq            = httptest.NewRequest("GET", "http://localhost:8080/v1/api/tasks/1", nil)
		methodNotAllowedReq = httptest.NewRequest("PUT", "http://localhost:8080/v1/api/tasks/1", nil)
	)

	type want struct {
		body   entity.Task
		status int
	}
	tests := []struct {
		name    string
		fields  Get
		request *http.Request
		want    want
	}{
		{
			name: "should get item successfully",
			fields: Get{ // no need to fill it with req and res because those get populated inside the handler itself
				TaskService: taskService,
			},
			request: mux.SetURLVars(validReq, map[string]string{"id": testUpdateTask.ID}),
			want: want{
				body:   testGetTask,
				status: http.StatusOK,
			},
		},
		{
			name: "should fail to get task with StatusMethodNotAllowed",
			fields: Get{
				TaskService: taskService,
			},
			request: mux.SetURLVars(methodNotAllowedReq, map[string]string{"id": testUpdateTask.ID}),
			want: want{
				body:   entity.Task{},
				status: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "should fail because no id provided as path parameter",
			fields: Get{
				TaskService: taskService,
			},
			request: mux.SetURLVars(validReq, map[string]string{"id": ""}),
			want: want{
				body:   entity.Task{},
				status: http.StatusBadRequest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We set the ID path parameter here for testing purposes, because it is not being read from the url
			//tt.request = mux.SetURLVars(tt.request, map[string]string{"id": testGetTask.ID})

			// Having a recorder for each separate request is safer!
			response := httptest.NewRecorder()

			u := Get{
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
