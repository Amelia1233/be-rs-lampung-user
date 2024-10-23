package routers

import (
	"be-rs-lampung-user/auth/handler"
	"be-rs-lampung-user/auth/repository"
	"be-rs-lampung-user/auth/usecase"
	"be-rs-lampung-user/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

type Routes struct {
	Db *gorm.DB
	R  *gin.Engine
}

func (r Routes) Routers() {
	middleware.Add(r.R, middleware.CORSMiddleware())
	v1 := r.R.Group("user")

	// Inisialisasi repository
	repo := repository.NewAuthRepository(r.Db)

	// Inisialisasi usecase
	authUsecase := usecase.NewAuthUsecase(repo)

	// Inisialisasi handler
	authHandler := handler.NewAuthHandler(authUsecase)

	// Setup rute login
	v1.POST("/login", authHandler.Login)
	v1.GET("/me", middleware.AuthMiddleware(), authHandler.GetUserData)
	v1.POST("/refresh", authHandler.RefreshToken)

	// Inisialisasi admin statis
	if err := authUsecase.InitStaticAdmins(); err != nil {
		log.Printf("Gagal menginisialisasi admin: %v", err)
	}
}
