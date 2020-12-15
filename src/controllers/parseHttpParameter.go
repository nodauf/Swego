package controllers

import (
    "net/http"
    "net/url"
)

func ParseHttpParameter(w http.ResponseWriter, req *http.Request){
        query, errParseQuery := url.ParseQuery(req.URL.RawQuery)

        if errParseQuery == nil {
            if len(query["embedded"]) > 0{ // Manage embedded files
                embeddedRequest(w, req)
                return
            }else if len(query["newFolder"]) > 0{
                createFolder(w, req)
               return
            }
        }
        handleFile(w, req)
}
