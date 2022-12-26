package config

import "testing"

func TestGetEnv(t *testing.T) {
	t.Setenv("TEST", "testValue")
	type args struct {
		key      string
		fallback string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should return testValue because env variable is set",
			args: args{
				key:      "TEST",
				fallback: "FallBackValue",
			},
			want: "testValue",
		},
		{
			name: "should return fallBackValue because env variable is not set",
			args: args{
				key:      "NonSetVariable",
				fallback: "FallBackValue",
			},
			want: "FallBackValue",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEnv(tt.args.key, tt.args.fallback); got != tt.want {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
