package handlers

import (
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDelete_ServeHTTP(t *testing.T) {
	var (
		taskService         = newMockTaskService(tasksDatabase)
		validReq            = httptest.NewRequest("DELETE", "http://localhost:8080/v1/api/tasks/1", nil)
		methodNotAllowedReq = httptest.NewRequest("PUT", "http://localhost:8080/v1/api/tasks/1", nil)
	)

	type want struct {
		status int
	}
	tests := []struct {
		name    string
		fields  DeleteTask
		request *http.Request
		want    want
	}{
		{
			name: "should delete item successfully",
			fields: DeleteTask{
				TaskService: taskService,
			},
			request: mux.SetURLVars(validReq, map[string]string{"id": testUpdateTask.ID}),
			want: want{
				status: http.StatusNoContent,
			},
		},
		{
			name: "should fail to delete task with StatusMethodNotAllowed",
			fields: DeleteTask{
				TaskService: taskService,
			},
			request: mux.SetURLVars(methodNotAllowedReq, map[string]string{"id": testUpdateTask.ID}),
			want: want{
				status: http.StatusMethodNotAllowed,
			},
		},
		{
			name: "should fail because no id provided as path parameter",
			fields: DeleteTask{
				TaskService: taskService,
			},
			request: mux.SetURLVars(validReq, map[string]string{"id": ""}),
			want: want{
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

			d := DeleteTask{
				TaskService: tt.fields.TaskService,
			}
			d.ServeHTTP(response, tt.request)

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
		})
	}
}
