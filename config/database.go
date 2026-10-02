// package config

// import (
// 	"fmt"
// 	"log"
// 	"os"

// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// func ConnectDatabase() *gorm.DB {
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
// 		os.Getenv("DB_USER"),
// 		os.Getenv("DB_PASS"),
// 		os.Getenv("DB_HOST"),
// 		os.Getenv("DB_PORT"),
// 		os.Getenv("DB_NAME"),
// 	)

// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal("Gagal terkoneksi ke database:", err)
// 	}
// 	return db
// }

package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")

	// 1. Validasi Keberadaan Environment Variables
	if host == "" || port == "" {
		log.Fatal("CRITICAL ERROR: DB_HOST atau DB_PORT kosong! Vercel belum menyuntikkan variabel environment.")
	}

	// 2. Format DSN (Tambahkan &tls=true di bagian belakang khusus untuk Cloud Database)
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=true",
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=false",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		host,
		port,
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terkoneksi ke database:", err)
	}
	
	log.Println("Sukses terhubung ke database:", host)
	return db
}