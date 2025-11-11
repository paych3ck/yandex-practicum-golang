package main

import (
	"log"
	"os"

	"sprint6/internal/server"
)

func main() {
	logger := log.New(os.Stdout, "http ", log.LstdFlags|log.Lshortfile)
	srv := server.New(logger)

	if err := srv.HTTP.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}
