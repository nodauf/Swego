package routers

import (
    "net/http"
    "net/url"
    "Swego/src/controllers"
)

func ParseHttpParameter(w http.ResponseWriter, req *http.Request){
        query, errParseQuery := url.ParseQuery(req.URL.RawQuery)

        if errParseQuery == nil {
            if len(query["embedded"]) > 0{ // Manage embedded files
                controllers.EmbeddedRequest(w, req)
                return
            }else if len(query["newFolder"]) > 0{
                controllers.CreateFolder(w, req)
               return
            }
        }
        controllers.HandleFile(w, req)
}
