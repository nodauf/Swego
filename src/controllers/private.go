package controllers

import (
	"encoding/base64"
	"net/http"
	"strings"
)

func BasicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Not authorized", 401)
			return
		}

		if pair[0] != *Username || pair[1] != *Password {
			http.Error(w, "Not authorized", 401)
			return
		}
		//fmt.Println("Serving: "+ path.Join((*Private), path.Clean(r.URL.Path)))
		//		h.ServeHTTP(w, r)
		HandleFile(w, r)
	}
}
