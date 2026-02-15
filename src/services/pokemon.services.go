package services

import (
	"PokeAPI/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func PokemonService(offset int) (*models.PokemonResponse, int, error) {

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon?offset=%d", offset)
	request, requestErr := http.NewRequest(http.MethodGet, url, nil)
	if requestErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur initialisation requete - %s", requestErr.Error())
	}

	response, responseErr := client.Do(request)
	if responseErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur envoi requete - %s", responseErr.Error())
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, response.StatusCode, fmt.Errorf("erreur reponse requete - code : %d, status : %s", response.StatusCode, response.Status)
	}

	var listPokemon models.PokemonResponse

	decodeErr := json.NewDecoder(response.Body).Decode(&listPokemon)
	if decodeErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur decodage des données - %s", decodeErr.Error())
	}

	return &listPokemon, response.StatusCode, nil
}

func GetTypes() ([]string, int, error) {
	client := http.Client{Timeout: 5 * time.Second}
	url := "https://pokeapi.co/api/v2/type"
	request, requestErr := http.NewRequest(http.MethodGet, url, nil)
	if requestErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur initialisation requete - %s", requestErr.Error())
	}

	response, responseErr := client.Do(request)
	if responseErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur envoi requete - %s", responseErr.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, response.StatusCode, fmt.Errorf("erreur reponse requete - code : %d, status : %s", response.StatusCode, response.Status)
	}

	var payload struct {
		Results []struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"results"`
	}

	decodeErr := json.NewDecoder(response.Body).Decode(&payload)
	if decodeErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur decodage des données - %s", decodeErr.Error())
	}

	types := make([]string, 0, len(payload.Results))
	for _, t := range payload.Results {
		types = append(types, t.Name)
	}

	return types, http.StatusOK, nil
}

// PokemonByType retourne la liste des pokemons pour un type donné (avec pagination côté client)
func PokemonByType(typeName string, offset int) (*models.PokemonResponse, int, error) {
	client := http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/type/%s", typeName)
	request, requestErr := http.NewRequest(http.MethodGet, url, nil)
	if requestErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur initialisation requete - %s", requestErr.Error())
	}

	response, responseErr := client.Do(request)
	if responseErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur envoi requete - %s", responseErr.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, response.StatusCode, fmt.Errorf("erreur reponse requete - code : %d, status : %s", response.StatusCode, response.Status)
	}

	var payload struct {
		Pokemon []struct {
			Pokemon struct {
				Name string `json:"name"`
				Url  string `json:"url"`
			} `json:"pokemon"`
		} `json:"pokemon"`
	}

	decodeErr := json.NewDecoder(response.Body).Decode(&payload)
	if decodeErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur decodage des données - %s", decodeErr.Error())
	}

	// Convertir en models.Pokemon
	all := make([]models.Pokemon, 0, len(payload.Pokemon))
	for _, p := range payload.Pokemon {
		all = append(all, models.Pokemon{Name: p.Pokemon.Name, Url: p.Pokemon.Url})
	}

	// Pagination côté client (taille de page 20)
	limit := 20
	start := offset
	if start < 0 {
		start = 0
	}
	end := start + limit
	if start > len(all) {
		start = len(all)
	}
	if end > len(all) {
		end = len(all)
	}

	res := &models.PokemonResponse{Results: all[start:end]}
	return res, http.StatusOK, nil
}

func PokemonByNameOrID(name string) (*models.PokemonResponse, int, error) {
	client := http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)
	request, requestErr := http.NewRequest(http.MethodGet, url, nil)
	if requestErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur initialisation requete - %s", requestErr.Error())
	}

	response, responseErr := client.Do(request)
	if responseErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur envoi requete - %s", responseErr.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, response.StatusCode, fmt.Errorf("erreur reponse requete - code : %d, status : %s", response.StatusCode, response.Status)
	}

	var payload struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	decodeErr := json.NewDecoder(response.Body).Decode(&payload)
	if decodeErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur decodage des données - %s", decodeErr.Error())
	}

	p := models.Pokemon{
		Name: payload.Name,
		Url:  fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%d/", payload.ID),
	}
	res := &models.PokemonResponse{Results: []models.Pokemon{p}}
	return res, http.StatusOK, nil
}

func GetPokemonDetails(idOrName string) (*models.PokemonDetails, int, error) {
	client := http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", idOrName)
	request, requestErr := http.NewRequest(http.MethodGet, url, nil)
	if requestErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur initialisation requete - %s", requestErr.Error())
	}

	response, responseErr := client.Do(request)
	if responseErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur envoi requete - %s", responseErr.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, response.StatusCode, fmt.Errorf("erreur reponse requete - code : %d, status : %s", response.StatusCode, response.Status)
	}

	var payload struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Sprites struct {
			Other struct {
				OfficialArtwork struct {
					FrontDefault string `json:"front_default"`
				} `json:"official-artwork"`
			} `json:"other"`
		} `json:"sprites"`
		Types []struct {
			Type struct {
				Name string `json:"name"`
			} `json:"type"`
		} `json:"types"`
		Stats []struct {
			Stat struct {
				Name string `json:"name"`
			} `json:"stat"`
			BaseStat int `json:"base_stat"`
		} `json:"stats"`
		Abilities []struct {
			Ability struct {
				Name string `json:"name"`
			} `json:"ability"`
		} `json:"abilities"`
	}

	decodeErr := json.NewDecoder(response.Body).Decode(&payload)
	if decodeErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur decodage des données - %s", decodeErr.Error())
	}

	details := &models.PokemonDetails{
		ID:        payload.ID,
		Name:      payload.Name,
		Image:     payload.Sprites.Other.OfficialArtwork.FrontDefault,
		Types:     []string{},
		Stats:     map[string]int{},
		Abilities: []string{},
	}

	for _, t := range payload.Types {
		details.Types = append(details.Types, t.Type.Name)
	}
	for _, s := range payload.Stats {
		details.Stats[s.Stat.Name] = s.BaseStat
	}
	for _, a := range payload.Abilities {
		details.Abilities = append(details.Abilities, a.Ability.Name)
	}

	return details, http.StatusOK, nil
}

func MoveService(offset int) (*models.MoveResponse, int, error) {
	client := http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/move?offset=%d", offset)
	request, requestErr := http.NewRequest(http.MethodGet, url, nil)
	if requestErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur initialisation requete - %s", requestErr.Error())
	}

	response, responseErr := client.Do(request)
	if responseErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur envoi requete - %s", responseErr.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, response.StatusCode, fmt.Errorf("erreur reponse requete - code : %d, status : %s", response.StatusCode, response.Status)
	}

	var listMoves models.MoveResponse
	decodeErr := json.NewDecoder(response.Body).Decode(&listMoves)
	if decodeErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur decodage des données - %s", decodeErr.Error())
	}

	return &listMoves, response.StatusCode, nil
}

func MoveByNameOrID(name string) (*models.MoveResponse, int, error) {
	client := http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/move/%s", name)
	request, requestErr := http.NewRequest(http.MethodGet, url, nil)
	if requestErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur initialisation requete - %s", requestErr.Error())
	}

	response, responseErr := client.Do(request)
	if responseErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur envoi requete - %s", responseErr.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, response.StatusCode, fmt.Errorf("erreur reponse requete - code : %d, status : %s", response.StatusCode, response.Status)
	}

	var payload struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	decodeErr := json.NewDecoder(response.Body).Decode(&payload)
	if decodeErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur decodage des données - %s", decodeErr.Error())
	}

	m := models.Move{
		Name: payload.Name,
		Url:  fmt.Sprintf("https://pokeapi.co/api/v2/move/%d/", payload.ID),
	}
	res := &models.MoveResponse{Results: []models.Move{m}}
	return res, http.StatusOK, nil
}

func GetMoveDetails(idOrName string) (*models.MoveDetails, int, error) {
	client := http.Client{Timeout: 5 * time.Second}
	url := fmt.Sprintf("https://pokeapi.co/api/v2/move/%s", idOrName)
	request, requestErr := http.NewRequest(http.MethodGet, url, nil)
	if requestErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur initialisation requete - %s", requestErr.Error())
	}

	response, responseErr := client.Do(request)
	if responseErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur envoi requete - %s", responseErr.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, response.StatusCode, fmt.Errorf("erreur reponse requete - code : %d, status : %s", response.StatusCode, response.Status)
	}

	var payload struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Accuracy *int   `json:"accuracy"`
		Power    *int   `json:"power"`
		PP       int    `json:"pp"`
		Type     struct {
			Name string `json:"name"`
		} `json:"type"`
		DamageClass struct {
			Name string `json:"name"`
		} `json:"damage_class"`
		EffectEntries []struct {
			Effect      string `json:"effect"`
			ShortEffect string `json:"short_effect"`
			Language    struct {
				Name string `json:"name"`
			} `json:"language"`
		} `json:"effect_entries"`
	}

	decodeErr := json.NewDecoder(response.Body).Decode(&payload)
	if decodeErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur decodage des données - %s", decodeErr.Error())
	}

	details := &models.MoveDetails{
		ID:          payload.ID,
		Name:        payload.Name,
		Accuracy:    payload.Accuracy,
		Power:       payload.Power,
		PP:          payload.PP,
		Type:        payload.Type.Name,
		DamageClass: payload.DamageClass.Name,
		Effect:      "",
		ShortEffect: "",
	}

	for _, e := range payload.EffectEntries {
		if e.Language.Name == "en" {
			details.Effect = e.Effect
			details.ShortEffect = e.ShortEffect
			break
		}
	}

	return details, http.StatusOK, nil
}
func TeamsService(idOrName string) (*models.PokemonDetails, int, error) {
	// Reuse existing GetPokemonDetails to retrieve a single pokemon's details
	return GetPokemonDetails(idOrName)
}

// SearchPokemon retrieves the full list of pokemons and filters by name substring (case-insensitive)
func SearchPokemon(query string) (*models.PokemonResponse, int, error) {
	client := http.Client{Timeout: 5 * time.Second}
	// Retrieve all pokemons (limit large enough to include all)
	url := "https://pokeapi.co/api/v2/pokemon?limit=100000&offset=0"
	request, requestErr := http.NewRequest(http.MethodGet, url, nil)
	if requestErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur initialisation requete - %s", requestErr.Error())
	}

	response, responseErr := client.Do(request)
	if responseErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur envoi requete - %s", responseErr.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, response.StatusCode, fmt.Errorf("erreur reponse requete - code : %d, status : %s", response.StatusCode, response.Status)
	}

	var payload struct {
		Results []struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"results"`
	}

	decodeErr := json.NewDecoder(response.Body).Decode(&payload)
	if decodeErr != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("erreur decodage des données - %s", decodeErr.Error())
	}

	matches := make([]models.Pokemon, 0)
	q := strings.ToLower(strings.TrimSpace(query))
	for _, r := range payload.Results {
		if strings.Contains(strings.ToLower(r.Name), q) {
			matches = append(matches, models.Pokemon{Name: r.Name, Url: r.Url})
		}
	}

	res := &models.PokemonResponse{Results: matches}
	return res, http.StatusOK, nil
}
