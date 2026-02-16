package models

/*=========================================================
  MODELS : Structure de données
  =========================================================*/

//Pokemon
type PokemonResponse struct {
	Results []Pokemon `json:"results"`
}
type Pokemon struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

//Pokemon Details
type PokemonDetails struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	Image     string         `json:"image"`
	Types     []string       `json:"types"`
	Stats     map[string]int `json:"stats"`
	Abilities []string       `json:"abilities"`
}

//Move
type MoveResponse struct {
	Results []Move `json:"results"`
}
type Move struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

//Move Details
type MoveDetails struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Accuracy    *int   `json:"accuracy"`
	Power       *int   `json:"power"`
	PP          int    `json:"pp"`
	Type        string `json:"type"`
	DamageClass string `json:"damage_class"`
	Effect      string `json:"effect"`
	ShortEffect string `json:"short_effect"`
}

//Team
type TeamResponse struct {
	Results []Team `json:"results"`
}
type Team struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
