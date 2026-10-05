package usecase

import (
	"regexp"
	"strings"

	"github.com/seanalden-great/kecilung-resto-be/domain"
)

type categoryUsecase struct {
	repo domain.CategoryRepository
}

func NewCategoryUsecase(r domain.CategoryRepository) domain.CategoryUsecase {
	return &categoryUsecase{repo: r}
}

// === Helper Pembentuk Slug ===
func generateCategorySlug(title string) string {
	str := strings.ToLower(title)
	re := regexp.MustCompile("[^a-z0-9]+")
	str = re.ReplaceAllString(str, "-")
	return strings.Trim(str, "-")
}

func (u *categoryUsecase) GetAll() ([]domain.Category, error) {
	return u.repo.FetchAll()
}

func (u *categoryUsecase) GetByID(id uint) (domain.Category, error) {
	return u.repo.FindByID(id)
}

// === TAMBAHAN BARU ===
func (u *categoryUsecase) GetBySlug(slug string) (domain.Category, error) {
	return u.repo.FindBySlug(slug)
}

func (u *categoryUsecase) Create(c *domain.Category) error {
	return u.repo.Store(c)
}

// func (u *categoryUsecase) Update(id uint, c *domain.Category) error {
// 	existing, err := u.repo.FindByID(id)
// 	if err != nil {
// 		return err
// 	}
// 	c.ID = existing.ID
// 	return u.repo.Update(c)
// }

func (u *categoryUsecase) Update(id uint, c *domain.Category) error {
	existing, err := u.repo.FindByID(id)
	if err != nil {
		return err
	}
	c.ID = existing.ID
	c.Slug = generateCategorySlug(c.Name) // Otomatis update Slug jika nama berubah
	
	// Jika gambar kosong, pertahankan gambar lama (agar tidak hilang)
	if c.ImageURL == "" {
		c.ImageURL = existing.ImageURL
	}
	
	return u.repo.Update(c)
}

func (u *categoryUsecase) Delete(id uint) error {
	return u.repo.Delete(id)
}