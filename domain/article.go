package domain

import "time"

// Article struct untuk tabel articles
type Article struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Code        string         `json:"code"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug" gorm:"uniqueIndex"` // <-- TAMBAHAN BARU: Kolom Slug unik
	Description string         `json:"description"`
	Images      []ArticleImage `json:"images" gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE;"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// ArticleImage struct untuk multi-gambar
type ArticleImage struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ArticleID uint      `json:"article_id"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CursorResponse adalah format data balikan untuk Infinite Scroll
type ArticleCursorResponse struct {
	Data       []Article `json:"data"`
	NextCursor uint      `json:"next_cursor"` // ID terakhir dari batch saat ini, 0 jika habis
	HasMore    bool      `json:"has_more"`
}

// ArticleRepository interface
type ArticleRepository interface {
	GetAll() ([]Article, error)
	GetByID(id uint) (Article, error)
	GetBySlug(slug string) (Article, error) // <-- TAMBAHAN BARU: Cari berdasarkan Slug
	Create(article *Article) error
	Update(id uint, article *Article) error
	Delete(id uint) error

	// TAMBAHAN BARU:
	FetchWithCursor(cursor uint, limit int) ([]Article, error)
}

// ArticleUsecase interface
type ArticleUsecase interface {
	GetAll() ([]Article, error)
	GetByID(id uint) (Article, error)
	GetBySlug(slug string) (Article, error) // <-- TAMBAHAN BARU: Cari berdasarkan Slug
	Create(article *Article) error
	Update(id uint, article *Article) error
	Delete(id uint) error

	// TAMBAHAN BARU:
	GetWithCursor(cursor uint, limit int) (ArticleCursorResponse, error)
}