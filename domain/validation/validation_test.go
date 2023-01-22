package validation

import (
	"github.com/FirasYousfi/tasks-web-servcie/domain/entity"
	"reflect"
	"testing"
)

func TestValidateParams(t *testing.T) {
	type args struct {
		req *entity.TaskDescription
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.TaskDescription
		wantErr bool
	}{
		{
			name: "should succeed because of valid params",
			args: args{req: &entity.TaskDescription{
				Title:       "test",
				Description: "test desc",
				Priority:    5,
				Status:      "new",
			}},
			want: &entity.TaskDescription{
				Title:       "test",
				Description: "test desc",
				Priority:    5,
				Status:      "new",
			},
			wantErr: false,
		},
		{
			name: "should fail because of invalid empty title",
			args: args{req: &entity.TaskDescription{
				Title:       "",
				Description: "test desc",
				Priority:    5,
				Status:      "new",
			}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should fail because of invalid status",
			args: args{req: &entity.TaskDescription{
				Title:       "title",
				Description: "test desc",
				Priority:    5,
				Status:      "invalid",
			}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateParams(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ValidateParams() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateDescription(t *testing.T) {
	type args struct {
		description string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "should fail invalid description",
			args:    args{description: string(make([]rune, 505))}, //INFO this is how to make string from char array https://www.codegrepper.com/code-examples/go/go+make+string+from+char+array
			wantErr: true,
		},
		{
			name:    "should succeed because valid description",
			args:    args{description: "this is valid"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateDescription(tt.args.description); (err != nil) != tt.wantErr {
				t.Errorf("ValidateDescription() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePriority(t *testing.T) {
	type args struct {
		priority int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "should fail invalid priority",
			args:    args{priority: 500},
			wantErr: true,
		},
		{
			name:    "should succeed valid prioriy",
			args:    args{priority: 10},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidatePriority(tt.args.priority); (err != nil) != tt.wantErr {
				t.Errorf("ValidatePriority() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateTitle(t *testing.T) {
	type args struct {
		title string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "should fail invalid title",
			args:    args{title: ""},
			wantErr: true,
		},
		{
			name:    "should pays valid title",
			args:    args{title: "this is a valid title"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateTitle(tt.args.title); (err != nil) != tt.wantErr {
				t.Errorf("ValidateTitle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateStatus(t *testing.T) {
	type args struct {
		status entity.Status
	}
	tests := []struct {
		name    string
		args    args
		want    entity.Status
		wantErr bool
	}{
		{
			name:    "should pass because valid status",
			args:    args{status: entity.OnHold},
			want:    entity.OnHold,
			wantErr: false,
		},
		{
			name:    "should return error because invalid status",
			args:    args{status: "invalid-status"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateStatus(tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}
