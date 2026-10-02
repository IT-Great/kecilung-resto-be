// package api

// import (
// 	"log"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/seanalden-great/kecilung-resto-be/config"
// 	httpDelivery "github.com/seanalden-great/kecilung-resto-be/delivery/http"
// 	"github.com/seanalden-great/kecilung-resto-be/domain"
// 	"github.com/seanalden-great/kecilung-resto-be/repository"
// 	"github.com/seanalden-great/kecilung-resto-be/usecase"
// )

// var app *gin.Engine

// // Fungsi init() otomatis dieksekusi pertama kali oleh Vercel sebelum menerima request
// func init() {
// 	// Koneksi Database (menggunakan variabel environment Vercel)
// 	db := config.ConnectDatabase()

// 	// Migrasi Tabel
// 	err := db.AutoMigrate(&domain.Category{}, &domain.Menu{})
// 	if err != nil {
// 		log.Println("Gagal migrasi:", err)
// 	}

// 	// Inisialisasi Arsitektur
// 	catRepo := repository.NewCategoryRepository(db)
// 	catUseCase := usecase.NewCategoryUsecase(catRepo)

// 	menuRepo := repository.NewMenuRepository(db)
// 	menuUseCase := usecase.NewMenuUsecase(menuRepo)

// 	// Set Gin ke mode Release agar log lebih bersih di Vercel
// 	gin.SetMode(gin.ReleaseMode)
// 	app = gin.Default()

// 	// Middleware CORS
// 	app.Use(func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}
// 		c.Next()
// 	})

// 	// Daftarkan Routes
// 	httpDelivery.RegisterHandlers(app, catUseCase, menuUseCase)
// }

// // Handler adalah fungsi standar yang dicari oleh Vercel
// func Handler(w http.ResponseWriter, r *http.Request) {
// 	// Serahkan request HTTP dari Vercel ke dalam router Gin
// 	app.ServeHTTP(w, r)
// }

// package api

// import (
// 	"log"
// 	"net/http"
// 	"os"
// 	"sync" // 1. Tambahkan package sync bawaan Go

// 	"github.com/gin-gonic/gin"
// 	"github.com/seanalden-great/kecilung-resto-be/config"
// 	httpDelivery "github.com/seanalden-great/kecilung-resto-be/delivery/http"
// 	"github.com/seanalden-great/kecilung-resto-be/domain"
// 	"github.com/seanalden-great/kecilung-resto-be/repository"
// 	"github.com/seanalden-great/kecilung-resto-be/usecase"
// 	"golang.org/x/crypto/bcrypt"
// )

// var (
// 	app  *gin.Engine
// 	once sync.Once // 2. Daftarkan variabel penjaga gerbang
// )

// // Kita buat fungsi inisialisasi terpisah (BUKAN init() bawaan Go)
// func initApp() {
// 	// Koneksi Database
// 	db := config.ConnectDatabase()
	
// 	// Migrasi Tabel
// 	err := db.AutoMigrate(
// 		&domain.Category{}, 
// 		&domain.Menu{},
// 		&domain.GreetingMessage{}, 
// 		&domain.Catering{}, 
// 		&domain.CateringImage{}, 
// 		&domain.Booking{},
// 		&domain.Moment{},
// 		&domain.MomentImage{},
// 		&domain.MomentBooking{},
// 		&domain.Article{},         // <--- TAMBAH INI
// 		&domain.ArticleImage{},    // <--- TAMBAH INI
// 		&domain.ContactUs{}, // <--- TAMBAH INI
// 		&domain.Admin{}, // <--- TAMBAH INI
// 	)
// 	if err != nil {
// 		log.Println("Gagal migrasi:", err)
// 	}

// 	// Inisialisasi Arsitektur
// 	catRepo := repository.NewCategoryRepository(db)
// 	catUseCase := usecase.NewCategoryUsecase(catRepo)

// 	menuRepo := repository.NewMenuRepository(db)
// 	menuUseCase := usecase.NewMenuUsecase(menuRepo)

// 	cateringRepo := repository.NewCateringRepository(db)
// 	cateringUseCase := usecase.NewCateringUsecase(cateringRepo)

// 	momentRepo := repository.NewMomentRepository(db)
// 	momentUseCase := usecase.NewMomentUsecase(momentRepo)

// 	// === TAMBAHKAN INISIALISASI ARTICLE DI SINI ===
// 	articleRepo := repository.NewArticleRepository(db)
// 	articleUseCase := usecase.NewArticleUsecase(articleRepo)

// 	// Pada blok Inisialisasi Arsitektur, tambahkan baris berikut:
// 	contactRepo := repository.NewContactRepository(db)
// 	contactUseCase := usecase.NewContactUsecase(contactRepo)

// 	// Pada blok Inisialisasi Arsitektur, tambahkan baris berikut:
// 	authRepo := repository.NewAuthRepository(db)
// 	authUseCase := usecase.NewAuthUsecase(authRepo)

