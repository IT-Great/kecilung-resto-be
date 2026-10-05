package http

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"github.com/seanalden-great/kecilung-resto-be/utils"
	"github.com/xuri/excelize/v2"
)

func RegisterCateringHandlers(rg *gin.RouterGroup, u domain.CateringUsecase) {
	// Routes Catering
	catering := rg.Group("/catering")
	{
		catering.GET("/greeting", getGreeting(u))
		catering.PUT("/greeting", updateGreeting(u))
		catering.GET("/packages", getCaterings(u))
		// catering.GET("/packages/:id", getCateringPackageByID(u))
		// UBAH BARIS INI: Gunakan :param dan panggil fungsi pintar
		catering.GET("/packages/:param", getCateringDetail(u))
		catering.POST("/packages", createCateringPackage(u))
		// Tambahkan 2 baris ini di dalam RegisterHandlers -> catering := api.Group("/catering")
		catering.PUT("/packages/:id", updateCateringPackage(u))
		catering.DELETE("/packages/:id", deleteCateringPackage(u))
		catering.GET("/bookings", getBookings(u))
		catering.POST("/bookings", createBooking(u))
		catering.PUT("/bookings/:id/approve", approveBooking(u))
		catering.PUT("/bookings/:id/reject", rejectBooking(u))
		// catering.GET("/packages/:id/bookings", getApprovedCateringBookings(u))
		catering.GET("/bookings/approved/:id", getApprovedCateringBookings(u))
	}
}

// ====================================================================
// === HANDLERS CATERING (TAMBAHKAN KODE INI DI BAGIAN PALING BAWAH) ===
// ====================================================================

// Lalu buat handlernya di bawah:
func getApprovedCateringBookings(u domain.CateringUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		res, err := u.GetApprovedBookingsByCateringID(uint(id))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": res})
	}
}

func getGreeting(u domain.CateringUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		msg, err := u.GetGreeting()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": msg})
	}
}

func updateGreeting(u domain.CateringUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var msg domain.GreetingMessage
		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := u.UpdateGreeting(&msg); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Pesan sambutan berhasil diupdate"})
	}
}

func createBooking(u domain.CateringUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var b domain.Booking
		if err := c.ShouldBindJSON(&b); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if err := u.CreateBooking(&b); err != nil {
			// Menampilkan error jika jadwal bentrok
			c.JSON(409, gin.H{"error": err.Error()})
			return
		}

		// === TRIGGER NOTIFIKASI WEBSOCKET ===
		utils.Hub.Broadcast <- utils.WsMessage{
			Type:    "NEW_CATERING_BOOKING",
			Message: "Ada booking katering baru dari " + b.CustomerName,
		}

		c.JSON(201, gin.H{"message": "Booking berhasil diajukan"})
	}
}

func getBookings(u domain.CateringUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := u.GetAllBookings()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": res})
	}
}

func approveBooking(u domain.CateringUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := u.ApproveBooking(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Booking disetujui"})
	}
}

func rejectBooking(u domain.CateringUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := u.RejectBooking(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Booking ditolak"})
	}
}

func getCaterings(u domain.CateringUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, _ := u.GetAllCaterings()
		c.JSON(200, gin.H{"data": res})
	}
}

// func getCateringPackageByID(u domain.CateringUsecase) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		id, _ := strconv.Atoi(c.Param("id"))
// 		res, err := u.GetCateringByID(uint(id))
// 		if err != nil {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Data katering tidak ditemukan"})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"data": res})
// 	}
// }

// === GANTI getCateringPackageByID DENGAN FUNGSI PINTAR getCateringDetail ===
func getCateringDetail(u domain.CateringUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("param")

		// 1. Deteksi angka (ID untuk Admin)
		if id, err := strconv.Atoi(param); err == nil {
			res, err := u.GetCateringByID(uint(id))
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Data katering tidak ditemukan (ID)"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"data": res})
			return
		}

		// 2. Deteksi huruf (Slug untuk Publik)
		res, err := u.GetCateringBySlug(param)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data katering tidak ditemukan (Slug)"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": res})
	}
}

