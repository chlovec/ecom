package user

import (
	"bytes"
	"ecom/types"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestUserRegistration(t *testing.T){
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userStore := NewStore(db)
	handler := NewHandler(userStore)
	
	t.Run("should create new user if email does not exist", func(t *testing.T) {
		user := types.RegisterUserPayload{
			FirstName: "Alice",
			LastName:  "Doe",
			Email:     "alice.doe@example.com",
			Password:  "password",
		}

		mock.ExpectQuery("SELECT \\* FROM users WHERE email = \\?").
		WithArgs("alice.doe@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "firstName", "lastName", "email", "password", "createdAt"})) // No rows returned

		mock.ExpectExec("INSERT INTO users").
		WithArgs("Alice", "Doe", "alice.doe@example.com", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Mock successful insertion

		body, _ := json.Marshal(user)
		req, err := http.NewRequest(http.MethodPost, "/user", bytes.NewReader(body))

		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/user", handler.handleRegister).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, http.StatusCreated, rr.Code)
	})

	t.Run("should fail if email exists", func(t *testing.T) {
		user := types.RegisterUserPayload{
			FirstName: "Alice",
			LastName:  "Doe",
			Email:     "alice.doe@example.com",
			Password:  "password",
		}

		// Define expected rows
		createdAt := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		rows := sqlmock.NewRows([]string{"id", "firstName", "lastName", "email", "password", "createdAt"}).
			AddRow(1, "Alice", "Doe", "alice.doe@example.com", "password123", createdAt)
		mock.ExpectQuery("SELECT \\* FROM users WHERE email = \\?").
		WithArgs("alice.doe@example.com").
		WillReturnRows(rows) // No rows returned

		body, _ := json.Marshal(user)
		req, err := http.NewRequest(http.MethodPost, "/user", bytes.NewReader(body))

		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/user", handler.handleRegister).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		expectedResponse := `{"error":"user with email alice.doe@example.com already exist"}`
		assert.JSONEq(t, expectedResponse, rr.Body.String())
	})

	t.Run("should fail if payload is invalid", func(t *testing.T) {
		user := types.RegisterUserPayload{
			FirstName: "",
			LastName:  "",
			Email:     "",
			Password:  "",
		}

		body, _ := json.Marshal(user)
		req, err := http.NewRequest(http.MethodPost, "/user", bytes.NewReader(body))

		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/user", handler.handleRegister).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		expectedResponse := `{
			"error": "validation error",
			"details": [
				{
					"field": "FirstName",
					"message": "Field validation failed on the 'required' tag"
				},
				{
					"field": "LastName",
					"message": "Field validation failed on the 'required' tag"
				},
				{
					"field": "Email",
					"message": "Field validation failed on the 'required' tag"
				},
				{
					"field": "Password",
					"message": "Field validation failed on the 'required' tag"
				}
			]
		}`		
		assert.JSONEq(t, expectedResponse, rr.Body.String())
	})
}