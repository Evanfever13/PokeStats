package main

import (
	"PokeAPI/routers"
	"PokeAPI/templates"
	"fmt"
	"log"
	"net/http"
)

// Fonction Principal
func main() {
	//Charge les templates
	templates.Load()

	//Crée un routeur
	mux := routers.MainRouter()

	//Setup du localhost
	port := "8000"
	addr := fmt.Sprintf("localhost:%s", port)
	fmt.Printf("Serveur prêt sur http://%s\n", addr)
	err := http.ListenAndServe(addr, mux)

	//Gestion d'erreur
	if err != nil {
		log.Fatalf("Erreur lancement serveur : %s\n", err.Error())
	}
}
