package routers

import (
	"Swego/src/cmd"
	"Swego/src/controllers"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Router : Routing function
func Router(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("\"%s %s %s %s\" \"%s\" \"%s\"\n",
		req.RemoteAddr,
		req.Method,
		req.URL.String(),
		req.Proto,
		req.Referer(),
		req.UserAgent()) // TODO: Improve this crappy logging
	switch req.Method {
	case "GET":
		//controllers.HandleFile(w, req)
		//controllers.ParseHttpParameter(w, req)
		ParseHTTPParameter(w, req)

	case "POST":
		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		if cmd.Verbose {
			data, err := ioutil.ReadAll(req.Body)
			if err == nil && len(data) > 0 {
				fmt.Println(string(data))
			}
			req.ParseForm()
			for key, value := range req.Form {
				fmt.Printf("%s => %s \n", key, value)
			}

		}
		controllers.UploadFile(w, req)

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

// Use function to implement middleware on the router
// See https://gist.github.com/elithrar/7600878#comment-955958 for how to extend it to suit simple http.Handler's
func Use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)

	}
	return h
}
