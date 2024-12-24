package main

import (
	"backend/internal/auth"
	"backend/internal/calendar"
	"backend/internal/db"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	defer db.DB.Close()
	router := gin.Default()

	router.Use(func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    })

	router.POST("/register", auth.Register)
	router.POST("/login", auth.Login)

	authRoutes := router.Group("/")
	authRoutes.Use(auth.AuthMiddleware())
	{
		authRoutes.POST("/calendar/day", calendar.CreateDay)
		authRoutes.GET("/calendar/next", calendar.GetNextDay)
	}

	router.Run(":8181")
}
