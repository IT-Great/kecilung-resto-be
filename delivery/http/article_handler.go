package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"github.com/seanalden-great/kecilung-resto-be/utils"
)

// func RegisterArticleHandlers(rg *gin.RouterGroup, u domain.ArticleUsecase) {
// // === TAMBAHKAN ROUTES ARTICLE DI SINI ===
// 	article := rg.Group("/articles")
// 	{
// 		// === RUTE BARU UNTUK INFINITE SCROLL (PUBLIK) ===
// 		article.GET("/feed", getArticlesFeed(u))

// 		article.GET("", getArticles(u))
// 		article.GET("/:slug", getArticleBySlug(u))
// 		article.POST("", createArticle(u))
// 		article.PUT("/:id", updateArticle(u))
// 		article.DELETE("/:id", deleteArticle(u))
// 	}
// }

func RegisterArticleHandlers(rg *gin.RouterGroup, u domain.ArticleUsecase) {
	article := rg.Group("/articles")
	{
		article.GET("/feed", getArticlesFeed(u))
		article.GET("", getArticles(u))
		
		// UBAH BARIS INI: Gunakan :param dan panggil fungsi pendeteksi pintar
		article.GET("/:param", getArticleDetail(u)) 
		
		article.POST("", createArticle(u))
		article.PUT("/:id", updateArticle(u))
		article.DELETE("/:id", deleteArticle(u))
	}
}

// ====================================================================
// === HANDLERS ARTICLE (TAMBAHKAN KODE INI DI BAGIAN PALING BAWAH) ===
// ====================================================================

func getArticles(u domain.ArticleUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := u.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": res})
	}
}

// func getArticleByID(u domain.ArticleUsecase) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		id, _ := strconv.Atoi(c.Param("id"))
// 		res, err := u.GetByID(uint(id))
// 		if err != nil {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Data artikel tidak ditemukan"})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"data": res})
// 	}
// }

// // UBAH fungsi getArticleByID menjadi getArticleBySlug
// func getArticleBySlug(u domain.ArticleUsecase) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		slug := c.Param("slug") // Tangkap slug berupa huruf
// 		res, err := u.GetBySlug(slug)
// 		if err != nil {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Data artikel tidak ditemukan"})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"data": res})
// 	}
// }

// === FUNGSI PINTAR: BISA MENANGKAP ID (ANGKA) ATAU SLUG (HURUF) ===
func getArticleDetail(u domain.ArticleUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("param")

		// 1. Coba ubah param menjadi Angka (Integer)
		if id, err := strconv.Atoi(param); err == nil {
			// JIKA BERHASIL (Berarti ini ID dari halaman Admin)
			res, err := u.GetByID(uint(id))
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Data artikel tidak ditemukan (ID)"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"data": res})
			return
		}

		// 2. JIKA GAGAL DIUBAH KE ANGKA (Berarti ini berupa Huruf/Slug dari halaman Publik)
		res, err := u.GetBySlug(param)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data artikel tidak ditemukan (Slug)"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": res})
	}
}

func createArticle(u domain.ArticleUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.PostForm("code")
		name := c.PostForm("name")
		description := c.PostForm("description")

		var images []domain.ArticleImage

		// Deteksi multi-upload
		form, err := c.MultipartForm()
		if err == nil {
			files := form.File["images"]
			for _, file := range files {
				url, errUp := utils.UploadToCleverCloud(file)
				if errUp == nil {
					images = append(images, domain.ArticleImage{ImageURL: url})
				}
			}
		}

		article := domain.Article{Code: code, Name: name, Description: description, Images: images}
		if err := u.Create(&article); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Artikel berhasil ditambahkan", "data": article})
	}
}

func updateArticle(u domain.ArticleUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		code := c.PostForm("code")
		name := c.PostForm("name")
		description := c.PostForm("description")

		var images []domain.ArticleImage
		form, err := c.MultipartForm()
		if err == nil {
			files := form.File["images"]
			for _, file := range files {
				url, errUp := utils.UploadToCleverCloud(file)
				if errUp == nil {
					images = append(images, domain.ArticleImage{ImageURL: url})
				}
			}
		}

		article := domain.Article{Code: code, Name: name, Description: description, Images: images}
		if err := u.Update(uint(id), &article); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Artikel berhasil diupdate"})
	}
}

func deleteArticle(u domain.ArticleUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := u.Delete(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Artikel berhasil dihapus"})
	}
}

// === HANDLER BARU: INFINITE SCROLL ===
func getArticlesFeed(u domain.ArticleUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		cursorStr := c.DefaultQuery("cursor", "0")
		limitStr := c.DefaultQuery("limit", "5") // Default load 5 artikel per batch

		cursor, _ := strconv.Atoi(cursorStr)
		limit, _ := strconv.Atoi(limitStr)

		res, err := u.GetWithCursor(uint(cursor), limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat feed artikel"})
			return
		}
		
		c.JSON(http.StatusOK, res)
	}
}