/* Tiny web server in Golang for sharing a folder
Copyright (c) 2010-2014 Alexis ROBERT <alexis.robert@gmail.com>
Contains some code from Golang's http.ServeFile method, and
uses lighttpd's directory listing HTML template. */

package main

import (
	"fmt"
	"log"
	"net/http"
    "strconv"

	"SimpleHTTPServer-golang/src/controllers"
	"SimpleHTTPServer-golang/src/routers"
)

func main() {
	controllers.ParseArgs()
	if controllers.WebCommand.Parsed() {
        bind := strconv.Itoa(*controllers.Bind)
		http.Handle("/", routers.Use(routers.Router))
        fmt.Printf("Sharing %s on %s:%s ...\n", *controllers.Root_folder, *controllers.IP, bind)
		if *controllers.Private != "" {
			http.Handle("/private/", routers.Use(routers.Router, controllers.BasicAuth))
            fmt.Printf("Sharing private %s on %s:%s ...\n", *controllers.Private, *controllers.IP, bind)
		}
		var err error
		// Check if HTTPS or not
		if *controllers.Tls {
			err = http.ListenAndServeTLS(*controllers.IP+":"+bind, *controllers.Certificate, *controllers.Key, nil)
		} else {
            err = http.ListenAndServe(*controllers.IP+":"+bind, nil)
		}
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	} else if controllers.RunCommand.Parsed() {
		if *controllers.Binary != "" {
			controllers.RunEmbeddedBinary(*controllers.Binary, *controllers.Args)
		} else {
			controllers.PrintEmbeddedFiles()
		}
	}
}
