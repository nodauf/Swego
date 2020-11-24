/* Tiny web server in Golang for sharing a folder
Copyright (c) 2010-2014 Alexis ROBERT <alexis.robert@gmail.com>
Contains some code from Golang's http.ServeFile method, and
uses lighttpd's directory listing HTML template. */

package main

import "net/http"
import "fmt"
//import "SimpleHTTPServer-golang/src/utils"
import "SimpleHTTPServer-golang/src/routers"
import "SimpleHTTPServer-golang/src/controllers"


func main() {
        controllers.ParseArgs()
        http.Handle("/"+*controllers.Public, routers.Use(routers.Router))

        fmt.Printf("Sharing %s/%s on %s ...\n", *controllers.Root_folder, *controllers.Public, *controllers.Bind)
        if *controllers.Private != "" {
            http.Handle("/"+*controllers.Private+"/", routers.Use(routers.Router,controllers.BasicAuth))
            fmt.Printf("Sharing %s/%s on %s ...\n", *controllers.Root_folder, *controllers.Private, *controllers.Bind)
        }
        http.ListenAndServe(":"+(*controllers.Bind), nil)
}



