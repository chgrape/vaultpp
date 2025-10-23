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
	"github.com/chgrape/vaultpp/internal/vault"
	"github.com/hashicorp/vault/api"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf(".env file not found: %v", err)
		return
	}

	client, err := api.NewClient(&api.Config{
		Address: "http://127.0.0.1:8200",
	})
	if err != nil {
		log.Fatalf("Connection to Vault couldn't be established: %v", err)
		return
	}

	err = vault.AuthVault(client)
	if err != nil {
		log.Fatalf("Error during Vault authentication: %v", err)
		return
	}

	vaultService := &vault.VaultService{Client: client}
	authVault, err := middleware.NewVaultProvider(vaultService)
	if err != nil {
		log.Fatalf("Error getting JWT key from Vault: %v", err)
		return
	}

	pass, err := vaultService.FetchSecret("db", "postgres_pass")
	if err != nil {
		log.Fatalf("Couldn't get DB credentials from Vault: %v", err)
		return
	}

	cfg := db.Config{
		Host: os.Getenv("POSTGRES_HOST"),
		User: os.Getenv("POSTGRES_USER"),
		Pass: pass,
		Port: os.Getenv("POSTGRES_PORT"),
		DB:   os.Getenv("POSTGRES_DB"),
	}

	pool, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("Connection to database couldn't be established")
	}

	userRepo := repository.UserRepository{DB: pool}
	userService := service.UserService{Repo: &userRepo}
	userHandler := handler.UserHandler{Service: &userService, JWTProvider: authVault}

	itemRepo := repository.ItemRepository{DB: pool}
	itemService := service.ItemService{Repo: &itemRepo}
	itemHandler := handler.ItemHandler{Service: &itemService}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /items", authVault.AuthMiddleware(itemHandler.HandleGetItems))
	mux.HandleFunc("GET /item", itemHandler.HandleGetItem)
	mux.HandleFunc("POST /item", itemHandler.HandleCreateItems)
	mux.HandleFunc("PUT /item", itemHandler.HandleUpdateItem)
	mux.HandleFunc("DELETE /item", itemHandler.HandleDeleteItem)

	mux.HandleFunc("POST /register", userHandler.Register)
	mux.HandleFunc("POST /login", userHandler.Login)
	mux.HandleFunc("GET /me", authVault.AuthMiddleware(userHandler.Details))

	err = http.ListenAndServe(":3333", mux)
	if err != nil {
		panic("Couldn't start HTTP server")
	}
}
