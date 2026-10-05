package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"github.com/seanalden-great/kecilung-resto-be/utils"
)

func RegisterCategoryHandlers(rg *gin.RouterGroup, u domain.CategoryUsecase) {
	cat := rg.Group("/categories")
	{
		cat.GET("", fetchCategories(u))
		// === TAMBAHAN BARU: Smart Handler Detail (Bisa ID, Bisa Slug) ===
		cat.GET("/:param", getCategoryDetail(u))
		cat.POST("", createCategory(u))
		cat.PUT("/:id", updateCategory(u))
		cat.DELETE("/:id", deleteCategory(u))
	}
}

// === HANDLERS KATEGORI ===
func fetchCategories(u domain.CategoryUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := u.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": res})
	}
}

func createCategory(u domain.CategoryUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil input teks dari form-data
		code := c.PostForm("code")
		name := c.PostForm("name")
		description := c.PostForm("description")

		var imageURL string

		// 2. Cek apakah ada file gambar yang dilampirkan
		file, err := c.FormFile("image")
		if err == nil {
			uploadedURL, errUpload := utils.UploadToCleverCloud(file)
			if errUpload != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "S3 Error: " + errUpload.Error()})
				return
			}
			imageURL = uploadedURL
		}

		// 3. Susun struct Kategori
		cat := domain.Category{
			Code:        code,
			Name:        name,
			Description: description,
			ImageURL:    imageURL,
		}

		if err := u.Create(&cat); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Kategori berhasil dibuat", "data": cat})
	}
}

func updateCategory(u domain.CategoryUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))

		// 1. Ambil data kategori lama
		existingCat, err := u.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan"})
			return
		}

		// 2. Ambil data baru
		code := c.PostForm("code")
		name := c.PostForm("name")
		description := c.PostForm("description")

		// Set default gambar ke gambar lama
		imageURL := existingCat.ImageURL

		// 3. Cek apakah ada gambar baru
		file, err := c.FormFile("image")
		if err == nil {
			uploadedURL, errUpload := utils.UploadToCleverCloud(file)
			if errUpload != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "S3 Error: " + errUpload.Error()})
				return
			}
			imageURL = uploadedURL
		}

		cat := domain.Category{
			Code:        code,
			Name:        name,
			Description: description,
			ImageURL:    imageURL,
		}

		if err := u.Update(uint(id), &cat); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Kategori berhasil diupdate"})
	}
}

func deleteCategory(u domain.CategoryUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := u.Delete(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Kategori berhasil dihapus"})
	}
}

// === TAMBAHKAN FUNGSI INI DI PALING BAWAH ===
func getCategoryDetail(u domain.CategoryUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("param")

		// 1. Cek jika Angka (ID)
		if id, err := strconv.Atoi(param); err == nil {
			res, err := u.GetByID(uint(id))
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan (ID)"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"data": res})
			return
		}

		// 2. Cek jika Huruf (Slug)
		res, err := u.GetBySlug(param)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Kategori tidak ditemukan (Slug)"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": res})
	}
}