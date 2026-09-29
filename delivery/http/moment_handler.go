package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"github.com/seanalden-great/kecilung-resto-be/utils"
)

func RegisterMomentHandlers(rg *gin.RouterGroup, u domain.MomentUsecase) {
	moment := rg.Group("/moments")
	{
		moment.GET("/packages", getMoments(u))
		moment.GET("/packages/:id", getMomentByID(u))
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

func getMomentByID(u domain.MomentUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		res, err := u.GetMomentByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Data moment tidak ditemukan"})
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
