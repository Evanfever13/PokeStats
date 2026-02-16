package controllers

import (
	"PokeAPI/helpers"
	"PokeAPI/models"
	"PokeAPI/services"
	"PokeAPI/templates"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

/*=========================================================
  CONTROLEUR : Gestion des données à demender aux Endpoint
  =========================================================*/

// Page Home
func HomeDisplay(w http.ResponseWriter, r *http.Request) {
	//Pas de Requete à faire...
	templates.RenderTemplate(w, r, "Home", nil)
}

// Page About
func AboutDisplay(w http.ResponseWriter, r *http.Request) {
	//Pas de Requete à faire..
	templates.RenderTemplate(w, r, "About", nil)
}

// Page Pokemon
func PokemonDisplay(w http.ResponseWriter, r *http.Request) {
	//Recuperation de la Query
	query := r.URL.Query()
	search := strings.TrimSpace(query.Get("search"))
	selectedType := strings.TrimSpace(query.Get("type"))

	//Calcul de l'offset necessaire
	offset := 0
	if vals, ok := query["offset"]; ok && len(vals) > 0 {
		fmt.Sscanf(vals[0], "%d", &offset)
		if offset < 0 {
			offset = 0
		}
	}

	//Definition des variables
	var data *models.PokemonResponse
	var status int
	var err error

	//Recherche (elle est prioritaire par rapport au tri)
	if search != "" {
		if _, convErr := strconv.Atoi(search); convErr == nil {
			data, status, err = services.PokemonByNameOrID(search)
		} else {
			data, status, err = services.PokemonByNameOrID(search)
		}

		//Tri
	} else {
		if selectedType != "" {
			data, status, err = services.PokemonByType(strings.ToLower(selectedType), offset)
		} else {
			data, status, err = services.PokemonService(offset)
		}
	}

	//Verfication que le server renvoit un 200 OK
	if status != http.StatusOK || err != nil {
		helpers.RedirectToError(w, r, status, "Erreur lors de la récupération des données")
		if err != nil {
			log.Print(err.Error())
		}
		return
	}

	//Systeme de pagination//
	prev := offset - 20
	if prev < 0 {
		prev = 0
	}
	next := offset + 20

	//(Source : StackOverflow)
	//Sert à stocker les resultat, la pagination, recherche, ect...
	payload := map[string]interface{}{
		"Results":      data.Results,
		"Offset":       offset,
		"PrevOffset":   prev,
		"NextOffset":   next,
		"Search":       search,
		"SelectedType": selectedType,
	}

	templates.RenderTemplate(w, r, "Pokemon", payload)
}

// Page Pokemon Details
func PokemonDetailsDisplay(w http.ResponseWriter, r *http.Request) {
	//Recuperation de la Query
	query := r.URL.Query()
	id := strings.TrimSpace(query.Get("id"))
	if id == "" {
		helpers.RedirectToError(w, r, http.StatusBadRequest, "Missing pokemon id")
		return
	}

	//Recuperation des Details du Pokemon
	details, status, err := services.GetPokemonDetails(id)
	if status != http.StatusOK || err != nil {
		helpers.RedirectToError(w, r, status, "Erreur lors de la récupération des données")
		if err != nil {
			log.Print(err.Error())
		}
		return
	}

	templates.RenderTemplate(w, r, "PokemonDetails", details)
}

// Page Moves
func MovesDisplay(w http.ResponseWriter, r *http.Request) {
	//Recuperation de la Query
	query := r.URL.Query()
	search := strings.TrimSpace(query.Get("search"))

	//Calcul de l'offset necessaire
	offset := 0
	if vals, ok := query["offset"]; ok && len(vals) > 0 {
		fmt.Sscanf(vals[0], "%d", &offset)
		if offset < 0 {
			offset = 0
		}
	}

	//Definition des Variables
	var data *models.MoveResponse
	var status int
	var err error

	//Recherche
	if search != "" {
		data, status, err = services.MoveByNameOrID(strings.ToLower(search))
	} else {
		data, status, err = services.MoveService(offset)
	}

	//Verfication que le server renvoit un 200 OK
	if status != http.StatusOK || err != nil {
		helpers.RedirectToError(w, r, status, "Erreur lors de la récupération des données")
		if err != nil {
			log.Print(err.Error())
		}
		return
	}

	//Systeme de Pagination
	prev := offset - 20
	if prev < 0 {
		prev = 0
	}
	next := offset + 20

	//(Source: StackOverflow)
	//Sert à stocker les resultat, la pagination, recherche, ect...
	payload := map[string]interface{}{
		"Results":    data.Results,
		"Offset":     offset,
		"PrevOffset": prev,
		"NextOffset": next,
		"Search":     search,
	}

	templates.RenderTemplate(w, r, "Moves", payload)
}

// Page Moves Details
func MovesDetailsDisplay(w http.ResponseWriter, r *http.Request) {
	//Recuperation de la Query
	query := r.URL.Query()
	id := strings.TrimSpace(query.Get("id"))
	if id == "" {
		helpers.RedirectToError(w, r, http.StatusBadRequest, "Missing move id")
		return
	}

	//Recuperation des Details de l'Attaque
	details, status, err := services.GetMoveDetails(id)
	if status != http.StatusOK || err != nil {
		helpers.RedirectToError(w, r, status, "Erreur lors de la récupération des données")
		if err != nil {
			log.Print(err.Error())
		}
		return
	}

	templates.RenderTemplate(w, r, "MoveDetails", details)
}

// Page Teams (Favoris)
func TeamsDisplay(w http.ResponseWriter, r *http.Request) {
	//Recuperation de la Quary
	query := r.URL.Query()

	//Remove des Favoris
	removeID := strings.TrimSpace(query.Get("remove"))
	if removeID != "" {
		templates.RemoveTeamID(removeID)
		http.Redirect(w, r, "/teams", http.StatusSeeOther)
		return
	}

	//Add des Favoris
	id := strings.TrimSpace(query.Get("add"))
	if id != "" {
		for _, v := range templates.ListTeams {
			if v == id {
				helpers.RedirectToError(w, r, http.StatusBadRequest, "Team already added")
				return
			}
		}
		templates.AddTeamID(id)
	}

	//Definition des variable
	var teamDetails []*models.PokemonDetails

	//Parcours des Favoris et recuperation des données des Pokemons
	for _, Pokemon := range templates.ListTeams {
		if Pokemon == "" {
			continue
		}
		details, status, err := services.TeamsService(Pokemon)

		//Verfication que le server renvoit un 200 OK
		if status != http.StatusOK || err != nil {
			helpers.RedirectToError(w, r, status, "Erreur lors de la récupération des données")
			if err != nil {
				log.Print(err.Error())
			}
			return
		}

		//Ajouts chaque de Pokemon dans la liste
		teamDetails = append(teamDetails, details)
	}

	//(Source: StackOverflow)
	//Sert à ne pas dupliquer le header à chaque nouveau favori...
	payload := map[string]interface{}{
		"Header": true,
		"Teams":  teamDetails,
	}
	templates.RenderTemplate(w, r, "Teams", payload)
}
