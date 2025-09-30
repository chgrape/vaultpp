package main

import (
	"log"
	"net/http"
	"os"

	"github.com/chgrape/vaultpp/internal/db"
	"github.com/chgrape/vaultpp/internal/handler"
	"github.com/chgrape/vaultpp/internal/middleware"
	"github.com/chgrape/vaultpp/internal/repository"
	"github.com/chgrape/vaultpp/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(".env file not found")
	}

	cfg := db.Config{
		Host: os.Getenv("POSTGRES_HOST"),
		User: os.Getenv("POSTGRES_USER"),
		Pass: os.Getenv("POSTGRES_PASS"),
		Port: os.Getenv("POSTGRES_PORT"),
		DB:   os.Getenv("POSTGRES_DB"),
	}

	pool, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Connection to database couldn't be established")
	}

	userRepo := repository.UserRepository{DB: pool}
	userService := service.UserService{Repo: &userRepo}
	userHandler := handler.UserHandler{Service: &userService}

	itemRepo := repository.ItemRepository{DB: pool}
	itemService := service.ItemService{Repo: &itemRepo}
	itemHandler := handler.ItemHandler{Service: &itemService}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /items", middleware.AuthMiddleware(itemHandler.HandleGetItems))
	mux.HandleFunc("GET /item", itemHandler.HandleGetItem)
	mux.HandleFunc("POST /item", itemHandler.HandleCreateItems)
	mux.HandleFunc("PUT /item", itemHandler.HandleUpdateItem)
	mux.HandleFunc("DELETE /item", itemHandler.HandleDeleteItem)

	mux.HandleFunc("POST /register", userHandler.Register)
	mux.HandleFunc("POST /login", userHandler.Login)

	err = http.ListenAndServe(":3333", mux)
	if err != nil {
		panic("Couldn't start HTTP server")
	}
}
