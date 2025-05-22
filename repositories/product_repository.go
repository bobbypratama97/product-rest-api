package repositories

import (
	"fmt"
	"os"
	"strings"

	"github.com/bobbypratama97/product-rest-api/models"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	//load env variables
	dbName := os.Getenv("DB_DATABASE_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", dbUsername, dbPassword, dbHost, dbName)
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

func GetProducts(sortParam string,page int, limit int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64
	query := DB.Model(&models.Product{})

	sortFields := strings.Split(sortParam, ",")
	appliedSorts := make(map[string]bool)

	//list of valid sorting options
	sortOptions := map[string]struct {
		Column string
		Order  string
	}{
		"price_asc":  {"price", "ASC"},
		"price_desc": {"price", "DESC"},
		"name_asc":   {"name", "ASC"},
		"name_desc":  {"name", "DESC"},
	}

	//iterate through sorting options
	for _, field := range sortFields {
		field = strings.TrimSpace(field) //remove any unwanted empty space
		if opt, ok := sortOptions[field]; ok {
			// we can only have one sort option per column, so we can't sort same column with 2 different options
			if appliedSorts[opt.Column] {
				continue
			}
			//applied the sorting
			query = query.Order(fmt.Sprintf("%s %s", opt.Column, opt.Order))
			//temp variable to marks that the sorting on that column has been applied
			appliedSorts[opt.Column] = true
		}
	}

	if len(appliedSorts) == 0 {
		query = query.Order("created_at DESC")
	}

	query.Count(&total)

	// Pagination
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}
	return products, total,nil
}

func InsertProduct(name string, price float64, quantity int, description string) error {
	product := models.Product{
		Name:        name,
		Price:       price,
		Quantity:    quantity,
		Description: description,
	}
	if err := DB.Create(&product).Error; err != nil {
		return err
	}
	return nil
}