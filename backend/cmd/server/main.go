package main

import (
	"log"
	"net/http"
	"os"

	"xwitter/internal/auth"
	"xwitter/internal/controller"
	"xwitter/internal/db"
	"xwitter/internal/repository"
)

func main() {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-only-secret"
	}

	origin := os.Getenv("CLIENT_ORIGIN")
	if origin == "" {
		origin = "http://localhost:8080"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	pool := db.Connect()
	defer pool.Close()

	ctrl := controller.New(
		repository.NewUserRepository(pool),
		repository.NewPostRepository(pool),
		auth.NewService(jwtSecret),
	)

	log.Printf("Xwitter API listening on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, controller.CORS(origin, ctrl.Routes())); err != nil {
		log.Fatal(err)
	}
}
