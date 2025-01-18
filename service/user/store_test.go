package user

import (
	"ecom/types"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := NewStore(db)

	// Define expected behavior
	mock.ExpectExec("INSERT INTO users").
		WithArgs("John", "Doe", "john.doe@example.com", "password123").
		WillReturnResult(sqlmock.NewResult(1, 1)) // Mock successful insertion

	// Call the method
	err = store.CreateUser(types.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "password123",
	})

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail(t *testing.T) {
	// Mock the database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := NewStore(db)

	t.Run("should return user when email exists", func(t *testing.T) {
		// Define expected rows
		createdAt := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		rows := sqlmock.NewRows([]string{"id", "firstName", "lastName", "email", "password", "createdAt"}).
			AddRow(1, "John", "Doe", "john.doe@example.com", "password123", createdAt)

		// Mock query
		mock.ExpectQuery("SELECT \\* FROM users WHERE email = \\?").
			WithArgs("john.doe@example.com").
			WillReturnRows(rows)

		// Call the method
		user, err := store.GetUserByEmail("john.doe@example.com")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, 1, user.ID)
		assert.Equal(t, "John", user.FirstName)
		assert.Equal(t, "Doe", user.LastName)
		assert.Equal(t, "john.doe@example.com", user.Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when email does not exist", func(t *testing.T) {
		// Mock query
		mock.ExpectQuery("SELECT \\* FROM users WHERE email = \\?").
			WithArgs("john.doe@example.com").
			WillReturnRows(sqlmock.NewRows([]string{"id", "firstName", "lastName", "email", "password", "createdAt"}))

		// Call the method
		user, err := store.GetUserByEmail("john.doe@example.com")

		// Assert
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.EqualError(t, err, "user not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	
}

func TestGetUserByID(t *testing.T) {
	// Mock the database
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := NewStore(db)

	t.Run("should return user when id exists", func(t *testing.T) {
		// Define expected rows
		createdAt := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		rows := sqlmock.NewRows([]string{"id", "firstName", "lastName", "email", "password", "createdAt"}).
			AddRow(1, "John", "Doe", "john.doe@example.com", "password123", createdAt)

		// Mock query
		mock.ExpectQuery("SELECT \\* FROM users WHERE id = \\?").
			WithArgs(1).
			WillReturnRows(rows)

		// Call the method
		user, err := store.GetUserByID(1)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, 1, user.ID)
		assert.Equal(t, "John", user.FirstName)
		assert.Equal(t, "Doe", user.LastName)
		assert.Equal(t, "john.doe@example.com", user.Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when id does not exist", func(t *testing.T) {
		// Mock query with no rows returned
		mock.ExpectQuery("SELECT \\* FROM users WHERE id = \\?").
			WithArgs(99). // Simulate a non-existing user ID
			WillReturnRows(sqlmock.NewRows([]string{"id", "firstName", "lastName", "email", "password", "createdAt"})) // No rows

		// Call the method
		user, err := store.GetUserByID(99)

		// Assert
		assert.Nil(t, user)
		assert.Error(t, err)
		assert.EqualError(t, err, "user not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}