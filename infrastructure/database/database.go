package database

import (
	"fmt"
	"github.com/FirasYousfi/tasks-web-servcie/config"
	"github.com/FirasYousfi/tasks-web-servcie/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

const (
	maxIdleConns = 10
	maxOpenConns = 100
)

// DB is the actual database connection instance, defined as interface to be implemented by any possible database. Also, it improves testability.
var DB DbConnector // Needs to be initialized before we can assign to it

// DbConnector is an interface defining the contract needed to create a new Database connection
type DbConnector interface {
	SetDBConn() error
	SetDBConf(dbConfig *config.DbConfig)
	GetDBConn() *gorm.DB //Because this avoids casting in main, since we cannot do DB.Conf because DB is actually an interface
	GetDBConf() *config.DbConfig
}

// Database is the struct holding the current Database configuration, it implements the DB connector struct
type Database struct {
	Conf *config.DbConfig
	Conn *gorm.DB
}

// InitializeDB Reads the user defined DB configuration and initializes the connection. If another DB struct implements the DbConnector interface, this InitializeDB function needs to be updated.
func InitializeDB() error {
	DB = &Database{} //otherwise DB would be nil
	DB.SetDBConf(&config.Config.DB)
	// Database used to persist tasks
	return DB.SetDBConn()

}

func (d *Database) SetDBConn() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", d.Conf.Host,
		d.Conf.User, d.Conf.Password, d.Conf.Name, d.Conf.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to Database: %w", err)
	}

	// AutoMigrate will create tables, missing foreign keys, constraints, columns and indexes. You can give it multiple structs.
	err = db.AutoMigrate(&entity.Task{})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(maxIdleConns)

	// SetMaxOpenConns sets the maximum number of open connections to the Database.
	sqlDB.SetMaxOpenConns(maxOpenConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	d.Conn = db

	return nil
}

// GetDBConf returns the actual gormDB connection that can be used by the repository out of the interface given as input in SetupHandlers in main
func (d *Database) GetDBConf() *config.DbConfig {
	return d.Conf
}

// GetDBConn returns the actual gormDB connection that can be used by the repository out of the interface given as input in SetupHandlers in main
func (d *Database) GetDBConn() *gorm.DB {
	return d.Conn
}

// SetDBConf defines the way to set the configuration for the specific database instance
func (d *Database) SetDBConf(dbConfig *config.DbConfig) {
	d.Conf = dbConfig
}
