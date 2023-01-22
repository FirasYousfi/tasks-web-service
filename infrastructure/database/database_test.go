package database

import (
	"fmt"
	"github.com/FirasYousfi/tasks-web-servcie/config"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

type mockDb struct {
	conf *config.DbConfig
	conn *gorm.DB
}

func newMockDb() *mockDb {
	return &mockDb{}
}

func (m *mockDb) SetDBConn() error {
	//here is how to go through values of a struct https://stackoverflow.com/questions/18926303/iterate-through-the-fields-of-a-struct-in-go
	v := reflect.ValueOf(*m.conf)

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).IsZero() {
			return fmt.Errorf("not all conf is initialized")
		}
	}
	return nil
}

func (m *mockDb) SetDBConf(dbConfig *config.DbConfig) {
	m.conf = dbConfig
}

func (m *mockDb) GetDBConn() *gorm.DB {
	return m.conn
}

func (m *mockDb) GetDBConf() *config.DbConfig {
	return m.conf
}

func TestDatabase_InitializeDB(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "should initialize DB instance with default values",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DB = newMockDb()
			config.BuildConfig()
			if err := InitializeDB(); (err != nil) != tt.wantErr {
				t.Errorf("InitializeDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
