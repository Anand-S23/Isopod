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
	err := godotenv.Load()
	if err != nil {
		log.Printf("warning: .env not loaded (%v); using existing environment variables only", err)
	}

	env, err := config.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	ctxTimeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	db, err := database.InitDB(ctx, env.DB_URI)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := store.NewStore(db, store.NewPgUserRepo(db))
	controller := controller.NewController(store, env, ctx)
	baseRouter := router.NewRouter(controller)
	router := router.NewCorsRouter(baseRouter, env.ALLOWED_ORIGIN)

	log.Println("isopod backend running on port: ", env.PORT)
	log.Fatal(http.ListenAndServe(":"+env.PORT, router))
}
