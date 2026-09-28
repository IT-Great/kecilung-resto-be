package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"github.com/seanalden-great/kecilung-resto-be/utils"
)

func RegisterAuthHandlers(rg *gin.RouterGroup, u domain.AuthUsecase) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", login(u))
		auth.GET("/profile/:id", getProfile(u))
		auth.PUT("/profile/:id", updateProfile(u))

		// TAMBAHAN ROUTES LUPA PASSWORD
		auth.POST("/forgot-password", requestOTP(u))
		auth.POST("/verify-otp", verifyOTP(u))
		auth.POST("/reset-password", resetPassword(u))
	}
}

// === HANDLERS AUTH (Letakkan di bagian bawah file) ===
func login(u domain.AuthUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		token, admin, err := u.Login(req.Username, req.Password)
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Login berhasil", "token": token, "data": admin})
	}
}

func getProfile(u domain.AuthUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		admin, err := u.GetProfile(uint(id))
		if err != nil {
			c.JSON(404, gin.H{"error": "Profile tidak ditemukan"})
			return
		}
		c.JSON(200, gin.H{"data": admin})
	}
}

func updateProfile(u domain.AuthUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		
		name := c.PostForm("name")
		username := c.PostForm("username")
		email := c.PostForm("email") // TAMBAHAN: Ambil nilai email dari form data
		password := c.PostForm("password") // Opsional

		var imageURL string
		file, err := c.FormFile("image")
		if err == nil {
			uploadedURL, errUpload := utils.UploadToCleverCloud(file)
			if errUpload == nil {
				imageURL = uploadedURL
			}
		}

		adminData := domain.Admin{
			Name:     name,
			Username: username,
			Email:    email, // TAMBAHAN: Masukkan ke struct
			Password: password,
		}
		
		if imageURL != "" {
			adminData.ImageURL = imageURL
		}

		if err := u.UpdateProfile(uint(id), &adminData); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Profil berhasil diperbarui"})
	}
}

func requestOTP(u domain.AuthUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct { Email string `json:"email"` }
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Format email tidak valid"})
			return
		}
		if err := u.RequestOTP(req.Email); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Kode OTP telah dikirim ke email Anda"})
	}
}

func verifyOTP(u domain.AuthUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct { 
			Email string `json:"email"`
			OTP   string `json:"otp"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Input tidak valid"})
			return
		}
		if err := u.VerifyOTP(req.Email, req.OTP); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "OTP valid, silakan atur password baru"})
	}
}

func resetPassword(u domain.AuthUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct { 
			Email       string `json:"email"`
			NewPassword string `json:"new_password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Input tidak valid"})
			return
		}
		if err := u.ResetPassword(req.Email, req.NewPassword); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "Password berhasil diubah, silakan login"})
	}
}