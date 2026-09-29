package http

import (
	"github.com/gin-gonic/gin"
	"github.com/seanalden-great/kecilung-resto-be/domain"
	"github.com/seanalden-great/kecilung-resto-be/utils"
)

// PERUBAHAN 1: Tambahkan catUsecase domain.CateringUsecase di parameter ini
func RegisterHandlers(r *gin.Engine, cu domain.CategoryUsecase, mu domain.MenuUsecase, catUsecase domain.CateringUsecase, momentUsecase domain.MomentUsecase, articleUsecase domain.ArticleUsecase, contactUsecase domain.ContactUsecase, authUsecase domain.AuthUsecase) {

	// === JALANKAN HUB WEBSOCKET DI BACKGROUND (GOROUTINE) ===
	go utils.Hub.Run()
	
	api := r.Group("/api")

	// === ENDPOINT WEBSOCKET ===
	api.GET("/ws", utils.ServeWS)

	RegisterCategoryHandlers(api, cu)
	RegisterMenuHandlers(api, mu)
	RegisterCateringHandlers(api, catUsecase)
	RegisterMomentHandlers(api, momentUsecase)
	RegisterArticleHandlers(api, articleUsecase)
	RegisterContactHandlers(api, contactUsecase)
	RegisterAuthHandlers(api, authUsecase)
}
