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

var listTemp *template.Template
var ListTeams []string

func Load() {
	funcMap := template.FuncMap{
		"imageURL": func(u string) string {
			id := path.Base(strings.TrimSuffix(u, "/"))
			return fmt.Sprintf("https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/other/official-artwork/%s.png", id)
		},
		"pokemonID": func(u string) string {
			return path.Base(strings.TrimSuffix(u, "/"))
		},
	}

	listTemplates, errTemplates := template.New("").Funcs(funcMap).ParseGlob("./../../templates/*.html")
	if errTemplates != nil {
		log.Fatalf("Erreur chargement des templates : %s", errTemplates.Error())
	}
	listTemp = listTemplates
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, name string, data interface{}) {

	var buffer bytes.Buffer

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
	_, _ = buffer.WriteTo(w)
}

func AddTeamID(id string) {
	if id == "" {
		return
	}
	for _, v := range ListTeams {
		if v == id {
			return
		}
	}
	ListTeams = append(ListTeams, id)
}

func RemoveTeamID(id string) {
	if id == "" {
		return
	}
	newList := make([]string, 0, len(ListTeams))
	for _, v := range ListTeams {
		if v != id {
			newList = append(newList, v)
		}
	}
	ListTeams = newList
}
