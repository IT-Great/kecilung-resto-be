package repository

import (
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"gorm.io/gorm"
)

type authRepo struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) domain.AuthRepository {
	return &authRepo{db}
}

func (r *authRepo) GetAdminByUsername(username string) (domain.Admin, error) {
	var admin domain.Admin
	err := r.db.Where("username = ?", username).First(&admin).Error
	return admin, err
}

func (r *authRepo) GetAdminByID(id uint) (domain.Admin, error) {
	var admin domain.Admin
	err := r.db.First(&admin, id).Error
	return admin, err
}

func (r *authRepo) UpdateAdmin(id uint, data *domain.Admin) error {
	return r.db.Model(&domain.Admin{}).Where("id = ?", id).Updates(data).Error
}

func (r *authRepo) CreateDefaultAdmin(admin *domain.Admin) error {
	var count int64
	r.db.Model(&domain.Admin{}).Count(&count)
	if count == 0 {
		return r.db.Create(admin).Error
	}
	return nil // Jangan buat lagi jika sudah ada admin
}

// Tambahkan fungsi ini di dalam file repository/auth_repository.go
func (r *authRepo) GetAdminByEmail(email string) (domain.Admin, error) {
	var admin domain.Admin
	err := r.db.Where("email = ?", email).First(&admin).Error
	return admin, err
}