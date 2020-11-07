package repositories

import (
	"errors"

	"github.com/baybaraandrey/ws_notification/internal/notification/entities"
	database "github.com/baybaraandrey/ws_notification/pkg/postgresql"
)

const (
	UserTable = "users_user"
)

type UserRepository interface {
	GetUserByUsername(username string) (*entities.User, error)
}

// NewUserRepository creates new user repository
func NewUserRepository() UserRepository {
	return &userRepository{}
}

type userRepository struct{}

func (u *userRepository) GetUserByUsername(username string) (*entities.User, error) {
	var user entities.User
	db := database.OpenConnection()

	rows, err := db.Query("SELECT id,username,password FROM users_user WHERE username=$1",
		username)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		rows.Scan(&user.ID, &user.Username, &user.Password)
	}

	if user.ID <= 0 {
		return nil, errors.New("User not found")
	}

	defer db.Close()
	defer rows.Close()

	return &user, nil
}