// 	// BUAT AKUN ADMIN DEFAULT (username: admin, password: password123)
// 	// Hanya akan berjalan jika tabel admins masih kosong
// 	hashedPass, _ := bcrypt.GenerateFromPassword([]byte("adminkecilung_!1@2#3"), bcrypt.DefaultCost)
// 	authRepo.CreateDefaultAdmin(&domain.Admin{
// 		Name:     "Admin Kecilung",
// 		Username: "adminkecilung",
// 		Password: string(hashedPass),
// 	})

// 	gin.SetMode(gin.ReleaseMode)
// 	app = gin.Default()

// 	// Middleware CORS
// 	app.Use(func(c *gin.Context) {
// 		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
// 		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}
// 		c.Next()
// 	})

// 	// Endpoint Debug (Untuk memastikan Vercel benar-benar membaca ENV)
// 	app.GET("/api/debug-env", func(c *gin.Context) {
// 		c.JSON(200, gin.H{
// 			"host": os.Getenv("DB_HOST"),
// 			"port": os.Getenv("DB_PORT"),
// 			"user": os.Getenv("DB_USER"),
// 		})
// 	})

// 	// Daftarkan Routes
// 	httpDelivery.RegisterHandlers(app, catUseCase, menuUseCase, cateringUseCase, momentUseCase, articleUseCase, contactUseCase, authUseCase)
// }

// // // Handler ini adalah pintu masuk utama Vercel
// // func Handler(w http.ResponseWriter, r *http.Request) {
// // 	// Jika aplikasi belum diinisialisasi (request pertama), maka inisialisasi sekarang
// // 	if app == nil {
// // 		initApp()
// // 	}
// // 	// Teruskan request ke Gin Router
// // 	app.ServeHTTP(w, r)
// // }

// // 3. Modifikasi fungsi Handler utama Vercel
// func Handler(w http.ResponseWriter, r *http.Request) {
// 	// once.Do memastikan fungsi initApp hanya dieksekusi 1x seumur hidup instance
// 	// Meskipun ratusan request datang menabrak bersamaan, yang lain akan menunggu hingga eksekusi pertama selesai
// 	once.Do(func() {
// 		initApp()
// 	})

// 	// Teruskan request ke Gin Router
// 	app.ServeHTTP(w, r)
// }

package api

import (
	"html/template" // <-- Tambahkan package ini untuk memproses file .tmpl
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/seanalden-great/kecilung-resto-be/config"
	httpDelivery "github.com/seanalden-great/kecilung-resto-be/delivery/http"
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"github.com/seanalden-great/kecilung-resto-be/repository"
	"github.com/seanalden-great/kecilung-resto-be/usecase"
	"golang.org/x/crypto/bcrypt"
)

var (
	app  *gin.Engine
	once sync.Once
)

func initApp() {
	// Koneksi Database
	db := config.ConnectDatabase()
	
	// Migrasi Tabel
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
		log.Println("Gagal migrasi:", err)
	}

	// Inisialisasi Arsitektur
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

	// BUAT AKUN ADMIN DEFAULT
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte("adminkecilung_!1@2#3"), bcrypt.DefaultCost)
	authRepo.CreateDefaultAdmin(&domain.Admin{
		Name:     "Admin Kecilung",
		Username: "adminkecilung",
		Password: string(hashedPass),
	})

	gin.SetMode(gin.ReleaseMode)
	app = gin.Default()

	// === PEMUATAN TEMPLATE HTML (AMAN UNTUK VERCEL) ===
	// Vercel kadang membaca dari root, kadang dari folder /api.
	// Kode ini akan mencari file secara otomatis di kedua tempat tersebut.
	tmpl, err := template.ParseFiles("welcome.tmpl")
	if err != nil {
		tmpl, _ = template.ParseFiles("../welcome.tmpl")
	}
	
	// Jika file ditemukan, daftarkan ke Gin
	if tmpl != nil {
		app.SetHTMLTemplate(tmpl)
	} else {
		log.Println("Peringatan: File welcome.tmpl tidak ditemukan oleh mesin Vercel")
	}

	// Middleware CORS
	app.Use(func(c *gin.Context) {
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
	app.GET("/", func(c *gin.Context) {
		if tmpl != nil {
			// Jika file welcome.tmpl berhasil terbaca, render halaman!
			// Kita menyuntikkan teks "1.0 (Vercel)" ke dalam variabel {{ .version }}
			c.HTML(http.StatusOK, "welcome.tmpl", gin.H{
				"version": "1.0 (Vercel Serverless)",
			})
		} else {
			// Jika file hilang (gagal diupload Vercel), tampilkan pesan biasa (Fallback)
			c.String(http.StatusOK, "Kecilung Resto API - System Online (Vercel)")
		}
	})

	// Endpoint Debug 
	app.GET("/api/debug-env", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"host": os.Getenv("DB_HOST"),
			"port": os.Getenv("DB_PORT"),
			"user": os.Getenv("DB_USER"),
		})
	})

	// Daftarkan Routes
	httpDelivery.RegisterHandlers(app, catUseCase, menuUseCase, cateringUseCase, momentUseCase, articleUseCase, contactUseCase, authUseCase)
}

// Modifikasi fungsi Handler utama Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	once.Do(func() {
		initApp()
	})

	app.ServeHTTP(w, r)
}