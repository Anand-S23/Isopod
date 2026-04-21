package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Anand-S23/isopod/internal/config"
	"github.com/Anand-S23/isopod/internal/controller"
	"github.com/Anand-S23/isopod/internal/database"
	"github.com/Anand-S23/isopod/internal/router"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: .env not loaded (%v); using existing environment variables only", err)
	}

	env, err := config.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	ctxTimeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	db, err := database.InitDB(ctx, env.DATABASE_URL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	controller := controller.NewController(ctx, env.PRODUCTION)
	baseRouter := router.NewRouter(controller)
	router := router.NewCorsRouter(baseRouter, "http://localhost:3000") // TODO: allowedOrigin from env

	log.Println("isopod backend running on port: ", env.PORT)
	log.Fatal(http.ListenAndServe(":"+env.PORT, router))
}
