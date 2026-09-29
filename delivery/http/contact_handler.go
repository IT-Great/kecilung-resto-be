package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"github.com/seanalden-great/kecilung-resto-be/utils"
)

func RegisterContactHandlers(rg *gin.RouterGroup, u domain.ContactUsecase) {
	// === TAMBAHKAN ROUTES CONTACT DI SINI ===
	contact := rg.Group("/contacts")
	{
		contact.POST("", createContact(u))
		contact.GET("", getContacts(u))            // Opsional untuk Admin
		contact.PUT("/:id/reply", replyContact(u)) // <--- TAMBAH BARIS INI
	}
}

// ====================================================================
// === HANDLERS CONTACT (TAMBAHKAN KODE INI DI BAGIAN PALING BAWAH) ===
// ====================================================================

func createContact(u domain.ContactUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var contactData domain.ContactUs
		if err := c.ShouldBindJSON(&contactData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := u.Create(&contactData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// === TRIGGER NOTIFIKASI WEBSOCKET ===
		utils.Hub.Broadcast <- utils.WsMessage{
			Type:    "NEW_CONTACT_MSG",
			Message: "Pesan baru dari " + contactData.FullName,
		}
		
		c.JSON(http.StatusCreated, gin.H{"message": "Pesan Anda berhasil dikirim", "data": contactData})
	}
}

func getContacts(u domain.ContactUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := u.GetAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": res})
	}
}

// Tambahkan handler baru ini di bawah getContacts:
func replyContact(u domain.ContactUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))

		// Kita ambil text balasan dari JSON body
		var payload struct {
			Reply string `json:"reply"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format data tidak valid"})
			return
		}

		if err := u.ReplyMessage(uint(id), payload.Reply); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Pesan berhasil dibalas"})
	}
}