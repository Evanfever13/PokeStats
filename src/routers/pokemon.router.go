package routers

import (
	"PokeAPI/controllers"
	"net/http"
)

func MainRouter() *http.ServeMux {
	mainRouter := http.NewServeMux()
	PokeRouter(mainRouter)
	fileServer := http.FileServer(http.Dir("./../../assets"))
	mainRouter.Handle("/static/", http.StripPrefix("/static/", fileServer))
	return mainRouter
}

func PokeRouter(router *http.ServeMux) {
	router.HandleFunc("/", controllers.HomeDisplay)
	router.HandleFunc("/pokemon", controllers.PokemonDisplay)
	router.HandleFunc("/moves", controllers.MovesDisplay)
	router.HandleFunc("/about", controllers.AboutDisplay)
	router.HandleFunc("/teams", controllers.TeamsDisplay)
	router.HandleFunc("/moves/details", controllers.MovesDetailsDisplay)
	router.HandleFunc("/pokemon/details", controllers.PokemonDetailsDisplay)
}
