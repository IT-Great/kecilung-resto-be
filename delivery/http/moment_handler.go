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

func RegisterMomentHandlers(rg *gin.RouterGroup, u domain.MomentUsecase) {
	moment := rg.Group("/moments")
	{
		moment.GET("/packages", getMoments(u))
		// moment.GET("/packages/:id", getMomentByID(u))
		// UBAH BARIS INI: Gunakan :param alih-alih :id untuk getMomentDetail
		moment.GET("/packages/:param", getMomentDetail(u))
		moment.POST("/packages", createMomentPackage(u))
		moment.PUT("/packages/:id", updateMomentPackage(u))
		moment.DELETE("/packages/:id", deleteMomentPackage(u))

		moment.GET("/bookings", getMomentBookings(u))
		moment.POST("/bookings", createMomentBooking(u))
		moment.PUT("/bookings/:id/approve", approveMomentBooking(u))
		moment.PUT("/bookings/:id/reject", rejectMomentBooking(u))
		moment.GET("/packages/:id/bookings", getApprovedMomentBookings(u))
	}
}

func createMomentBooking(u domain.MomentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var b domain.MomentBooking
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
			Type:    "NEW_MOMENT_BOOKING",
			Message: "Ada booking moment baru dari " + b.CustomerName,
		}

		c.JSON(201, gin.H{"message": "Booking berhasil diajukan"})
	}
}

func getMomentBookings(u domain.MomentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := u.GetAllBookings()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": res})
	}
}

func approveMomentBooking(u domain.MomentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := u.ApproveBooking(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Booking disetujui"})
	}
}

func rejectMomentBooking(u domain.MomentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := u.RejectBooking(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Booking ditolak"})
	}
}

func getMoments(u domain.MomentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, _ := u.GetAllMoments()
		c.JSON(200, gin.H{"data": res})
	}
}

// func getMomentByID(u domain.MomentUsecase) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		id, _ := strconv.Atoi(c.Param("id"))
// 		res, err := u.GetMomentByID(uint(id))
// 		if err != nil {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Data moment tidak ditemukan"})
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"data": res})
// 	}
// }

// === GANTI getMomentByID DENGAN FUNGSI PINTAR getMomentDetail ===
func getMomentDetail(u domain.MomentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("param")

		// 1. Coba deteksi jika input berupa ID (Angka) -> Dipakai oleh Admin Panel
		if id, err := strconv.Atoi(param); err == nil {
			res, err := u.GetMomentByID(uint(id))
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Data moment tidak ditemukan (ID)"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"data": res})
			return
		}

		// 2. Jika bukan angka, anggap sebagai Slug (Huruf) -> Dipakai oleh Frontend Nuxt Publik
		res, err := u.GetMomentBySlug(param)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data moment tidak ditemukan (Slug)"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": res})
	}
}

func createMomentPackage(u domain.MomentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.PostForm("name")
		description := c.PostForm("description")

		var images []domain.MomentImage

		// Deteksi multi-upload
		form, err := c.MultipartForm()
		if err == nil {
			files := form.File["images"]
			for _, file := range files {
				url, errUp := utils.UploadToCleverCloud(file)
				if errUp == nil {
					images = append(images, domain.MomentImage{ImageURL: url})
				}
			}
		}

		moment := domain.Moment{Name: name, Description: description, Images: images}
		if err := u.CreateMoment(&moment); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, gin.H{"message": "Moment ditambahkan"})
	}
}

func updateMomentPackage(u domain.MomentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		name := c.PostForm("name")
		description := c.PostForm("description")

		var images []domain.MomentImage
		form, err := c.MultipartForm()
		if err == nil {
			files := form.File["images"]
			for _, file := range files {
				url, errUp := utils.UploadToCleverCloud(file)
				if errUp == nil {
					images = append(images, domain.MomentImage{ImageURL: url})
				}
			}
		}

		moment := domain.Moment{Name: name, Description: description, Images: images}
		if err := u.UpdateMoment(uint(id), &moment); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Moment berhasil diupdate"})
	}
}

func deleteMomentPackage(u domain.MomentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		if err := u.DeleteMoment(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Moment berhasil dihapus"})
	}
}

// === HANDLER BARU: EXPORT TO EXCEL ===
func exportMomentBookings(u domain.MomentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		bookings, err := u.GetAllBookings()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data booking"})
			return
		}

		f := excelize.NewFile()
		defer func() { _ = f.Close() }()

		sheet := "Sheet1"
		f.SetSheetName("Sheet1", "Rekap Moment")
		sheet = "Rekap Moment"

		f.SetCellValue(sheet, "A1", "ID Booking")
		f.SetCellValue(sheet, "B1", "Nama Pelanggan")
		f.SetCellValue(sheet, "C1", "No Telepon")
		f.SetCellValue(sheet, "D1", "Nama Paket")
		f.SetCellValue(sheet, "E1", "Detail Acara")
		f.SetCellValue(sheet, "F1", "Jumlah Pax")
		f.SetCellValue(sheet, "G1", "Waktu Mulai")
		f.SetCellValue(sheet, "H1", "Waktu Selesai")
		f.SetCellValue(sheet, "I1", "Status")

		style, _ := f.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true, Color: "FFFFFF"},
			Fill: excelize.Fill{Type: "pattern", Color: []string{"F97316"}, Pattern: 1}, 
		})
		f.SetRowStyle(sheet, 1, 1, style)
		f.SetColWidth(sheet, "A", "A", 12)
		f.SetColWidth(sheet, "B", "C", 20)
		f.SetColWidth(sheet, "D", "E", 30)
		f.SetColWidth(sheet, "G", "H", 20)

		row := 2
		for _, b := range bookings {
			if b.Status != "APPROVED" {
				continue
			}
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("BOK-MOM-%d", b.ID))
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), b.CustomerName)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), b.Phone)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), b.Moment.Name)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), b.Description)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), b.MemberCount)
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), b.BookingDate.Format("02 Jan 2006 15:04"))
			f.SetCellValue(sheet, fmt.Sprintf("H%d", row), b.BookingEndDate.Format("02 Jan 2006 15:04"))
			f.SetCellValue(sheet, fmt.Sprintf("I%d", row), b.Status)
			row++
		}

		var buf bytes.Buffer
		if err := f.Write(&buf); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyusun file Excel"})
			return
		}

		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", "attachment; filename=Rekap_Booking_Moment_Approved.xlsx")
		c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
	}
}
