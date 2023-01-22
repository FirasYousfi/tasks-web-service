package service

import (
	"errors"
	"github.com/FirasYousfi/tasks-web-servcie/application/interfaces"
	"github.com/FirasYousfi/tasks-web-servcie/domain/entity"
	"reflect"
	"testing"
	"time"
)

const (
	testID              = "testID"
	testFullUpdateID    = "testFullUpdateID"
	testPartialUpdateID = "testPartialUpdateID"
)

var (
	TaskRequestInstance = entity.TaskDescription{
		Title:       "test",
		Description: "test",
		Priority:    5,
		Status:      "new",
	}
	FullUpdateRequest = entity.TaskDescription{
		Title:       "testFullUpdate",
		Description: "testFullUpdate",
		Priority:    1,
		Status:      "new",
	}
	PartialUpdateRequest = entity.TaskDescription{
		Title:  "testPartialUpdate",
		Status: "Active",
	}
)

type mockTaskRepository struct {
}

func (m mockTaskRepository) Create(task *entity.Task) error {
	if reflect.DeepEqual(task.TaskDescription, TaskRequestInstance) {
		return nil
	}
	return errors.New("unexpected values passed to the repository")
}

func (m mockTaskRepository) DeleteByID(id string) error {
	_, err := m.FindByID(id)
	return err
}

func (m mockTaskRepository) Update(fields map[string]interface{}, id string) error {
	fullUpdateValues := map[string]interface{}{"title": FullUpdateRequest.Title, "description": FullUpdateRequest.Description,
		"priority": FullUpdateRequest.Priority, "status": FullUpdateRequest.Status}

	partialUpdateValues := map[string]interface{}{"title": PartialUpdateRequest.Title, "status": PartialUpdateRequest.Status}
	//Task should be found
	_, err := m.FindByID(id)
	if err != nil {
		return err
	}
	//checking with id to test full-update functionality
	if id == testFullUpdateID && !reflect.DeepEqual(fields, fullUpdateValues) {
		return errors.New("not all fields are being updated correctly")
	}
	//checking with id to test partial-update functionality
	if id == testPartialUpdateID && !reflect.DeepEqual(fields, partialUpdateValues) {
		return errors.New("the specified fields are not being updated")
	}
	return nil
}

func (m mockTaskRepository) FindAll() ([]*entity.Task, error) {
	return []*entity.Task{{
		ID:              "test1",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
		TaskDescription: entity.TaskDescription{},
	}, {
		ID:              "test2",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
		TaskDescription: entity.TaskDescription{},
	}}, nil
}

func (m mockTaskRepository) FindByID(id string) (*entity.Task, error) {
	if id == testID {
		return &entity.Task{
			ID:              id,
			TaskDescription: TaskRequestInstance,
		}, nil
	} else if id == testFullUpdateID {
		return &entity.Task{
			ID:              testFullUpdateID,
			TaskDescription: FullUpdateRequest,
		}, nil
	} else if id == testPartialUpdateID {
		return &entity.Task{
			ID:              testPartialUpdateID,
			TaskDescription: PartialUpdateRequest,
		}, nil
	}
	return nil, errors.New("instance with specified ID not found")
}

