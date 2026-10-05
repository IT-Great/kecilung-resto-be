package domain

type Category struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Code        string `json:"code" gorm:"size:50;uniqueIndex"` // Tambahkan size:50 di sini
	Name        string `json:"name" gorm:"size:255"`            // Menjadi varchar(255)
	Slug        string `json:"slug" gorm:"uniqueIndex"` // <--- TAMBAHAN BARU: Kolom Slug Unik
	Description string `json:"description"`                     // Tetap longtext
	ImageURL    string `json:"image_url"` // === TAMBAHAN BARU ===
}

type CategoryRepository interface {
	FetchAll() ([]Category, error)
	FindByID(id uint) (Category, error)
	FindBySlug(slug string) (Category, error) // <--- TAMBAHAN BARU
	Store(c *Category) error
	Update(c *Category) error
	Delete(id uint) error
}

type CategoryUsecase interface {
	GetAll() ([]Category, error)
	GetByID(id uint) (Category, error)
	GetBySlug(slug string) (Category, error) // <--- TAMBAHAN BARU
	Create(c *Category) error
	Update(id uint, c *Category) error
	Delete(id uint) error
}