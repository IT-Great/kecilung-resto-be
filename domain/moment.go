package domain

import "time"

type Moment struct {
	ID          uint          `json:"id" gorm:"primaryKey"`
	Name        string        `json:"name"`
	Slug        string        `json:"slug" gorm:"uniqueIndex"` // <--- TAMBAHAN: Kolom Slug Unik
	Description string        `json:"description"`
	Images      []MomentImage `json:"images" gorm:"foreignKey:MomentID;constraint:OnDelete:CASCADE;"`
}

type MomentImage struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	MomentID uint   `json:"moment_id"`
	ImageURL string `json:"image_url"`
}

type MomentBooking struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	MomentID       uint      `json:"moment_id"`
	Moment         Moment    `json:"moment" gorm:"foreignKey:MomentID"`
	CustomerName   string    `json:"customer_name"`
	Phone          string    `json:"phone"`
	Description    string    `json:"description"`  // BARU: Keterangan kegiatan
	MemberCount    int       `json:"member_count"` // BARU: Jumlah orang
	BookingDate    time.Time `json:"booking_date"`
	BookingEndDate time.Time `json:"booking_end_date"` // BARU: Waktu Selesai
	Status         string    `json:"status" gorm:"default:'PENDING'"`
}

type MomentRepository interface {
	FetchAllMoments() ([]Moment, error)
	FindMomentByID(id uint) (Moment, error)
	FindMomentBySlug(slug string) (Moment, error) // <--- TAMBAHAN: Fungsi Cari via Slug
	CreateMoment(m *Moment) error
	UpdateMoment(m *Moment) error
	DeleteMoment(id uint) error

	FetchBookings() ([]MomentBooking, error)
	CreateBooking(b *MomentBooking) error
	UpdateBookingStatus(id uint, status string) error
	// CheckApprovedBookingExists(date time.Time) (bool, error)
	CheckTimeConflict(start, end time.Time) (bool, error)
	FetchApprovedBookingsByMomentID(id uint) ([]MomentBooking, error)
}

type MomentUsecase interface {
	GetAllMoments() ([]Moment, error)
	GetMomentByID(id uint) (Moment, error)
	GetMomentBySlug(slug string) (Moment, error) // <--- TAMBAHAN: Fungsi Cari via Slug
	CreateMoment(m *Moment) error
	UpdateMoment(id uint, m *Moment) error
	DeleteMoment(id uint) error

	GetAllBookings() ([]MomentBooking, error)
	CreateBooking(b *MomentBooking) error
	ApproveBooking(id uint) error
	RejectBooking(id uint) error

	GetApprovedBookingsByMomentID(id uint) ([]MomentBooking, error)
}
