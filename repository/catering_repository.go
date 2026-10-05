package repository

import (
	"time"
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"gorm.io/gorm"
)

type cateringRepo struct {
	db *gorm.DB
}

func NewCateringRepository(db *gorm.DB) domain.CateringRepository {
	return &cateringRepo{db}
}

// -- Greeting --
func (r *cateringRepo) GetGreeting() (domain.GreetingMessage, error) {
	var msg domain.GreetingMessage
	err := r.db.FirstOrCreate(&msg, domain.GreetingMessage{ID: 1, Message: "Selamat Datang di Layanan Katering Kami"}).Error
	return msg, err
}

func (r *cateringRepo) UpdateGreeting(msg *domain.GreetingMessage) error {
	msg.ID = 1
	return r.db.Save(msg).Error
}

// -- Catering --
func (r *cateringRepo) FetchAllCaterings() ([]domain.Catering, error) {
	var caterings []domain.Catering
	err := r.db.Preload("Images").Find(&caterings).Error
	return caterings, err
}

func (r *cateringRepo) FindCateringByID(id uint) (domain.Catering, error) {
	var catering domain.Catering
	err := r.db.Preload("Images").First(&catering, id).Error
	return catering, err
}

// === TAMBAHAN BARU: Cari via Slug ===
func (r *cateringRepo) FindCateringBySlug(slug string) (domain.Catering, error) {
	var catering domain.Catering
	err := r.db.Preload("Images").Where("slug = ?", slug).First(&catering).Error
	return catering, err
}

func (r *cateringRepo) CreateCatering(c *domain.Catering) error {
	return r.db.Create(c).Error
}

// === TAMBAHAN BARU: UPDATE CATERING ===
// func (r *cateringRepo) UpdateCatering(c *domain.Catering) error {
// 	// Gunakan transaksi agar jika gagal di tengah jalan, database di-rollback
// 	return r.db.Transaction(func(tx *gorm.DB) error {
// 		// 1. Update data teks utama
// 		if err := tx.Model(&domain.Catering{}).Where("id = ?", c.ID).Updates(map[string]interface{}{
// 			"name":        c.Name,
// 			"description": c.Description,
// 		}).Error; err != nil {
// 			return err
// 		}

// 		// 2. Jika ada gambar baru yang diunggah, hapus gambar lama dan masukkan yang baru
// 		if len(c.Images) > 0 {
// 			if err := tx.Where("catering_id = ?", c.ID).Delete(&domain.CateringImage{}).Error; err != nil {
// 				return err
// 			}
// 			for i := range c.Images {
// 				c.Images[i].CateringID = c.ID
// 			}
// 			if err := tx.Create(&c.Images).Error; err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	})
// }

func (r *cateringRepo) UpdateCatering(c *domain.Catering) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existing domain.Catering
		if err := tx.First(&existing, c.ID).Error; err != nil {
			return err
		}

		// Update kolom dasar (secara eksplisit, termasuk slug)
		existing.Name = c.Name
		existing.Description = c.Description
		existing.Slug = c.Slug // <-- PASTI TERSIMPAN

		if err := tx.Save(&existing).Error; err != nil {
			return err
		}

		if len(c.Images) > 0 {
			if err := tx.Where("catering_id = ?", c.ID).Delete(&domain.CateringImage{}).Error; err != nil {
				return err
			}
			for i := range c.Images {
				c.Images[i].CateringID = c.ID
			}
			if err := tx.Create(&c.Images).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// === TAMBAHAN BARU: DELETE CATERING ===
func (r *cateringRepo) DeleteCatering(id uint) error {
	// Berkat constraint:OnDelete:CASCADE di struct, menghapus Catering otomatis menghapus semua CateringImage miliknya
	return r.db.Delete(&domain.Catering{}, id).Error
}

// -- Booking --
func (r *cateringRepo) FetchBookings() ([]domain.Booking, error) {
	var bookings []domain.Booking
	err := r.db.Preload("Catering").Order("booking_date desc").Find(&bookings).Error
	return bookings, err
}

func (r *cateringRepo) CreateBooking(b *domain.Booking) error {
	return r.db.Create(b).Error
}

func (r *cateringRepo) UpdateBookingStatus(id uint, status string) error {
	return r.db.Model(&domain.Booking{}).Where("id = ?", id).Update("status", status).Error
}

// func (r *cateringRepo) CheckApprovedBookingExists(date time.Time) (bool, error) {
// 	var count int64
// 	// Cek apakah ada booking berstatus APPROVED pada hari (tanggal) yang sama
// 	err := r.db.Model(&domain.Booking{}).
// 		Where("DATE(booking_date) = DATE(?) AND status = ?", date, "APPROVED").
// 		Count(&count).Error
// 	return count > 0, err
// }

// Ganti fungsi CheckApprovedBookingExists dengan:
func (r *cateringRepo) CheckTimeConflict(start, end time.Time) (bool, error) {
	var count int64
	// Logika Overlap:
	// Booking Baru tumpang tindih dengan Booking Lama (APPROVED) JIKA:
	// (StartBaru < EndLama) AND (EndBaru > StartLama)
	err := r.db.Model(&domain.Booking{}).
		Where("status = ?", "APPROVED").
		Where("booking_date < ? AND booking_end_date > ?", end, start).
		Count(&count).Error
	return count > 0, err
}

// Di fungsi FetchBookings, pastikan mengambil booking berdasarkan ID spesifik 
// agar kita bisa menampilkannya di halaman detail paket tertentu.
// Tambahkan fungsi baru ini:
func (r *cateringRepo) FetchApprovedBookingsByCateringID(id uint) ([]domain.Booking, error) {
	var bookings []domain.Booking
	err := r.db.Where("catering_id = ? AND status = ?", id, "APPROVED").Order("booking_date asc").Find(&bookings).Error
	return bookings, err
}


