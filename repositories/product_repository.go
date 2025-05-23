package repositories

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/bobbypratama97/product-rest-api/models"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func GetProducts(sortParam string, page int, limit int) ([]models.Product, models.MetaData, error) {
	key := "products:erajaya"

	//fetch data from redis first
	cached, err := RedisClient.Get(Ctx, key).Result()
	if err == nil {
		var cachedProducts []models.Product
		if err := json.Unmarshal([]byte(cached), &cachedProducts); err == nil {
			paginated := sortAndPaginate(cachedProducts, sortParam, page, limit)
			meta := models.MetaData{
				Page:  page,
				Limit: limit,
				Total: int64(len(cachedProducts)),
				TotalPages: int(math.Ceil(float64(len(cachedProducts)) / float64(limit))),
				Cache: true,
			}
			return paginated, meta , nil
		}
	}

	// if redis cache is empty, fetch from DB
	var products []models.Product
	if err := DB.Order("created_at DESC").Find(&products).Error; err != nil {
		return nil, models.MetaData{}, err
	}

	// insert data to redis using simple goroutine to avoid blocking the main process
	go func(p []models.Product) {
		if data, err := json.Marshal(p); err == nil {
			RedisClient.Set(Ctx, key, data, 15*time.Minute)
		}
	}(products)

	paginated := sortAndPaginate(products, sortParam, page, limit)
	meta := models.MetaData{
		Page:  page,
		Limit: limit,
		Total: int64(len(products)),
		TotalPages: int(math.Ceil(float64(len(products)) / float64(limit))),
		Cache: false,
	}
	return paginated, meta, nil
}

//sorting and pagination process
func sortAndPaginate(products []models.Product, sortParam string, page int, limit int) []models.Product {
	sortFields := strings.Split(sortParam, ",")

	for _, field := range sortFields {
		field = strings.TrimSpace(field)
		switch field {
			case "price_asc":
				sort.Slice(products, func(i, j int) bool {
					return products[i].Price < products[j].Price
				})
			case "price_desc":
				sort.Slice(products, func(i, j int) bool {
					return products[i].Price > products[j].Price
				})
			case "name_asc":
				sort.Slice(products, func(i, j int) bool {
					return products[i].Name < products[j].Name
				})
			case "name_desc":
				sort.Slice(products, func(i, j int) bool {
					return products[i].Name > products[j].Name
				})
			}
	}

	// pagination
	start := (page - 1) * limit
	if start > len(products) {
		start = len(products)
	}
	end := start + limit
	if end > len(products) {
		end = len(products)
	}

	return products[start:end]
}

func InsertProduct(name string, price float64, quantity int, description string) error {
	product := models.Product{
		Name:        name,
		Price:       price,
		Quantity:    quantity,
		Description: description,
	}
	// upsert by product name, update the data if the product name already exists
	if err := DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		UpdateAll: true,
	}).Create(&product).Error; err != nil {
		return err
	}

	// clear redis cache
	cacheKey := "products:erajaya"
	RedisClient.Del(Ctx, cacheKey)

	// refresh cache with updated data
	var products []models.Product
	if err := DB.Order("created_at DESC").Find(&products).Error; err == nil {
		if data, err := json.Marshal(products); err == nil {
			RedisClient.Set(Ctx, cacheKey, data, 15*time.Minute)
		}
	}

	return nil
}