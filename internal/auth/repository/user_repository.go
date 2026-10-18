package repository

import (
	"todolist-project/internal/auth/model"
)

type UserRepository interface {
	// CreateUser creates a new user in the database.
	CreateUser(user *model.User) error

	// GetUserByEmail retrieves a user by their email address.
	GetUserByEmail(email string) (*model.User, error)

	// GetUserByID retrieves a user by their ID.
	GetUserByID(id string) (*model.User, error)

	GetAllUsers() ([]*model.User, error)
}

type userRepositoryMem struct {
	users map[string]*model.User
}

// type newUserRepositoryMem struct {
// 	users map[string]*model.User
// }

func NewUserRepositoryMem() UserRepository {
	return &userRepositoryMem{
		users: make(map[string]*model.User),
	}
}

func (r *userRepositoryMem) CreateUser(user *model.User) error {
	r.users[user.Email] = user
	return nil
}

func (r *userRepositoryMem) GetUserByEmail(email string) (*model.User, error) {
	user, exists := r.users[email]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (r *userRepositoryMem) GetUserByID(id string) (*model.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, nil
}

func (r *userRepositoryMem) GetAllUsers() ([]*model.User, error) {
	var users []*model.User
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}
