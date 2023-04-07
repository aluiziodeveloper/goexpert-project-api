package main

import (
	"github.com/aluiziodeveloper/goexpert-project-api/configs"
	"github.com/aluiziodeveloper/goexpert-project-api/internal/entity"
	"github.com/aluiziodeveloper/goexpert-project-api/internal/infra/database"
	"github.com/aluiziodeveloper/goexpert-project-api/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
)

func main() {
	_, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&entity.Product{}, &entity.User{})
	if err != nil {
		panic(err)
	}

	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/products", productHandler.CreateProduct)
	r.Get("/products", productHandler.GetProducts)
	r.Get("/products/{id}", productHandler.GetProduct)
	r.Put("/products/{id}", productHandler.UpdateProduct)
	r.Delete("/products/{id}", productHandler.DeleteProduct)

	err = http.ListenAndServe(":8000", r)
	if err != nil {
		panic(err)
	}
}
