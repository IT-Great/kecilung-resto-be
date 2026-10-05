package usecase

import (
	"errors"
	"regexp"
	"strings"

	"github.com/seanalden-great/kecilung-resto-be/domain"
)

type momentUsecase struct {
	repo domain.MomentRepository
}

func NewMomentUsecase(r domain.MomentRepository) domain.MomentUsecase {
	return &momentUsecase{repo: r}
}

// === Helper Pembentuk Slug ===
func generateMomentSlug(title string) string {
	str := strings.ToLower(title)
	re := regexp.MustCompile("[^a-z0-9]+")
	str = re.ReplaceAllString(str, "-")
	return strings.Trim(str, "-")
}

func (u *momentUsecase) GetAllMoments() ([]domain.Moment, error) { return u.repo.FetchAllMoments() }
func (u *momentUsecase) GetMomentByID(id uint) (domain.Moment, error) { return u.repo.FindMomentByID(id) }
// === TAMBAHAN BARU: GetMomentBySlug ===
func (u *momentUsecase) GetMomentBySlug(slug string) (domain.Moment, error) {
	return u.repo.FindMomentBySlug(slug)
}
func (u *momentUsecase) CreateMoment(m *domain.Moment) error { 
	m.Slug = generateMomentSlug(m.Name)
	return u.repo.CreateMoment(m) 
}

func (u *momentUsecase) UpdateMoment(id uint, m *domain.Moment) error {
	m.ID = id
	m.Slug = generateMomentSlug(m.Name)
	return u.repo.UpdateMoment(m)
}

func (u *momentUsecase) DeleteMoment(id uint) error { return u.repo.DeleteMoment(id) }

func (u *momentUsecase) GetAllBookings() ([]domain.MomentBooking, error) { return u.repo.FetchBookings() }

// func (u *momentUsecase) CreateBooking(b *domain.MomentBooking) error {
// 	exists, err := u.repo.CheckApprovedBookingExists(b.BookingDate)
// 	if err != nil { return err }
// 	if exists {
// 		return errors.New("Mohon maaf, jadwal pada tanggal tersebut sudah penuh dipesan.")
// 	}
// 	b.Status = "PENDING"
// 	return u.repo.CreateBooking(b)
// }

func (u *momentUsecase) CreateBooking(b *domain.MomentBooking) error {
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

func (u *momentUsecase) ApproveBooking(id uint) error { return u.repo.UpdateBookingStatus(id, "APPROVED") }
func (u *momentUsecase) RejectBooking(id uint) error { return u.repo.UpdateBookingStatus(id, "REJECTED") }

// Tambahkan fungsi ini
func (u *momentUsecase) GetApprovedBookingsByMomentID(id uint) ([]domain.MomentBooking, error) {
	return u.repo.FetchApprovedBookingsByMomentID(id)
}