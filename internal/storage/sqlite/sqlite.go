package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nandinigthub/students-api/internal/models"
)

type sqlite struct {
	Db *sql.DB
}

// new student db
func New(cfg *models.Config) (*sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	email TEXT,
	age INTEGER )`)

	if err != nil {
		return nil, err
	}

	return &sqlite{
		Db: db,
	}, nil

}

// create student db

func (s *sqlite) CreateStudent(name string, email string, age int) (int64, error) {

	stmt, err := s.Db.Prepare("INSERT INTO students(name,email,age)VALUES (?,?,?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastid, err := result.LastInsertId()
	if err != nil {
		return 0, err // empty value and err
	}

	return lastid, nil
}
