package usecase

import (
	"errors"
	"todolist-project/internal/auth/model"
	"todolist-project/internal/auth/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Register(req *model.RegisterRequest) error
	GetAllUsers() ([]*model.User, error)
}

type authUseCase struct {
	userRepo repository.UserRepository
}

func NewAuthUseCase(userRepo repository.UserRepository) AuthUseCase {
	return &authUseCase{
		userRepo: userRepo,
	}
}

func (u *authUseCase) Register(req *model.RegisterRequest) error {
	// 1. Cek apakah email sudah terdaftar di repository
	existingUser, err := u.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("email sudah terdaftar")
	}

	// 2. Hash password menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 3. Buat entity user baru
	newUser := &model.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		// Tambahkan field lain jika ada (misal Name atau ID/UUID)
	}

	// 4. Simpan user baru melalui repository
	return u.userRepo.CreateUser(newUser)
}

func (u *authUseCase) GetAllUsers() ([]*model.User, error) {
	return u.userRepo.GetAllUsers()
}
