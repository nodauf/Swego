package routers

import "fmt"
import "net/http"
import "SimpleHTTPServer-golang/src/controllers"

// Routing function
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
         ParseHttpParameter(w, req)

    case "POST":
        // Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
        if err := req.ParseForm(); err != nil {
            fmt.Fprintf(w, "ParseForm() err: %v", err)
            return
        }
                controllers.UploadFile(w,req)

    default:
        fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
    }
}

// See https://gist.github.com/elithrar/7600878#comment-955958 for how to extend it to suit simple http.Handler's
func Use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}
