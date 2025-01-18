package user

import (
	"bytes"
	"ecom/types"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
	
	t.Run("should create new user if payload is valid and email does not exist", func(t *testing.T) {
		// Mock expected behavior
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

		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/user", handler.handleRegister).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, http.StatusCreated, rr.Code)
	})
}