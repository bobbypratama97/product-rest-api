package controllers

import (
	"math"
	"net/http"
	"strconv"

	"github.com/bobbypratama97/product-rest-api/models"
	"github.com/bobbypratama97/product-rest-api/repositories"
	"github.com/bobbypratama97/product-rest-api/utilities"
	"github.com/gin-gonic/gin"
)


func GetProducts(ctx *gin.Context) {
	sortParam := ctx.Query("sorting")
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "10")

	// Convert to integers
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	products, total, err := repositories.GetProducts(sortParam,page,limit)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to fetch products"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	response := models.ProductResponse{
		Code: http.StatusOK,
		Meta: models.MetaData{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
		Message: "Successfully fetched products",
		Data:    products,
	}

	ctx.JSON(http.StatusOK, response)
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