package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/nandinigthub/students-api/internal/models"
	"github.com/nandinigthub/students-api/internal/utils/response"
)

var stu models.Student

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

		slog.Info("user created successfully with id=", slog.String("userid", fmt.Sprint(lastId)))

		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]int64{"id": lastId})
	}
}
