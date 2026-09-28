package usecase

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"github.com/seanalden-great/kecilung-resto-be/utils"
	"golang.org/x/crypto/bcrypt"
)

// GANTI INI DENGAN SECRET KEY YANG LEBIH AMAN (Bisa ditaruh di .env)
var jwtSecretKey = []byte("KecilungRestoSecretKey2026")

type authUsecase struct {
	repo domain.AuthRepository
}

func NewAuthUsecase(r domain.AuthRepository) domain.AuthUsecase {
	return &authUsecase{repo: r}
}

func (u *authUsecase) Login(username, password string) (string, domain.Admin, error) {
	admin, err := u.repo.GetAdminByUsername(username)
	if err != nil {
		return "", admin, errors.New("username tidak ditemukan")
	}

	// Bandingkan password input dengan password hash di DB
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err != nil {
		return "", admin, errors.New("password salah")
	}

	// Generate JWT Token (Berlaku 24 Jam)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   admin.ID,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecretKey)
	
	// Kosongkan password agar tidak terkirim kembali di JSON response
	admin.Password = "" 
	
	return tokenString, admin, err
}

func (u *authUsecase) GetProfile(id uint) (domain.Admin, error) {
	admin, err := u.repo.GetAdminByID(id)
	admin.Password = "" // Keamanan
	return admin, err
}

func (u *authUsecase) UpdateProfile(id uint, data *domain.Admin) error {
	// Jika ada password baru, Hash terlebih dahulu
	if data.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		data.Password = string(hashedPassword)
	}
	return u.repo.UpdateAdmin(id, data)
}

func (u *authUsecase) RequestOTP(email string) error {
	admin, err := u.repo.GetAdminByEmail(email)
	if err != nil {
		return errors.New("email tidak terdaftar")
	}

	otpCode := utils.GenerateOTP()
	expiredAt := time.Now().Add(10 * time.Minute) // Berlaku 10 Menit

	admin.ResetOTP = otpCode
	admin.OTPExpiredAt = expiredAt
	
	// Simpan OTP ke DB
	if err := u.repo.UpdateAdmin(admin.ID, &admin); err != nil {
		return errors.New("gagal menyimpan OTP")
	}

	// Kirim Email
	if err := utils.SendOTPEmail(email, otpCode); err != nil {
		return errors.New("gagal mengirim email, periksa konfigurasi SMTP")
	}

	return nil
}

func (u *authUsecase) VerifyOTP(email, otp string) error {
	admin, err := u.repo.GetAdminByEmail(email)
	if err != nil {
		return errors.New("email tidak valid")
	}

	if admin.ResetOTP != otp {
		return errors.New("kode OTP salah")
	}

	if time.Now().After(admin.OTPExpiredAt) {
		return errors.New("kode OTP telah kedaluwarsa")
	}

	return nil
}

func (u *authUsecase) ResetPassword(email, newPassword string) error {
	admin, err := u.repo.GetAdminByEmail(email)
	if err != nil {
		return errors.New("email tidak valid")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal mengamankan password")
	}

	admin.Password = string(hashedPassword)
	admin.ResetOTP = "" // Hapus OTP setelah sukses reset
	return u.repo.UpdateAdmin(admin.ID, &admin)
}