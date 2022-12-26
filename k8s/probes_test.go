package k8s

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLiveness_ServeHTTP(t *testing.T) {
	var validReq = httptest.NewRequest("GET", "http://localhost:8080/healthz", nil)

	type want struct {
		status int
	}
	tests := []struct {
		name    string
		request *http.Request
		want    want
	}{
		{
			name:    "should return status ok because valid liveness request",
			request: validReq,
			want: want{
				status: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Having a recorder for each separate request is safer!
			response := httptest.NewRecorder()

			l := Liveness{}

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
		})
	}
}

func TestReadiness_ServeHTTP_ServeHTTP(t *testing.T) {

	var validReq = httptest.NewRequest("GET", "http://localhost:8080/readyz", nil)

	type want struct {
		status int
	}
	tests := []struct {
		name    string
		fields  Readiness
		request *http.Request
		want    want
	}{
		{
			name: "should return bad status because DB not running yet",
			fields: Readiness{
				DB: nil,
			},
			request: validReq,
			want: want{
				status: http.StatusServiceUnavailable,
			},
		},
		{
			name: "should return bad status because DB not running yet",
			fields: Readiness{
				DB: mockDatabase(t),
			},
			request: validReq,
			want: want{
				status: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Having a recorder for each separate request is safer!
			response := httptest.NewRecorder()

			r := Readiness{
				DB: tt.fields.DB,
			}

			r.ServeHTTP(response, tt.request)

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

func mockDatabase(t *testing.T) *gorm.DB {
	sqlDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	// GORM allows to initialize *gorm.DB with an existing database connection https://gorm.io/docs/connecting_to_the_database.html
	//In this case it happens to be the mock db
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		t.Fatal(err)
	}
	return gormDB
}
