package main

import (
	"be-rs-lampung-user/connection"
	"be-rs-lampung-user/entity"
	"be-rs-lampung-user/routers"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	r := gin.Default()

	// Konfigurasi CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://88.222.214.98:9001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	db := connection.Connection()

	if !db.Migrator().HasColumn(&entity.User{}, "refresh_token") {
		if err := db.Migrator().AddColumn(&entity.User{}, "refresh_token"); err != nil {
			log.Printf("Gagal menambahkan kolom refresh_token: %v", err)
		}
	}

	if err := db.AutoMigrate(&entity.User{}); err != nil {
		log.Fatalf("Gagal melakukan migrasi database: %v", err)
	}

	routes := routers.Routes{
		Db: db,
		R:  r,
	}

	routes.Routers()

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
