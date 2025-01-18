package user

import (
	"bytes"
	"ecom/types"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestUserServiceHandlers(t *testing.T) {
	userStore := newMockUserStore()
	handler := NewHandler(userStore)

	t.Run("should pass if user data is valid", func(t *testing.T) {
		user := types.RegisterUserPayload{
			FirstName: "Alice",
			LastName:  "Doe",
			Email:     "alice.doe@example.com",
			Password:  "password",
		}
		body, _ := json.Marshal(user)
		req, err := http.NewRequest(http.MethodPost, "/user", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/user", handler.handleRegister).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d\nresponse: %s", http.StatusCreated, rr.Code, rr.Body.String())
		}
	})

	t.Run("should fail if user is not in request body", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/user", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/user", handler.handleRegister).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})
}

type mockUserStore struct {
	users map[string]types.User // Simulate a data store
}

func newMockUserStore() *mockUserStore {
	return &mockUserStore{users: make(map[string]types.User)}
}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	if user, exists := m.users[email]; exists {
		return &user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *mockUserStore) CreateUser(u types.User) error {
	if _, exists := m.users[u.Email]; exists {
		return fmt.Errorf("user already exists")
	}
	m.users[u.Email] = u
	return nil
}
