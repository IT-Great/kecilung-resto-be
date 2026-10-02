// package main

// import (
// 	"log"
// 	"os"

// 	"github.com/gin-gonic/gin"
// 	"github.com/joho/godotenv"
// 	"github.com/seanalden-great/kecilung-resto-be/config"
// 	httpDelivery "github.com/seanalden-great/kecilung-resto-be/delivery/http"
// 	"github.com/seanalden-great/kecilung-resto-be/domain"
// 	"github.com/seanalden-great/kecilung-resto-be/repository"
// 	"github.com/seanalden-great/kecilung-resto-be/usecase"
// )

// func main() {
// 	_ = godotenv.Load()

// 	db := config.ConnectDatabase()

// 	// MIGRATION: Category harus di atas Menu karena Menu bergantung pada ID Category
// 	err := db.AutoMigrate(&domain.Category{}, &domain.Menu{})
// 	if err != nil {
// 		log.Fatal("Gagal migrasi tabel:", err)
// 	}

// 	// WIRING ARCHITECTURE
// 	catRepo := repository.NewCategoryRepository(db)
// 	catUseCase := usecase.NewCategoryUsecase(catRepo)

// 	menuRepo := repository.NewMenuRepository(db)
// 	menuUseCase := usecase.NewMenuUsecase(menuRepo)

// 	// INIT ROUTER
// 	r := gin.Default()

// 	// MIDDLEWARE CORS
// 	r.Use(func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}
// 		c.Next()
// 	})

// 	// REGISTER ENDPOINTS
// 	httpDelivery.RegisterHandlers(r, catUseCase, menuUseCase)

// 	// r.Run(":8080")
// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "8080" // Fallback untuk lokal Laragon
// 	}
// 	r.Run(":" + port)
// }

package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/seanalden-great/kecilung-resto-be/config"
	httpDelivery "github.com/seanalden-great/kecilung-resto-be/delivery/http"
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"github.com/seanalden-great/kecilung-resto-be/repository"
	"github.com/seanalden-great/kecilung-resto-be/usecase"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load .env HANYA untuk lokal
	_ = godotenv.Load()

	// Koneksi Database
	db := config.ConnectDatabase()

	// === 1. MIGRATION ===
	// Pastikan semua tabel terbaru terdaftar di sini, persis seperti api/index.go
	err := db.AutoMigrate(
		&domain.Category{},
		&domain.Menu{},
		&domain.GreetingMessage{},
		&domain.Catering{},
		&domain.CateringImage{},
		&domain.Booking{},
		&domain.Moment{},
		&domain.MomentImage{},
		&domain.MomentBooking{},
		&domain.Article{},
		&domain.ArticleImage{},
		&domain.ContactUs{},
		&domain.Admin{},
	)
	if err != nil {
		log.Fatal("Gagal migrasi tabel:", err)
	}

	// === 2. WIRING ARCHITECTURE ===
	// Inisialisasi semua Usecase dan Repository, persis seperti api/index.go

	catRepo := repository.NewCategoryRepository(db)
	catUseCase := usecase.NewCategoryUsecase(catRepo)

	menuRepo := repository.NewMenuRepository(db)
	menuUseCase := usecase.NewMenuUsecase(menuRepo)

	cateringRepo := repository.NewCateringRepository(db)
	cateringUseCase := usecase.NewCateringUsecase(cateringRepo)

	momentRepo := repository.NewMomentRepository(db)
	momentUseCase := usecase.NewMomentUsecase(momentRepo)

	articleRepo := repository.NewArticleRepository(db)
	articleUseCase := usecase.NewArticleUsecase(articleRepo)

	contactRepo := repository.NewContactRepository(db)
	contactUseCase := usecase.NewContactUsecase(contactRepo)

	authRepo := repository.NewAuthRepository(db)
	authUseCase := usecase.NewAuthUsecase(authRepo)

	// === 3. INJEKSI AKUN ADMIN DEFAULT ===
	// Hanya akan berjalan jika tabel admins masih kosong (Sama seperti di Vercel)
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte("adminkecilung_!1@2#3"), bcrypt.DefaultCost)
	authRepo.CreateDefaultAdmin(&domain.Admin{
		Name:     "Admin Kecilung",
		Username: "adminkecilung",
		Password: string(hashedPass),
	})

	// === 4. INIT ROUTER ===
	r := gin.Default()

	// MIDDLEWARE CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// === HALAMAN WELCOME API ===
	r.GET("/", func(c *gin.Context) {
		htmlTemplate := `
		<!DOCTYPE html>
		<html lang="id">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Kecilung Resto API</title>
			<style>
				body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif; background-color: #1a202c; color: #a0aec0; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; }
				.container { text-align: center; }
				h1 { color: #edf2f7; font-size: 3rem; margin-bottom: 0.5rem; font-weight: 300; letter-spacing: -1px; }
				p { font-size: 1.2rem; }
				.version { margin-top: 2.5rem; font-size: 0.85rem; color: #718096; font-weight: 600; text-transform: uppercase; letter-spacing: 1px; }
				.status { display: inline-block; padding: 4px 12px; background-color: #2f855a; color: #c6f6d5; border-radius: 9999px; font-size: 0.85rem; font-weight: 600; margin-bottom: 1rem; }
			</style>
		</head>
		<body>
			<div class="container">
				<div class="status">System Online</div>
				<h1>Kecilung Resto API</h1>
				<p>Backend service is active and listening for requests.</p>
				<div class="version">Powered by Go & Gin Framework</div>
			</div>
		</body>
		</html>
		`
		c.Data(200, "text/html; charset=utf-8", []byte(htmlTemplate))
	})

	// Endpoint Debug (Opsional, untuk konsistensi)
	r.GET("/api/debug-env", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"host": os.Getenv("DB_HOST"),
			"port": os.Getenv("DB_PORT"),
			"user": os.Getenv("DB_USER"),
		})
	})

	// === 5. REGISTER ENDPOINTS ===
	// PERUBAHAN: Masukkan SEMUA parameter Usecase di sini agar error hilang
	httpDelivery.RegisterHandlers(r, catUseCase, menuUseCase, cateringUseCase, momentUseCase, articleUseCase, contactUseCase, authUseCase)

	// === 6. RUN SERVER ===
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback untuk lokal
	}
	
	log.Printf("Server lokal berjalan di http://localhost:%s\n", port)
	r.Run(":" + port)
}