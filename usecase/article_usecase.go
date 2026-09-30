package usecase

import "github.com/seanalden-great/kecilung-resto-be/domain"

type articleUsecase struct {
	repo domain.ArticleRepository
}

func NewArticleUsecase(repo domain.ArticleRepository) domain.ArticleUsecase {
	return &articleUsecase{repo}
}

func (u *articleUsecase) GetAll() ([]domain.Article, error) {
	return u.repo.GetAll()
}

func (u *articleUsecase) GetByID(id uint) (domain.Article, error) {
	return u.repo.GetByID(id)
}

func (u *articleUsecase) Create(article *domain.Article) error {
	return u.repo.Create(article)
}

func (u *articleUsecase) Update(id uint, article *domain.Article) error {
	return u.repo.Update(id, article)
}

func (u *articleUsecase) Delete(id uint) error {
	return u.repo.Delete(id)
}

// Tambahkan fungsi ini di bawah
func (u *articleUsecase) GetWithCursor(cursor uint, limit int) (domain.ArticleCursorResponse, error) {
	// Ambil data LIMIT + 1 untuk mengecek apakah masih ada halaman selanjutnya
	// Jika limit=5, kita ambil 6. Jika yang kembali 6, artinya "HasMore = true".
	articles, err := u.repo.FetchWithCursor(cursor, limit+1)
	if err != nil {
		return domain.ArticleCursorResponse{}, err
	}

	hasMore := false
	var resultData []domain.Article

	if len(articles) > limit {
		hasMore = true
		resultData = articles[:limit] // Potong elemen ekstra (ke-6)
	} else {
		resultData = articles
	}

	var nextCursor uint = 0
	if len(resultData) > 0 {
		// Cursor selanjutnya adalah ID dari elemen terakhir di array saat ini
		nextCursor = resultData[len(resultData)-1].ID
	}

	return domain.ArticleCursorResponse{
		Data:       resultData,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}