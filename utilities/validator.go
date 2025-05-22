package utilities

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func ValidateRequest(ctx *gin.Context, req interface{}) (string ,error) {
	//validate json structure
	if err := ctx.ShouldBindJSON(req); err != nil {
		return "Invalid Request Format!",err
	}
	//validate request body
	if err := Validate.Struct(req); err != nil {
		return "Validation Failed!",err
	}
	return "",nil
}