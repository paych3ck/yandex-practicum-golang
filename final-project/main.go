package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/paych3ck/final_project/pkg/api"
	"github.com/paych3ck/final_project/pkg/db"
)

const (
	PORT = 7540
)

func main() {
	if _, err := db.Init("scheduler.db"); err != nil {
		log.Fatalf("unable to initialize database: %v", err)
	}

	api.Init()

	webDir, err := filepath.Abs("web")
	if err != nil {
		log.Fatalf("unable to resolve web directory: %v", err)
	}

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	addr := fmt.Sprintf(":%d", PORT)
	log.Printf("serving %s on %s", webDir, addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
