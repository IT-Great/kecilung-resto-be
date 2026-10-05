package repository

import (
	"time"
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"gorm.io/gorm"
)

type momentRepo struct {
	db *gorm.DB
}

func NewMomentRepository(db *gorm.DB) domain.MomentRepository {
	return &momentRepo{db}
}

func (r *momentRepo) FetchAllMoments() ([]domain.Moment, error) {
	var moments []domain.Moment
	err := r.db.Preload("Images").Find(&moments).Error
	return moments, err
}

func (r *momentRepo) FindMomentByID(id uint) (domain.Moment, error) {
	var moment domain.Moment
	err := r.db.Preload("Images").First(&moment, id).Error
	return moment, err
}

// === TAMBAHAN BARU: Cari via Slug ===
func (r *momentRepo) FindMomentBySlug(slug string) (domain.Moment, error) {
	var moment domain.Moment
	err := r.db.Preload("Images").Where("slug = ?", slug).First(&moment).Error
	return moment, err
}

func (r *momentRepo) CreateMoment(m *domain.Moment) error {
	return r.db.Create(m).Error
}

// func (r *momentRepo) UpdateMoment(m *domain.Moment) error {
// 	return r.db.Transaction(func(tx *gorm.DB) error {
// 		if err := tx.Model(&domain.Moment{}).Where("id = ?", m.ID).Updates(map[string]interface{}{
// 			"name":        m.Name,
// 			"description": m.Description,
// 		}).Error; err != nil {
// 			return err
// 		}
// 		if len(m.Images) > 0 {
// 			if err := tx.Where("moment_id = ?", m.ID).Delete(&domain.MomentImage{}).Error; err != nil {
// 				return err
// 			}
// 			for i := range m.Images {
// 				m.Images[i].MomentID = m.ID
// 			}
// 			if err := tx.Create(&m.Images).Error; err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	})
// }

// === UBAH: Fungsi Update ===
func (r *momentRepo) UpdateMoment(m *domain.Moment) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existing domain.Moment
		if err := tx.First(&existing, m.ID).Error; err != nil {
			return err
		}

		// Update kolom dasar (secara eksplisit, termasuk slug)
		existing.Name = m.Name
		existing.Description = m.Description
		existing.Slug = m.Slug // <-- PASTI TERSIMPAN

		if err := tx.Save(&existing).Error; err != nil {
			return err
		}

		// Manajemen Gambar (hapus lama, simpan baru)
		if len(m.Images) > 0 {
			if err := tx.Where("moment_id = ?", m.ID).Delete(&domain.MomentImage{}).Error; err != nil {
				return err
			}
			for i := range m.Images {
				m.Images[i].MomentID = m.ID
			}
			if err := tx.Create(&m.Images).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *momentRepo) DeleteMoment(id uint) error {
	return r.db.Delete(&domain.Moment{}, id).Error
}

func (r *momentRepo) FetchBookings() ([]domain.MomentBooking, error) {
	var bookings []domain.MomentBooking
	err := r.db.Preload("Moment").Order("booking_date desc").Find(&bookings).Error
	return bookings, err
}

func (r *momentRepo) CreateBooking(b *domain.MomentBooking) error {
	return r.db.Create(b).Error
}

func (r *momentRepo) UpdateBookingStatus(id uint, status string) error {
	return r.db.Model(&domain.MomentBooking{}).Where("id = ?", id).Update("status", status).Error
}

// func (r *momentRepo) CheckApprovedBookingExists(date time.Time) (bool, error) {
// 	var count int64
// 	err := r.db.Model(&domain.MomentBooking{}).
// 		Where("DATE(booking_date) = DATE(?) AND status = ?", date, "APPROVED").
// 		Count(&count).Error
// 	return count > 0, err
// }

// Ganti fungsi CheckApprovedBookingExists dengan:
func (r *momentRepo) CheckTimeConflict(start, end time.Time) (bool, error) {
	var count int64
	// Logika Overlap:
	// Booking Baru tumpang tindih dengan Booking Lama (APPROVED) JIKA:
	// (StartBaru < EndLama) AND (EndBaru > StartLama)
	err := r.db.Model(&domain.MomentBooking{}).
		Where("status = ?", "APPROVED").
		Where("booking_date < ? AND booking_end_date > ?", end, start).
		Count(&count).Error
	return count > 0, err
}

// Di fungsi FetchBookings, pastikan mengambil booking berdasarkan ID spesifik 
// agar kita bisa menampilkannya di halaman detail paket tertentu.
// Tambahkan fungsi baru ini:
func (r *momentRepo) FetchApprovedBookingsByMomentID(id uint) ([]domain.MomentBooking, error) {
	var bookings []domain.MomentBooking
	err := r.db.Where("moment_id = ? AND status = ?", id, "APPROVED").Order("booking_date asc").Find(&bookings).Error
	return bookings, err
}

