package domain

import "time"

type Admin struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Username  string    `json:"username" gorm:"unique"`
	Email        string    `json:"email" gorm:"unique"` // BARU: Target pengiriman OTP
	Password  string    `json:"password"` // Akan menyimpan Hashed Password
	ImageURL  string    `json:"image_url"`
	ResetOTP     string    `json:"-"`                   // BARU: OTP disembunyikan dari JSON response
	OTPExpiredAt time.Time `json:"-"`                   // BARU: Batas waktu OTP
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthRepository interface {
	GetAdminByUsername(username string) (Admin, error)
	GetAdminByEmail(email string) (Admin, error) // BARU
	GetAdminByID(id uint) (Admin, error)
	UpdateAdmin(id uint, data *Admin) error
	CreateDefaultAdmin(admin *Admin) error // Opsional untuk Inisialisasi awal
}

type AuthUsecase interface {
	Login(username, password string) (string, Admin, error) // Returns Token JWT
	GetProfile(id uint) (Admin, error)
	UpdateProfile(id uint, data *Admin) error
	// TAMBAHAN FITUR FORGOT PASSWORD
	RequestOTP(email string) error
	VerifyOTP(email, otp string) error
	ResetPassword(email, newPassword string) error
}