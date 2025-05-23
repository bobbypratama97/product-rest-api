package models

import (
	"encoding/json"
	"time"
)

type Product struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"-" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `json:"-" gorm:"column:updated_at;autoUpdateTime"`
}


type ProductRequest struct {
	Name        string  `json:"name" validate:"required"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Description string  `json:"description" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required,min=1"`
}

// customize JSON output
func (p Product) MarshalJSON() ([]byte, error) {
	type productAlias Product
	alias := productAlias(p)
	loc, _ := time.LoadLocation("Asia/Jakarta")

	ordered := struct {
		ID          uint    `json:"id"`
		Name        string  `json:"name"`
		Price       float64 `json:"price"`
		Description string  `json:"description"`
		Quantity    int     `json:"quantity"`
		CreatedAt   string  `json:"created_at"`
		UpdatedAt   string  `json:"updated_at"`
	}{
		ID:          alias.ID,
		Name:        alias.Name,
		Price:       alias.Price,
		Description: alias.Description,
		Quantity:    alias.Quantity,
		CreatedAt:   alias.CreatedAt.In(loc).Format("2006-01-02 15:04:05"),
		UpdatedAt:   alias.UpdatedAt.In(loc).Format("2006-01-02 15:04:05"),
	}

	return json.Marshal(ordered)
}

//unmarshal JSON from redis to format the date
func (p *Product) UnmarshalJSON(data []byte) error {
	type Alias Product
	formattedColumn := &struct {
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &formattedColumn); err != nil {
		return err
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
 createdAt, err := time.ParseInLocation("2006-01-02 15:04:05", formattedColumn.CreatedAt, loc)
	if err != nil {
			return err
	}
	updatedAt, err := time.ParseInLocation("2006-01-02 15:04:05", formattedColumn.UpdatedAt, loc)
	if err != nil {
			return err
	}
	p.CreatedAt = createdAt.UTC()
	p.UpdatedAt = updatedAt.UTC()

	return nil
}

type ProductResponse struct {
	Code    int         `json:"code"`
	Meta    MetaData    `json:"meta"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type MetaData struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int `json:"total_pages"`
	Cache      bool `json:"cache"`
}
