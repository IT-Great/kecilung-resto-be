package repository

import (
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"gorm.io/gorm"
)

type categoryRepo struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) domain.CategoryRepository {
	return &categoryRepo{db}
}

// func (r *categoryRepo) FetchAll() ([]domain.Category, error) {
// 	var categories []domain.Category
// 	err := r.db.Find(&categories).Error
// 	return categories, err
// }

func (r *categoryRepo) FetchAll() ([]domain.Category, error) {
	var categories []domain.Category
	// Urutkan berdasarkan code dari terkecil ke terbesar
	err := r.db.Order("code asc").Find(&categories).Error
	return categories, err
}

func (r *categoryRepo) FindByID(id uint) (domain.Category, error) {
	var category domain.Category
	err := r.db.First(&category, id).Error
	return category, err
}

// === TAMBAHAN BARU: Cari via Slug ===
func (r *categoryRepo) FindBySlug(slug string) (domain.Category, error) {
	var category domain.Category
	err := r.db.Where("slug = ?", slug).First(&category).Error
	return category, err
}

func (r *categoryRepo) Store(c *domain.Category) error {
	return r.db.Create(c).Error
}

func (r *categoryRepo) Update(c *domain.Category) error {
	return r.db.Save(c).Error
}

func (r *categoryRepo) Delete(id uint) error {
	return r.db.Delete(&domain.Category{}, id).Error
}