package main

import (
	"go-rest/config"
	"go-rest/controllers"
	"go-rest/utils"
	"log"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// koneksi database
	config.ConnectDatabase()

	// koneksi redis
	config.InitRedis()

	// router
	r := setupRouter()

	// jalankan server
	r.Run(":8080")
}

func setupRouter() *gin.Engine {

	r := gin.Default()

	r.Static("/uploads", "./uploads")

	// kelompok middleware auth
	auth := r.Group("/")
	auth.Use(utils.AuthMiddleware())
	auth.GET("/decode", controllers.Decode)
	auth.PUT("/password", controllers.ChangePassword)
	auth.PUT("/update", controllers.Update_user)
	auth.POST("/image", controllers.Update_image)
	auth.GET("/logout", controllers.Logout)

	// product
	r.GET("/api/products", controllers.Index_product)
	r.GET("/api/product/:id", controllers.Show)
	r.POST("/api/product", controllers.Create)
	// r.PUT("/api/product/:id", controllers.Update_product)
	r.DELETE("/api/product/:id", controllers.Delete)

	// events
	r.POST("/events", controllers.Index_event)
	r.GET("/events/:id", controllers.Get_event)

	// chapter
	r.GET("/chapter", controllers.Index_chapter)

	// auth
	r.POST("/login", controllers.Login)
	// r.GET("/decode", controllers.Decode)

	r.POST("/forgot", controllers.ForgotPassword)
	r.POST("/otp", controllers.RegOTP)
	r.POST("/register", controllers.Register)
	// r.GET("/verify/:id/:otp", controllers.Verify)
	r.POST("/verify", controllers.Verify)

	// redis test
	r.GET("/redis", controllers.TestRedis)
	r.GET("/cekredis", controllers.DebugSessions)
	r.GET("/city", controllers.Index_city)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "GO REST API",
			"status":  "running",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "GO REST API",
			"status":  "running",
		})
	})

	return r
}
