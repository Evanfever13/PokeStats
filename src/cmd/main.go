package main

import (
	"PokeAPI/routers"
	"PokeAPI/templates"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	templates.Load()
	mux := routers.MainRouter()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf("localhost:%s", port)
	fmt.Printf("Serveur prêt sur http://%s\n", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatalf("Erreur lancement serveur : %s\n", err.Error())
	}
}
