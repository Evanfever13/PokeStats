package helpers

import (
	"net/http"
	"net/url"
	"strconv"
)

/*=========================================================
  HELPER : Redirection a la page error.html
  =========================================================*/

// (Source : Moodle)
// Permet la redirection à la page error.html
func RedirectToError(w http.ResponseWriter, r *http.Request, code int, message string) {
	//Recupere le Code et le Message
	params := url.Values{}
	if code > 0 {
		params.Set("code", strconv.Itoa(code))
	}
	if message != "" {
		params.Set("message", message)
	}

	//Redirige vers error.html
	pathTarget := "/error"
	if encodeParams := params.Encode(); encodeParams != "" {
		pathTarget += "?" + encodeParams
	}
	http.Redirect(w, r, pathTarget, http.StatusSeeOther)
}
