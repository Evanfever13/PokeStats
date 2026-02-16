package routers

import (
	"PokeAPI/controllers"
	"net/http"
)

/*=========================================================
  ROUTER : Gestion des Routes
  =========================================================*/

// Creation du routeur
func MainRouter() *http.ServeMux {
	//Mise en place des Routeurs
	mainRouter := http.NewServeMux()
	PokeRouter(mainRouter)

	//Gestion des fichiers statique
	fileServer := http.FileServer(http.Dir("./../../assets"))
	mainRouter.Handle("/static/", http.StripPrefix("/static/", fileServer))
	return mainRouter
}

// Relier les Routes et les controlleurs
func PokeRouter(router *http.ServeMux) {
	router.HandleFunc("/", controllers.HomeDisplay)
	router.HandleFunc("/pokemon", controllers.PokemonDisplay)
	router.HandleFunc("/moves", controllers.MovesDisplay)
	router.HandleFunc("/about", controllers.AboutDisplay)
	router.HandleFunc("/teams", controllers.TeamsDisplay)
	router.HandleFunc("/moves/details", controllers.MovesDetailsDisplay)
	router.HandleFunc("/pokemon/details", controllers.PokemonDetailsDisplay)
}
