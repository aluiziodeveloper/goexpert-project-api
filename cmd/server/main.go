package main

import (
	"github.com/aluiziodeveloper/goexpert-project-api/configs"
	"github.com/aluiziodeveloper/goexpert-project-api/internal/entity"
	"github.com/aluiziodeveloper/goexpert-project-api/internal/infra/database"
	"github.com/aluiziodeveloper/goexpert-project-api/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func main() {
	config, err := configs.LoadConfig(".")
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
	// Middlewares do Go-Chi Router
	//r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", config.TokenAuth)) // Disponibilizar o token para os handlers
	r.Use(middleware.WithValue("JwtExpiresIn", config.JwtExpiresIn))
	// Middleware personalizado
	r.Use(LogRequest)

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/", productHandler.GetProducts)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB)

	r.Post("/users", userHandler.Create)
	r.Post("/users/generate_token", userHandler.GetJWT)

	err = http.ListenAndServe(":8000", r)
	if err != nil {
		panic(err)
	}
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