func createCateringPackage(u domain.CateringUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.PostForm("name")
		description := c.PostForm("description")

		var images []domain.CateringImage

		// Deteksi multi-upload
		form, err := c.MultipartForm()
		if err == nil {
			files := form.File["images"]
			for _, file := range files {
				url, errUp := utils.UploadToCleverCloud(file)
				if errUp == nil {
					images = append(images, domain.CateringImage{ImageURL: url})
				}
			}
		}

		catering := domain.Catering{Name: name, Description: description, Images: images}
		if err := u.CreateCatering(&catering); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"message": "Katering ditambahkan"})
	}
}

func updateCateringPackage(u domain.CateringUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		name := c.PostForm("name")
		description := c.PostForm("description")

		var images []domain.CateringImage
		form, err := c.MultipartForm()
		if err == nil {
			files := form.File["images"]
			for _, file := range files {
				url, errUp := utils.UploadToCleverCloud(file)
				if errUp == nil {
					images = append(images, domain.CateringImage{ImageURL: url})
				}
			}
		}

		catering := domain.Catering{Name: name, Description: description, Images: images}
		if err := u.UpdateCatering(uint(id), &catering); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Katering berhasil diupdate"})
	}
}

func deleteCateringPackage(u domain.CateringUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := u.DeleteCatering(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Katering berhasil dihapus"})
	}
}

// Lalu buat handlernya di bawah:
func getApprovedMomentBookings(u domain.MomentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		res, err := u.GetApprovedBookingsByMomentID(uint(id))
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": res})
	}
}

// === HANDLER BARU: EXPORT TO EXCEL ===
func exportCateringBookings(u domain.CateringUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil semua data dari database
		bookings, err := u.GetAllBookings()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data booking"})
			return
		}

		// 2. Buat File Excel baru di Memori
		f := excelize.NewFile()
		defer func() {
			if err := f.Close(); err != nil {
				// Abaikan error saat menutup memori
			}
		}()

		sheet := "Sheet1"
		f.SetSheetName("Sheet1", "Rekap Katering")
		sheet = "Rekap Katering"

		// 3. Menulis Header Tabel
		f.SetCellValue(sheet, "A1", "ID Booking")
		f.SetCellValue(sheet, "B1", "Nama Pelanggan")
		f.SetCellValue(sheet, "C1", "No Telepon")
		f.SetCellValue(sheet, "D1", "Nama Paket")
		f.SetCellValue(sheet, "E1", "Detail Acara")
		f.SetCellValue(sheet, "F1", "Jumlah Pax")
		f.SetCellValue(sheet, "G1", "Waktu Mulai")
		f.SetCellValue(sheet, "H1", "Waktu Selesai")
		f.SetCellValue(sheet, "I1", "Status")

		// Styling Header (Latar Belakang Orange, Teks Putih Bold)
		style, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
			Fill: excelize.Fill{Type: "pattern", Color: []string{"F97316"}, Pattern: 1}, 
		})
		f.SetRowStyle(sheet, 1, 1, style)

		// Atur lebar kolom agar rapi
		f.SetColWidth(sheet, "A", "A", 12)
		f.SetColWidth(sheet, "B", "C", 20)
		f.SetColWidth(sheet, "D", "E", 30)
		f.SetColWidth(sheet, "G", "H", 20)

		// 4. Masukkan Data ke Excel (Hanya yang APPROVED)
		row := 2
		for _, b := range bookings {
			if b.Status != "APPROVED" {
				continue // Lewati yang berstatus PENDING / REJECTED
			}

			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("BOK-CAT-%d", b.ID))
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), b.CustomerName)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), b.Phone)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), b.Catering.Name)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), b.Description)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), b.MemberCount)
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), b.BookingDate.Format("02 Jan 2006 15:04"))
			f.SetCellValue(sheet, fmt.Sprintf("H%d", row), b.BookingEndDate.Format("02 Jan 2006 15:04"))
			f.SetCellValue(sheet, fmt.Sprintf("I%d", row), b.Status)
			row++
		}

		// 5. Ubah data Excel ke wujud Bytes Buffer
		var buf bytes.Buffer
		if err := f.Write(&buf); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyusun file Excel"})
			return
		}

		// 6. Set Header HTTP untuk memaksa Browser mendownload file (Streaming)
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", "attachment; filename=Rekap_Booking_Katering_Approved.xlsx")
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
	}
}