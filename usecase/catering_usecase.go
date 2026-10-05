package usecase

import (
	"errors"
	"regexp"
	"strings"

	"github.com/seanalden-great/kecilung-resto-be/domain"
)

type cateringUsecase struct {
	repo domain.CateringRepository
}

func NewCateringUsecase(r domain.CateringRepository) domain.CateringUsecase {
	return &cateringUsecase{repo: r}
}

// === Helper Pembentuk Slug ===
func generateCateringSlug(title string) string {
	str := strings.ToLower(title)
	re := regexp.MustCompile("[^a-z0-9]+")
	str = re.ReplaceAllString(str, "-")
	return strings.Trim(str, "-")
}

func (u *cateringUsecase) GetGreeting() (domain.GreetingMessage, error) { return u.repo.GetGreeting() }
func (u *cateringUsecase) UpdateGreeting(msg *domain.GreetingMessage) error { return u.repo.UpdateGreeting(msg) }
func (u *cateringUsecase) GetAllCaterings() ([]domain.Catering, error) { return u.repo.FetchAllCaterings() }
func (u *cateringUsecase) GetCateringByID(id uint) (domain.Catering, error) { return u.repo.FindCateringByID(id) }
// === TAMBAHAN BARU: GetCateringBySlug ===
func (u *cateringUsecase) GetCateringBySlug(slug string) (domain.Catering, error) {
	return u.repo.FindCateringBySlug(slug)
}
func (u *cateringUsecase) CreateCatering(c *domain.Catering) error { return u.repo.CreateCatering(c) }

// === TAMBAHAN BARU ===
func (u *cateringUsecase) UpdateCatering(id uint, c *domain.Catering) error {
	c.ID = id
	// KEMUNGKINAN BESAR BARIS INI TERLEWAT OLEH ANDA SEBELUMNYA:
	c.Slug = generateCateringSlug(c.Name)
	return u.repo.UpdateCatering(c)
}

func (u *cateringUsecase) DeleteCatering(id uint) error {
	return u.repo.DeleteCatering(id)
}
// =====================

func (u *cateringUsecase) GetAllBookings() ([]domain.Booking, error) { return u.repo.FetchBookings() }

// func (u *cateringUsecase) CreateBooking(b *domain.Booking) error {
// 	// CEK BENTROK: Apakah tanggal ini sudah ada yang APPROVED?
// 	exists, err := u.repo.CheckApprovedBookingExists(b.BookingDate)
// 	if err != nil { return err }
// 	if exists {
// 		return errors.New("Mohon maaf, jadwal pada tanggal tersebut sudah penuh dipesan.")
// 	}
// 	b.Status = "PENDING"
// 	return u.repo.CreateBooking(b)
// }

func (u *cateringUsecase) CreateBooking(b *domain.Booking) error {
	// Pastikan waktu selesai > waktu mulai
	if !b.BookingEndDate.After(b.BookingDate) {
		return errors.New("waktu selesai harus lebih besar dari waktu mulai")
	}

	exists, err := u.repo.CheckTimeConflict(b.BookingDate, b.BookingEndDate)
	if err != nil { return err }
	if exists {
		return errors.New("Mohon maaf, jadwal pada waktu tersebut bertabrakan dengan pesanan lain.")
	}
	b.Status = "PENDING"
	return u.repo.CreateBooking(b)
}

func (u *cateringUsecase) ApproveBooking(id uint) error {
	// Jika admin meng-approve, bisa saja ditambahkan validasi ulang di sini
	return u.repo.UpdateBookingStatus(id, "APPROVED")
}

func (u *cateringUsecase) RejectBooking(id uint) error {
	return u.repo.UpdateBookingStatus(id, "REJECTED")
}

// Tambahkan fungsi ini
func (u *cateringUsecase) GetApprovedBookingsByCateringID(id uint) ([]domain.Booking, error) {
	return u.repo.FetchApprovedBookingsByCateringID(id)
}