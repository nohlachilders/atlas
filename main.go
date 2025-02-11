package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("serving...")
	run()
}

func run() {
	mux := http.ServeMux{}
	server := http.Server{
		Addr:    ":8080",
		Handler: &mux,
	}
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	log.Fatal(server.ListenAndServe())
}

type Config struct {
	Port string
}
