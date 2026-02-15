package helpers

import (
	"net/http"
	"net/url"
	"strconv"
)

func RedirectToError(w http.ResponseWriter, r *http.Request, code int, message string) {
	params := url.Values{}
	if code > 0 {
		params.Set("code", strconv.Itoa(code))
	}
	if message != "" {
		params.Set("message", message)
	}

	pathTarget := "/error"
	if encodeParams := params.Encode(); encodeParams != "" {
		pathTarget += "?" + encodeParams
	}
	http.Redirect(w, r, pathTarget, http.StatusSeeOther)
}
