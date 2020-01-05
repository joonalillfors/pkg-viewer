package controllers

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// Global variable for database connection
var Database *sql.DB
