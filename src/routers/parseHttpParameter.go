package routers

import (
	"Swego/src/controllers"
	"net/http"
	"net/url"
)

// ParseHTTPParameter will parse the http parameter to execute specific action if one is present
func ParseHTTPParameter(w http.ResponseWriter, req *http.Request) {
	query, errParseQuery := url.ParseQuery(req.URL.RawQuery)

	if errParseQuery == nil {
		if len(query["embedded"]) > 0 { // Manage embedded files
			controllers.EmbeddedRequest(w, req)
			return
		} else if len(query["newFolder"]) > 0 {
			controllers.CreateFolder(w, req)
			return
		}
	}
	controllers.HandleFile(w, req)
}
