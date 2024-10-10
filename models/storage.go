package models

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetstudentbyId(id int64) (Student, error)
	GetallStudent() ([]Student, error)
}
