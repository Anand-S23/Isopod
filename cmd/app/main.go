package main

import (
	"log"
	"net/http"
	"time"
	"context"

	"github.com/Anand-S23/isopod/internal/config"
	"github.com/Anand-S23/isopod/internal/controller"
	"github.com/Anand-S23/isopod/internal/router"
)

func main() {
	env, err := config.LoadEnv()
    if err != nil {
        log.Fatal(err)
    }

    ctxTimeout := 5 * time.Second
    ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
    defer cancel()

    controller := controller.NewController(ctx, env.PRODUCTION)
	baseRouter := router.NewRouter(controller)
	router := router.NewCorsRouter(baseRouter, "http://localhost:3000")

    log.Println("isopod backend running on port: ", env.PORT)
    log.Fatal(http.ListenAndServe(":" + env.PORT, router))
}
