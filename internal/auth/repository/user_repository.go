package repository

import (
	"errors"
	"todolist-project/internal/auth/model"
)

type UserRepository interface {
	// CreateUser creates a new user in the database.
	CreateUser(user *model.User) error

	// GetUserByEmail retrieves a user by their email address.
	GetUserByEmail(email string) (*model.User, error)

	// GetUserByID retrieves a user by their ID.
	GetUserByID(id int) (*model.User, error)

	GetAllUsers() ([]*model.User, error)

	UpdateUser(user *model.User) error

	DeleteUser(id int) error
}

type userRepositoryMem struct {
	users  map[string]*model.User
	nextID int
}

// type newUserRepositoryMem struct {
// 	users map[string]*model.User
// }

func NewUserRepositoryMem() UserRepository {
	return &userRepositoryMem{
		users:  make(map[string]*model.User),
		nextID: 1,
	}
}

func (r *userRepositoryMem) CreateUser(user *model.User) error {
	// ID naik monoton, tidak bentrok walaupun ada user yang dihapus
	user.ID = r.nextID
	r.nextID++
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

func (r *userRepositoryMem) GetUserByID(id int) (*model.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}
	// User tidak ditemukan -> kembalikan (nil, nil).
	// Layer usecase/delivery yang menentukan status 404.
	return nil, nil
}

func (r *userRepositoryMem) GetAllUsers() ([]*model.User, error) {
	var users []*model.User
	for _, user := range r.users {
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepositoryMem) UpdateUser(user *model.User) error {
	// Cari key lama (email lama) berdasarkan ID user
	oldEmail, exists := "", false
	for email, stored := range r.users {
		if stored.ID == user.ID {
			oldEmail, exists = email, true
			break
		}
	}
	if !exists {
		return errors.New("user not found")
	}

	// Kalau email berubah, hapus key lama supaya tidak ada data duplikat
	if oldEmail != user.Email {
		delete(r.users, oldEmail)
	}
	r.users[user.Email] = user
	return nil
}

func (r *userRepositoryMem) DeleteUser(id int) error {
	for email, user := range r.users {
		if user.ID == id {
			delete(r.users, email)
			return nil
		}
	}
	return errors.New("user not found")
}
