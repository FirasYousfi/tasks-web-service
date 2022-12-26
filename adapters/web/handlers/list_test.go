package handlers

import (
	"encoding/json"
	"github.com/FirasYousfi/tasks-web-servcie/domain/entity"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestList_ServeHTTP(t *testing.T) {
	var (
		taskService         = newMockTaskService(tasksDatabase)
		validReq            = httptest.NewRequest("GET", "http://localhost:8080/v1/api/tasks", nil)
		methodNotAllowedReq = httptest.NewRequest("PUT", "http://localhost:8080/v1/api/tasks", nil)
	)

	type want struct {
		body   []*entity.Task
		status int
	}
	tests := []struct {
		name    string
		fields  List
		request *http.Request
		want    want
	}{
		{
			name: "should list tasks successfully",
			fields: List{
				TaskService: taskService,
			},
			request: validReq,
			want: want{
				body:   taskService.tasks,
				status: http.StatusOK,
			},
		},
		{
			name: "should fail to list tasks with StatusMethodNotAllowed",
			fields: List{
				TaskService: taskService,
			},
			request: methodNotAllowedReq,
			want: want{
				body:   []*entity.Task{},
				status: http.StatusMethodNotAllowed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Having a recorder for each separate request is safer!
			response := httptest.NewRecorder()

			l := List{
				res:         tt.fields.res,
				TaskService: tt.fields.TaskService,
			}

			l.ServeHTTP(response, tt.request)

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

			//assert response body
			got, err := readListBody(response)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.want.body) {
				t.Errorf("invalid response, expected: %v, got: %v", tt.want.body, got)
			}

		})
	}
}

func readListBody(response *httptest.ResponseRecorder) ([]*entity.Task, error) {
	if response.Body.Len() == 0 { // this way when method not allowed and we got no body in response, this function would not return an error
		return []*entity.Task{}, nil
	}

	var got []*entity.Task

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(responseData, &got)
	if err != nil {
		return nil, err
	}
	return got, nil
}
