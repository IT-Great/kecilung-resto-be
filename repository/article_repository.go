package repository

import (
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"gorm.io/gorm"
)

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) domain.ArticleRepository {
	return &articleRepository{db}
}

func (r *articleRepository) GetAll() ([]domain.Article, error) {
	var articles []domain.Article
	err := r.db.Preload("Images").Find(&articles).Error
	return articles, err
}

func (r *articleRepository) GetByID(id uint) (domain.Article, error) {
	var article domain.Article
	err := r.db.Preload("Images").First(&article, id).Error
	return article, err
}

func (r *articleRepository) Create(article *domain.Article) error {
	return r.db.Create(article).Error
}

// func (r *articleRepository) Update(id uint, article *domain.Article) error {
// 	var existing domain.Article
// 	if err := r.db.First(&existing, id).Error; err != nil {
// 		return err
// 	}

// 	// Jika ada gambar baru yang dikirim, hapus relasi gambar lama agar diganti dengan yang baru
// 	if len(article.Images) > 0 {
// 		r.db.Where("article_id = ?", id).Delete(&domain.ArticleImage{})
// 	}

// 	article.ID = existing.ID
// 	// FullSaveAssociations akan otomatis mengupdate dan menautkan data images baru
// 	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(article).Error
// }

func (r *articleRepository) Update(id uint, article *domain.Article) error {
	var existing domain.Article
	if err := r.db.First(&existing, id).Error; err != nil {
		return err
	}

	// 1. TIMPA DATA LAMA DENGAN DATA BARU SECARA EKSPLISIT
	existing.Code = article.Code
	existing.Name = article.Name
	existing.Description = article.Description
	existing.Slug = article.Slug // <--- INI KUNCI UTAMANYA AGAR SLUG TERSIMPAN

	// 2. JIKA ADA GAMBAR BARU, HAPUS GAMBAR LAMA
	if len(article.Images) > 0 {
		r.db.Where("article_id = ?", id).Delete(&domain.ArticleImage{})
		existing.Images = article.Images
	}

	// 3. GUNAKAN .Save() (BUKAN .Updates()) UNTUK MEMAKSA GORM MENULIS ULANG SEMUA KOLOM
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&existing).Error
}

func (r *articleRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Article{}, id).Error
}

// Tambahkan fungsi ini di bawah
func (r *articleRepository) FetchWithCursor(cursor uint, limit int) ([]domain.Article, error) {
	var articles []domain.Article
	
	query := r.db.Preload("Images").Order("id asc").Limit(limit)
	
	// Jika cursor > 0, ambil data yang ID-nya LEBIH BESAR dari cursor sebelumnya
	if cursor > 0 {
		query = query.Where("id > ?", cursor)
	}

	err := query.Find(&articles).Error
	return articles, err
}

// Tambahkan fungsi baru ini di bawah GetByID:
func (r *articleRepository) GetBySlug(slug string) (domain.Article, error) {
	var article domain.Article
	err := r.db.Preload("Images").Where("slug = ?", slug).First(&article).Error
	return article, err
}