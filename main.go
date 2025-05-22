package main

import (
	"fmt"
	"log"
	"time"

	"github.com/bobbypratama97/product-rest-api/controllers"
	"github.com/bobbypratama97/product-rest-api/repositories"

	"github.com/gin-gonic/gin"
)

func main() {
	// init DB connection
	err := repositories.InitDB()
	if err != nil {
			log.Fatalf("Failed to connect to DB: %v", err)
	}
	// init server
	server := gin.Default()

	server.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"serverTime": time.Now().UTC().Unix(),
			"status":     "Product Rest API ready to Go!",
		})
	})

	server.GET("/products", controllers.GetProducts)
	server.POST("/products", controllers.InsertProduct)

	err = server.Run(":5000")
	if err != nil {
		fmt.Println(err)
	}

}
