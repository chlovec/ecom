package user

import (
	"database/sql"
	"ecom/types"
	"fmt"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec("INSERT INTO users(firstName, lastName, email, password) VALUES(?, ?, ?, ?)", user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	query := "SELECT * FROM users WHERE email = ?"
	rows, err := s.db.Query(query, email)
	if err != nil {
		return nil, err
	}

	return getUserFromRows(rows)
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	query := "SELECT * FROM users WHERE id = ?"
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	return getUserFromRows(rows)
}

func getUserFromRows(rows *sql.Rows) (*types.User, error) {
	u := new(types.User)
	for rows.Next() {
		var err error
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func scanRowsIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)

	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}