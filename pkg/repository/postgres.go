package repository

import (
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

const (
	usersTable      = "users"
	todoListsTable  = "todo_lists"
	usersListsTable = "users_lists"
	todoItemsTable  = "todo_items"
	listItemsTable  = "list_items"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	dataSourceName := os.Getenv("DATABASE_URL")
	if dataSourceName == "" {
		dataSourceName = viper.GetString("DATABASE_URL")
	}

	db, err := sqlx.Open("postgres", dataSourceName) // fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
