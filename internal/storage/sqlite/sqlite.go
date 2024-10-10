package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nandinigthub/students-api/models"
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

// get student by id
func (s *sqlite) GetstudentbyId(id int64) (models.Student, error) {

	slog.Info("getting student by id")
	stmt, err := s.Db.Prepare("SELECT * FROM students WHERE id == ? LIMIT 1 ")
	if err != nil {
		return models.Student{}, err
	}

	defer stmt.Close()

	var stu models.Student

	err = stmt.QueryRow(id).Scan(&stu.Id, &stu.Name, &stu.Email, &stu.Age)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Student{}, fmt.Errorf("no student with this id %s", fmt.Sprint(id))
		}

		return models.Student{}, fmt.Errorf("query error: %w", err)
	}

	return stu, nil

}

// get all students
func (s *sqlite) GetallStudent() ([]models.Student, error) {
	slog.Info("getting all student")

	stmt, err := s.Db.Prepare("SELECT * FROM students")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var stu []models.Student

	for rows.Next() {
		var student models.Student

		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}

		stu = append(stu, student)

	}
	return stu, nil

}