func TestNewTaskService(t *testing.T) {
	type args struct {
		repo interfaces.ITaskRepository
	}
	tests := []struct {
		name string
		args args
		want interfaces.ITaskService
	}{
		{
			name: "should pass",
			args: args{mockTaskRepository{}},
			want: &TaskService{TaskRepository: mockTaskRepository{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTaskService(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTaskService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskService_Create(t1 *testing.T) {
	type fields struct {
		TaskRepository interfaces.ITaskRepository
	}
	type args struct {
		req *entity.TaskDescription
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.TaskDescription
		wantErr bool
	}{
		{
			name:    "should pass and create task",
			fields:  fields{TaskRepository: mockTaskRepository{}},
			args:    args{req: &TaskRequestInstance},
			want:    &TaskRequestInstance,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TaskService{
				TaskRepository: tt.fields.TaskRepository,
			}
			got, err := t.Create(tt.args.req)
			if (err != nil) != tt.wantErr {
				t1.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//the stuff that is added by gorm is already tested on repo side, here we only need to check if the request is propagated
			if !reflect.DeepEqual(got.TaskDescription, *tt.want) {
				t1.Errorf("Create() got = %v, want %v", got.TaskDescription, tt.want)
			}
		})
	}
}

func TestTaskService_DeleteByID(t1 *testing.T) {
	type fields struct {
		TaskRepository interfaces.ITaskRepository
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
			name:    "should pass because instance to delete is found",
			fields:  fields{TaskRepository: mockTaskRepository{}},
			args:    args{testID},
			wantErr: false,
		},
		{
			name:    "should fail because instance to delete is not found",
			fields:  fields{TaskRepository: mockTaskRepository{}},
			args:    args{"non-existing-ID"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TaskService{
				TaskRepository: tt.fields.TaskRepository,
			}
			if err := t.DeleteByID(tt.args.id); (err != nil) != tt.wantErr {
				t1.Errorf("DeleteByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTaskService_Get(t1 *testing.T) {
	type fields struct {
		TaskRepository interfaces.ITaskRepository
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*entity.Task
		wantErr bool
	}{
		{
			name:   "should pass and return the instances fetched by the repository",
			fields: fields{TaskRepository: mockTaskRepository{}},
			want: []*entity.Task{{
				ID:              "test1",
				CreatedAt:       time.Time{},
				UpdatedAt:       time.Time{},
				TaskDescription: entity.TaskDescription{},
			}, {
				ID:              "test2",
				CreatedAt:       time.Time{},
				UpdatedAt:       time.Time{},
				TaskDescription: entity.TaskDescription{},
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TaskService{
				TaskRepository: tt.fields.TaskRepository,
			}
			got, err := t.Get()
			if (err != nil) != tt.wantErr {
				t1.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskService_GetByID(t1 *testing.T) {
	type fields struct {
		TaskRepository interfaces.ITaskRepository
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
			name:   "should pass and return task with specified ID",
			fields: fields{TaskRepository: mockTaskRepository{}},
			args:   args{testID},
			want: &entity.Task{
				ID:              testID,
				CreatedAt:       time.Time{},
				UpdatedAt:       time.Time{},
				TaskDescription: TaskRequestInstance,
			},
			wantErr: false,
		},
		{
			name:    "should fail because there is no task with specified ID",
			fields:  fields{TaskRepository: mockTaskRepository{}},
			args:    args{"invalidID"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TaskService{
				TaskRepository: tt.fields.TaskRepository,
			}
			got, err := t.GetByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t1.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("GetByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskService_UpdateFully(t1 *testing.T) {
	type fields struct {
		TaskRepository interfaces.ITaskRepository
	}
	type args struct {
		req *entity.TaskDescription
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Task
		wantErr bool
	}{
		{
			name:   "should update instance fully",
			fields: fields{TaskRepository: mockTaskRepository{}},
			args: args{
				req: &FullUpdateRequest,
				id:  testFullUpdateID,
			},
			want: &entity.Task{
				ID:              testFullUpdateID,
				TaskDescription: FullUpdateRequest,
			},
			wantErr: false,
		},
		{
			name:   "should fail to update instance since it is not found",
			fields: fields{TaskRepository: mockTaskRepository{}},
			args: args{
				req: &FullUpdateRequest,
				id:  "invalid-id",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TaskService{
				TaskRepository: tt.fields.TaskRepository,
			}
			got, err := t.UpdateFully(tt.args.req, tt.args.id)
			if (err != nil) != tt.wantErr {
				t1.Errorf("UpdateFully() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("UpdateFully() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTaskService_UpdatePartial(t1 *testing.T) {
	type fields struct {
		TaskRepository interfaces.ITaskRepository
	}
	type args struct {
		req *entity.TaskDescription
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Task
		wantErr bool
	}{
		{
			name:   "should update instance partially",
			fields: fields{TaskRepository: mockTaskRepository{}},
			args: args{
				req: &PartialUpdateRequest,
				id:  testPartialUpdateID,
			},
			want: &entity.Task{
				ID:              testPartialUpdateID,
				TaskDescription: PartialUpdateRequest,
			},
			wantErr: false,
		},
		{
			name:   "should fail to update instance since it is not found",
			fields: fields{TaskRepository: mockTaskRepository{}},
			args: args{
				req: &PartialUpdateRequest,
				id:  "invalid-id",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &TaskService{
				TaskRepository: tt.fields.TaskRepository,
			}
			got, err := t.UpdatePartial(tt.args.req, tt.args.id)
			if (err != nil) != tt.wantErr {
				t1.Errorf("UpdatePartial() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("UpdatePartial() got = %v, want %v", got, tt.want)
			}
		})
	}
}
