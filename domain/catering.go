package domain

import "time"

type GreetingMessage struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	Message string `json:"message"`
}

type Catering struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	Name        string          `json:"name"`
	Slug        string          `json:"slug" gorm:"uniqueIndex"` // <--- TAMBAHAN: Kolom Slug Unik
	Description string          `json:"description"`
	Images      []CateringImage `json:"images" gorm:"foreignKey:CateringID;constraint:OnDelete:CASCADE;"`
}

type CateringImage struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	CateringID uint   `json:"catering_id"`
	ImageURL   string `json:"image_url"`
}

type Booking struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	CateringID     uint      `json:"catering_id"`
	Catering       Catering  `json:"catering" gorm:"foreignKey:CateringID"`
	CustomerName   string    `json:"customer_name"`
	Phone          string    `json:"phone"`
	Description    string    `json:"description"`  // BARU: Keterangan kegiatan
	MemberCount    int       `json:"member_count"` // BARU: Jumlah orang
	BookingDate    time.Time `json:"booking_date"`
	BookingEndDate time.Time `json:"booking_end_date"`                // BARU: Waktu Selesai
	Status         string    `json:"status" gorm:"default:'PENDING'"` // PENDING, APPROVED, REJECTED
}

type CateringRepository interface {
	GetGreeting() (GreetingMessage, error)
	UpdateGreeting(msg *GreetingMessage) error

	FetchAllCaterings() ([]Catering, error)
	FindCateringByID(id uint) (Catering, error)
	FindCateringBySlug(slug string) (Catering, error) // <--- TAMBAHAN: Fungsi Cari via Slug
	CreateCatering(c *Catering) error
	UpdateCatering(c *Catering) error // Fungsi Update Baru
	DeleteCatering(id uint) error

	FetchBookings() ([]Booking, error)
	CreateBooking(b *Booking) error
	UpdateBookingStatus(id uint, status string) error
	// CheckApprovedBookingExists(date time.Time) (bool, error)
	CheckTimeConflict(start, end time.Time) (bool, error)
	FetchApprovedBookingsByCateringID(id uint) ([]Booking, error)
}

type CateringUsecase interface {
	GetGreeting() (GreetingMessage, error)
	UpdateGreeting(msg *GreetingMessage) error

	GetAllCaterings() ([]Catering, error)
	GetCateringByID(id uint) (Catering, error)
	GetCateringBySlug(slug string) (Catering, error) // <--- TAMBAHAN: Fungsi Cari via Slug
	CreateCatering(c *Catering) error
	UpdateCatering(id uint, c *Catering) error // Fungsi Update Baru
	DeleteCatering(id uint) error

	GetAllBookings() ([]Booking, error)
	CreateBooking(b *Booking) error
	ApproveBooking(id uint) error
	RejectBooking(id uint) error

	GetApprovedBookingsByCateringID(id uint) ([]Booking, error)
}
