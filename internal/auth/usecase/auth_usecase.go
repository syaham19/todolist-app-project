package usecase

import (
	"errors"
	"todolist-project/internal/auth/model"
	"todolist-project/internal/auth/repository"

	"time"

	"golang.org/x/crypto/bcrypt"
)

// Error sentinel agar layer delivery bisa membedakan kasus
// "data tidak ditemukan" (404) dengan error internal (500).
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyExist = errors.New("email sudah terdaftar")
)

type AuthUseCase interface {
	Register(req *model.RegisterRequest) error
	GetAllUsers() ([]*model.User, error)
	GetUserByID(id int) (*model.User, error)
	UpdateUser(id int, req *model.UpdateUser) error
	DeleteUser(id int) error
	Login(req *model.LoginRequest) (*model.LoginResponse, error)
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
		return ErrEmailAlreadyExist
	}

	// 2. Hash password menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 3. Buat entity user baru
	now := time.Now()
	newUser := &model.User{
		CreatedAt: now,
		UpdatedAt: now,
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
	}

	// 4. Simpan user baru melalui repository
	return u.userRepo.CreateUser(newUser)
}

func (u *authUseCase) GetAllUsers() ([]*model.User, error) {
	return u.userRepo.GetAllUsers()
}

func (u *authUseCase) GetUserByID(id int) (*model.User, error) {
	return u.userRepo.GetUserByID(id)
}

// UpdateUser mengubah data user berdasarkan id.
// req hanya berisi field yang mau diubah, field yang kosong tidak menimpa data lama.
func (u *authUseCase) UpdateUser(id int, req *model.UpdateUser) error {
	// 1. Ambil user lama
	existing, err := u.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrUserNotFound
	}

	// 2. Salin data lama, lalu terapkan perubahan
	updated := *existing

	if req.Name != "" {
		updated.Name = req.Name
	}

	if req.Email != "" && req.Email != updated.Email {
		// Pastikan email baru tidak dipakai user lain
		other, err := u.userRepo.GetUserByEmail(req.Email)
		if err != nil {
			return err
		}
		if other != nil && other.ID != id {
			return ErrEmailAlreadyExist
		}
		updated.Email = req.Email
	}

	if req.Password != "" {
		// Password baru harus di-hash ulang
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		updated.Password = string(hashedPassword)
	}

	updated.UpdatedAt = time.Now()

	// 3. Simpan ke repository (repository menerima 1 argumen: *model.User)
	return u.userRepo.UpdateUser(&updated)
}

func (u *authUseCase) DeleteUser(id int) error {
	// Cek dulu user-nya ada atau tidak supaya bisa dibalas 404
	existing, err := u.userRepo.GetUserByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrUserNotFound
	}
	return u.userRepo.DeleteUser(id)
}

func (u *authUseCase) Login(req *model.LoginRequest) (*model.LoginResponse, error) {
	// 1. Ambil user berdasarkan email
	user, err := u.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("Email atau password salah")
	}

	// 2. Bandingkan password yang di-hash dengan password yang diberikan
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 3. Jika berhasil, kembalikan response login (misal token JWT)
	response := &model.LoginResponse{
		Message: "Login successful",
		// Token JWT bisa ditambahkan di sini nanti
	}

	return response, nil
}
