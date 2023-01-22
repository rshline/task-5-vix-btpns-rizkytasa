package router

import (
	// "os"
	// "github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"github.com/rshline/task-5-vix-btpns-rizkytasa/controllers"
	"gorm.io/gorm"
)

//Function to initialize routes
func InitRoutes(db *gorm.DB) *gin.Engine {
	
	// godotenv.Load(".env")
	
	router := gin.Default()

	// proxy := os.Getenv("PROXY")
	// router.SetTrustedProxies([]string{proxy})

	router.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	router.POST("/users/register", controllers.Register)
	router.POST("/users/login", controllers.Login)
	router.PUT("/users/:userId", controllers.UpdateUser)
	router.DELETE("/users/:userId", controllers.DeleteUser)

	router.POST("/photos", controllers.AddPhoto)
	router.GET("/photos", controllers.GetPhoto)
	router.PUT("/:photoId", controllers.UpdatePhoto)
	router.DELETE("/:photoId", controllers.DeletePhoto)

	return router
}