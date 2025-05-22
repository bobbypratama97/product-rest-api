# 🛒 Product REST API

A simple RESTful API built with Go (Gin + GORM) to manage product data. Supports product creation and listing (including pagination and sorting)

---

## 🚀 Getting Started

### 1. Clone the Project
```bash
git clone https://github.com/your-username/product-rest-api.git
cd product-rest-api
```

### 2. Create .env file
Copy the env file from .env.example

### 3. Install Required Dependencies
```bash
go mod tidy
```

### 4. Run the Application
```bash
go run main.go
gin --appPort 5000 --all -i run server.go
```
For more convinient way, you can use the gin syntax so the project will auto reload everything there is a changes detected.


## 🧱 Architecture Overview
This project follows a clean architecture pattern with clear separation of responsibilities across every layers.

`📂 /controllers`
Handles HTTP request routing, input validation, and API response formatting.

`📂 /repositories`
Contains logic to query and manipulate the database using GORM. Each functions are isolated for readability, testability, and reusability.

`📂 /models`
Defines the structure of database models and handles things like custom JSON formatting.

`📂 /utilities`
Contains utility/helper functions such as middlewares, validators, and other reusable tools. This design ensures maintainability, scalability, and testability all accross the project.

## 🔗 API Documentation

## 👨‍🚀 Author
Built with 🫶 by Bobby Pratama.
