package repositories

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/bobbypratama97/product-rest-api/models"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() error {
	//load env variables
	dbName := os.Getenv("DB_DATABASE_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", dbUsername, dbPassword, dbHost, dbName)
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	return nil
}

func GetProducts() ([]models.Product, error) {
	rows, err := db.Query("SELECT id, name, price, description, quantity, created_at, updated_at FROM products ORDER BY CREATED_AT DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Description, &p.Quantity,&p.CreatedAt,&p.UpdatedAt)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		//date formatting
		p.CreatedAtStr = p.CreatedAt.Format("2006-01-02 15:04:05")
		p.UpdatedAtStr = p.UpdatedAt.Format("2006-01-02 15:04:05")
		products = append(products, p)
	}

	return products, nil
}

func InsertProduct(name string, price float64, quantity int, description string) error {
	_,err := db.Exec("INSERT INTO products (name, price, description, quantity) VALUES (?, ?, ?, ?)", name, price, description, quantity)
	if err != nil {
		return err
	}
	return nil
}
