package main

import (
	"log"
	"net/http"

	"github.com/Anand-S23/isopod/internal/config"
)

func main() {
	env, err := config.LoadEnv()
    if err != nil {
        log.Fatal(err)
    }
	
    log.Println("isopod backend running on port: ", env.PORT);
    http.ListenAndServe(":" + env.PORT, nil)
}

