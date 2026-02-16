package templates

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// Decaration des Variables global
var listTemp *template.Template
var ListTeams []string

// (Source : Moodle)
// Charge les templates
func Load() {
	//(Source : Cyril)
	//Permet de charger les images des pokemon dans la page Pokemon
	funcMap := template.FuncMap{
		"imageURL": func(u string) string {
			id := path.Base(strings.TrimSuffix(u, "/"))
			return fmt.Sprintf("https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/%s.png", id)
		},
		"pokemonID": func(u string) string {
			return path.Base(strings.TrimSuffix(u, "/"))
		},
	}

	//Cherche les template dans "./../../templates/*.html"
	listTemplates, errTemplates := template.New("").Funcs(funcMap).ParseGlob("./../../templates/*.html")

	//Gestion des erreur
	if errTemplates != nil {
		log.Fatalf("Erreur chargement des templates : %s", errTemplates.Error())
	}

	listTemp = listTemplates
}

// (Source : Moodle)
// Afficher les templates
func RenderTemplate(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	//Crée un buffer temporaire
	var buffer bytes.Buffer

	//Cherche des erreur liée au chargement de la page
	errRender := listTemp.ExecuteTemplate(&buffer, name, data)
	if errRender != nil {
		http.Redirect(
			w,
			r,
			fmt.Sprintf(
				"/error?code=%d&message=%s",
				http.StatusInternalServerError,
				url.QueryEscape("Erreur lors du chargement de la page"),
			),
			http.StatusSeeOther,
		)
		return
	}
	//Affiche la page
	_, _ = buffer.WriteTo(w)
}

// Ajoute un pokemon au Favoris
func AddTeamID(id string) {
	//Test si l'id n'est pas vide
	if id == "" {
		return
	}

	//Cherche si l'id n'est pas deja dans la liste
	for _, v := range ListTeams {
		if v == id {
			return
		}
	}

	//Ajoute le Pokemon au favoris
	ListTeams = append(ListTeams, id)

	//Limite le maximum de Pokemon favoris à 6
	if len(ListTeams) == 7 {
		RemoveTeamID(id)
	}
}

// Supprime un pokemon au Favoris
func RemoveTeamID(id string) {
	//Test si l'id n'est pas vide
	if id == "" {
		return
	}

	//Crée une liste temporaire égale aux favoris sans le Pokemon à supprimer
	newList := make([]string, 0, len(ListTeams))
	for _, v := range ListTeams {
		if v != id {
			newList = append(newList, v)
		}
	}

	//Update les favoris
	ListTeams = newList
}
