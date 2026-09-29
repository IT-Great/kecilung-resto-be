package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"github.com/seanalden-great/kecilung-resto-be/utils"
)

func RegisterCateringHandlers(rg *gin.RouterGroup, u domain.CateringUsecase) {
	// Routes Catering
	catering := rg.Group("/catering")
	{
		catering.GET("/greeting", getGreeting(u))
		catering.PUT("/greeting", updateGreeting(u))
		catering.GET("/packages", getCaterings(u))
		catering.GET("/packages/:id", getCateringPackageByID(u))
		catering.POST("/packages", createCateringPackage(u))
		// Tambahkan 2 baris ini di dalam RegisterHandlers -> catering := api.Group("/catering")
		catering.PUT("/packages/:id", updateCateringPackage(u))
		catering.DELETE("/packages/:id", deleteCateringPackage(u))
		catering.GET("/bookings", getBookings(u))
		catering.POST("/bookings", createBooking(u))
		catering.PUT("/bookings/:id/approve", approveBooking(u))
		catering.PUT("/bookings/:id/reject", rejectBooking(u))
		catering.GET("/packages/:id/bookings", getApprovedCateringBookings(u))
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

func getCateringPackageByID(u domain.CateringUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		res, err := u.GetCateringByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data katering tidak ditemukan"})
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