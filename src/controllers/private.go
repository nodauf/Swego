package controllers

import (
	"Swego/src/cmd"
	"Swego/src/utils"
	"encoding/base64"
	"net/http"
	"path/filepath"
	"strings"
)

// BasicAuth function to handle private directory and print basic authentication
func BasicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join((cmd.RootFolder), filepath.Clean(r.URL.Path))
		isSubdirectory, err := utils.PathSubElem(cmd.PrivateFolder, path)
		utils.Check(err, "fail to check subdirectory")

		if isSubdirectory {
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

			if pair[0] != cmd.Username || pair[1] != cmd.Password {
				http.Error(w, "Not authorized", 401)
				return
			}
			//fmt.Println("Serving: "+ path.Join((*Private), path.Clean(r.URL.Path)))
			//		h.ServeHTTP(w, r)
			//HandleFile(w, r)
		}
		h(w, r)
	}
}
