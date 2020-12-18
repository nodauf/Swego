/* Tiny web server in Golang for sharing a folder
Copyright (c) 2010-2014 Alexis ROBERT <alexis.robert@gmail.com>
Contains some code from Golang's http.ServeFile method, and
uses lighttpd's directory listing HTML template. */

package main

import (
    "net/http"
    "fmt"
    "log"

    "SimpleHTTPServer-golang/src/routers"
    "SimpleHTTPServer-golang/src/controllers"
)

func main() {
        controllers.ParseArgs()
        if controllers.WebCommand.Parsed(){
            http.Handle("/", routers.Use(routers.Router))

            fmt.Printf("Sharing %s on %s ...\n", *controllers.Root_folder, *controllers.Bind)
            if *controllers.Private != "" {
                http.Handle("/private/", routers.Use(routers.Router,controllers.BasicAuth))
                fmt.Printf("Sharing private %s on %s ...\n",  *controllers.Private, *controllers.Bind)
            }
            var err error
            // Check if HTTPS or not
            if *controllers.Tls {
                err = http.ListenAndServeTLS(":"+(*controllers.Bind), *controllers.Certificate, *controllers.Key,nil)
            }else{
                err = http.ListenAndServe(":"+(*controllers.Bind), nil)
            }
            if err != nil {
                log.Fatal("ListenAndServe: ", err)
            }
        } else if controllers.Run{
            if controllers.Binary != ""{
                controllers.RunEmbeddedBinary(controllers.Binary, controllers.Args)
            }else{
                controllers.PrintEmbeddedFiles()
            }
        }
}



