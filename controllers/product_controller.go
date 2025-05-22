package controllers

import (
	"net/http"

	"github.com/bobbypratama97/product-rest-api/models"
	"github.com/bobbypratama97/product-rest-api/repositories"
	"github.com/bobbypratama97/product-rest-api/utilities"
	"github.com/gin-gonic/gin"
)

func GetProducts(ctx *gin.Context) {
	products, err := repositories.GetProducts()
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to fetch products"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":   http.StatusOK,
		"meta": gin.H{},
		"message": "Successfully fetched products",
		"data": products,
	})
}

func InsertProduct(ctx *gin.Context) {
	var req models.ProductRequest
	// validate request body
	if message,err := utilities.ValidateRequest(ctx, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": message,
			"result": err.Error(),
		})
		return
	}
	err := repositories.InsertProduct(req.Name, req.Price, req.Quantity, req.Description)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"message": "Error inserting product",
			"result": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"code": http.StatusCreated,
		"message": "Successfully created product",
	})

}