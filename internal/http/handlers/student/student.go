package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	// "strings"

	"github.com/go-playground/validator"
	"github.com/nandinigthub/students-api/internal/utils/response"
	"github.com/nandinigthub/students-api/models"
)

var stu models.Student

// home page
func Home() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Home page of student api")
		err := response.WriteJson(w, http.StatusOK, stu)
		if err != nil {
			response.WriteJson(w, http.StatusBadGateway, response.ErrorMessage(err))
			return
		}
	}
}

// adding new student
func New(s models.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating student")

		err := json.NewDecoder(r.Body).Decode(&stu)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.ErrorMessage(err))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadGateway, response.ErrorMessage(err))
			return
		}

		// validating request

		if err := validator.New().Struct(stu); err != nil {
			validateErr := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErr))
			return
		}

		lastId, err := s.CreateStudent(
			stu.Name,
			stu.Email,
			stu.Age,
		)

		slog.Info("user created successfully with id = ", slog.String("userid", fmt.Sprint(lastId)))

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}

// getting student by id
func GetstudentbyId(s models.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		// parts := strings.Split(r.URL.Path, "/")

		// // Check if the URL contains at least 3 parts (like /students/{id})
		// if len(parts) < 3 {
		// 	http.Error(w, "Invalid request path", http.StatusBadRequest)
		// 	return
		// }

		// // Extract the id from the last part of the URL
		// id := parts[len(parts)-1]
		slog.Info("getting student by", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ErrorMessage(err))
			return
		}

		student, err := s.GetstudentbyId(intId)
		if err != nil {
			slog.Error("getting err", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.ErrorMessage(err))
			return
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}
