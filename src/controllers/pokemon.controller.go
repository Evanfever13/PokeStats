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

func HomeDisplay(w http.ResponseWriter, r *http.Request) {
	templates.RenderTemplate(w, r, "Home", nil)
}

func AboutDisplay(w http.ResponseWriter, r *http.Request) {
	templates.RenderTemplate(w, r, "About", nil)
}

func PokemonDisplay(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	search := strings.TrimSpace(query.Get("search"))
	selectedType := strings.TrimSpace(query.Get("type"))
	offset := 0
	if vals, ok := query["offset"]; ok && len(vals) > 0 {
		fmt.Sscanf(vals[0], "%d", &offset)
		if offset < 0 {
			offset = 0
		}
	}

	var data *models.PokemonResponse
	var status int
	var err error
	if search != "" {
		if _, convErr := strconv.Atoi(search); convErr == nil {
			data, status, err = services.PokemonByNameOrID(search)
		} else {
			data, status, err = services.SearchPokemon(strings.ToLower(search))
		}
	} else {
		if selectedType != "" {
			data, status, err = services.PokemonByType(strings.ToLower(selectedType), offset)
		} else {
			data, status, err = services.PokemonService(offset)
		}
	}
	if status != http.StatusOK || err != nil {
		helpers.RedirectToError(w, r, status, "Erreur lors de la récupération des données")
		if err != nil {
			log.Print(err.Error())
		}
		return
	}
	prev := offset - 20
	if prev < 0 {
		prev = 0
	}
	next := offset + 20
	//Le payload je l'ai trouvé sur StackOverflow
	payload := map[string]interface{}{
		"Results":      data.Results,
		"Offset":       offset,
		"PrevOffset":   prev,
		"NextOffset":   next,
		"Search":       search,
		"Types":        nil,
		"SelectedType": selectedType,
	}
	if types, st, _ := services.GetTypes(); st == http.StatusOK {
		payload["Types"] = types
	}

	templates.RenderTemplate(w, r, "Pokemon", payload)
}

func PokemonDetailsDisplay(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id := strings.TrimSpace(query.Get("id"))
	if id == "" {
		helpers.RedirectToError(w, r, http.StatusBadRequest, "Missing pokemon id")
		return
	}

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

func MovesDisplay(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	search := strings.TrimSpace(query.Get("search"))
	offset := 0
	if vals, ok := query["offset"]; ok && len(vals) > 0 {
		fmt.Sscanf(vals[0], "%d", &offset)
		if offset < 0 {
			offset = 0
		}
	}

	var data *models.MoveResponse
	var status int
	var err error
	if search != "" {
		data, status, err = services.MoveByNameOrID(strings.ToLower(search))
	} else {
		data, status, err = services.MoveService(offset)
	}
	if status != http.StatusOK || err != nil {
		helpers.RedirectToError(w, r, status, "Erreur lors de la récupération des données")
		if err != nil {
			log.Print(err.Error())
		}
		return
	}

	prev := offset - 20
	if prev < 0 {
		prev = 0
	}
	next := offset + 20

	payload := map[string]interface{}{
		"Results":    data.Results,
		"Offset":     offset,
		"PrevOffset": prev,
		"NextOffset": next,
		"Search":     search,
	}
	templates.RenderTemplate(w, r, "Moves", payload)
}

func MovesDetailsDisplay(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id := strings.TrimSpace(query.Get("id"))
	if id == "" {
		helpers.RedirectToError(w, r, http.StatusBadRequest, "Missing move id")
		return
	}

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

func TeamsDisplay(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	removeID := strings.TrimSpace(query.Get("remove"))
	if removeID != "" {
		templates.RemoveTeamID(removeID)
		http.Redirect(w, r, "/teams", http.StatusSeeOther)
		return
	}

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

	var teamDetails []*models.PokemonDetails
	for _, Pokemon := range templates.ListTeams {
		if Pokemon == "" {
			continue
		}
		details, status, err := services.TeamsService(Pokemon)
		if status != http.StatusOK || err != nil {
			helpers.RedirectToError(w, r, status, "Erreur lors de la récupération des données")
			if err != nil {
				log.Print(err.Error())
			}
			return
		}
		teamDetails = append(teamDetails, details)
	}

	payload := map[string]interface{}{
		"Header": true,
		"Teams":  teamDetails,
	}
	templates.RenderTemplate(w, r, "Teams", payload)
}
