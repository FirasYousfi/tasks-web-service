package repository

import (
	"database/sql/driver"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/FirasYousfi/tasks-web-servcie/application/interfaces"
	"github.com/FirasYousfi/tasks-web-servcie/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"reflect"
	"regexp"
	"testing"
	"time"
)

type Suite struct {
	gormDB     *gorm.DB
	mock       sqlmock.Sqlmock
	repository interfaces.ITaskRepository
}

// AnyTime is to be used for with args in sqlmock for timestamps
type AnyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func (s *Suite) SetupSuite() {
	sqlDB, mock, err := sqlmock.New()
	s.mock = mock
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	// GORM allows to initialize *gorm.DB with an existing database connection https://gorm.io/docs/connecting_to_the_database.html
	s.gormDB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	s.repository = NewTaskRepository(s.gormDB)
}

func TestNewTaskRepository(t *testing.T) {
	db := &gorm.DB{}
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want interfaces.ITaskRepository
	}{
		{
			name: "should create new task repository correctly",
			args: args{db: db},
			want: &TaskRepository{db: db},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTaskRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTaskRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskRepository_Create(t *testing.T) {
	var testSuite Suite
	testSuite.SetupSuite()
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		task *entity.Task
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "should work correctly and create element",
			fields: fields{db: testSuite.gormDB},
			args: args{task: &entity.Task{
				ID:        "1",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
				TaskDescription: entity.TaskDescription{
					Title:       "test",
					Description: "test",
					Priority:    5,
					Status:      "New",
				},
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskRepo := &TaskRepository{
				db: tt.fields.db,
			}
			testSuite.mock.ExpectBegin()
			testSuite.mock.ExpectExec(regexp.QuoteMeta(
				`INSERT INTO "tasks"`)).
				WithArgs(tt.args.task.ID, AnyTime{}, AnyTime{}, tt.args.task.Title, tt.args.task.Description, tt.args.task.Priority, tt.args.task.Status).
				WillReturnResult(sqlmock.NewResult(0, 0))
			testSuite.mock.ExpectCommit()

			if err := taskRepo.Create(tt.args.task); (err != nil) != tt.wantErr {
				fmt.Println("here is the error: ", err)
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTaskRepository_DeleteByID(t1 *testing.T) {
	var testSuite Suite
	testSuite.SetupSuite()
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "should pass and delete Item",
			fields:  fields{db: testSuite.gormDB},
			args:    args{id: "3"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TaskRepository{
				db: tt.fields.db,
			}

			testSuite.mock.ExpectBegin()
			testSuite.mock.ExpectExec(regexp.QuoteMeta(
				`DELETE FROM "tasks" WHERE id = $1`)).
				WithArgs(tt.args.id).
				WillReturnResult(sqlmock.NewResult(0, 0))
			testSuite.mock.ExpectCommit()

			if err := t.DeleteByID(tt.args.id); (err != nil) != tt.wantErr {
				t1.Errorf("DeleteByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTaskRepository_FindAll(t1 *testing.T) {
	var testSuite Suite
	testSuite.SetupSuite()
	type fields struct {
		db *gorm.DB
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*entity.Task
		wantErr bool
	}{
		{
			name:   "should return instances correctly",
			fields: fields{db: testSuite.gormDB},
			want: []*entity.Task{{
				ID: "1",
				TaskDescription: entity.TaskDescription{
					Title:       "test",
					Description: "test",
					Priority:    5,
					Status:      "New",
				},
			}, {
				ID: "2",
				TaskDescription: entity.TaskDescription{
					Title:       "test",
					Description: "test",
					Priority:    5,
					Status:      "New",
				},
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TaskRepository{
				db: tt.fields.db,
			}
			columns := []string{"id", "title", "description", "priority", "status"}
			rows := testSuite.mock.NewRows(columns).AddRow("1", "test", "test", 5, "New").
				AddRow("2", "test", "test", 5, "New")

			testSuite.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "tasks"`)).WillReturnRows(rows)

			got, err := t.FindAll()
			if (err != nil) != tt.wantErr {
				t1.Errorf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("FindAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskRepository_FindByID(t1 *testing.T) {
	var testSuite Suite
	testSuite.SetupSuite()
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Task
		wantErr bool
	}{
		{
			name:   "should return instance correctly",
			fields: fields{db: testSuite.gormDB},
			args:   args{id: "1234"},
			want: &entity.Task{
				ID:              "1234",
				CreatedAt:       time.Time{},
				UpdatedAt:       time.Time{},
				TaskDescription: entity.TaskDescription{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TaskRepository{
				db: tt.fields.db,
			}

			testSuite.mock.ExpectQuery(regexp.QuoteMeta(
				`SELECT * FROM "tasks" WHERE id = $1 ORDER BY "tasks"."id" LIMIT 1`)).
				WithArgs(tt.args.id).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).
					AddRow(tt.args.id))

			got, err := t.FindByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t1.Errorf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("FindByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskRepository_Update(t1 *testing.T) {
	var testSuite Suite
	testSuite.SetupSuite()

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		fields map[string]interface{}
		id     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "should pass and update entry",
			fields: fields{db: testSuite.gormDB},
			args: args{
				fields: map[string]interface{}{"description": "updated-description", "priority": 5, "status": "New"},
				id:     "1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TaskRepository{
				db: tt.fields.db,
			}
			testSuite.mock.ExpectBegin()
			testSuite.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "tasks" SET "description"=$1,"priority"=$2,"status"=$3,"updated_at"=$4 WHERE id = $5`)).
				WithArgs(tt.args.fields["description"], tt.args.fields["priority"], tt.args.fields["status"], AnyTime{}, tt.args.id).
				WillReturnResult(sqlmock.NewResult(1, 1))
			testSuite.mock.ExpectCommit()

			if err := t.Update(tt.args.fields, tt.args.id); (err != nil) != tt.wantErr {
				t1.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
