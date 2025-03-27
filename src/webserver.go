/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"Swego/src/cmd"
	"Swego/src/controllers"
	"Swego/src/routers"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	cmd.Execute()
	if cmd.Web {
		if cmd.Oneliners {
			go controllers.CliOnelinersMenu()
		}

		bind := strconv.Itoa(cmd.Bind)
		http.Handle("/", routers.Use(routers.Router, controllers.IPFiltering, controllers.BasicAuth))
		fmt.Printf("Sharing %s on %s:%s ...\n", cmd.RootFolder, cmd.IP, bind)
		if cmd.PrivateFolder != "" {
			http.Handle("/private/", routers.Use(routers.Router, controllers.IPFiltering, controllers.BasicAuth))
			fmt.Printf("Sharing private %s on %s:%s ...\n", cmd.PrivateFolder, cmd.IP, bind)
		}
		var err error
		// Check if HTTPS or not
		if cmd.TLS {

			srv := &http.Server{
				Addr:         cmd.IP + ":" + bind,
				WriteTimeout: 0,
				ReadTimeout:  0,
				IdleTimeout:  5 * time.Second,
				TLSConfig:    &cmd.TLSConfig,
			}
			listenTLS, err := tls.Listen("tcp", cmd.IP+":"+bind, &cmd.TLSConfig)
			if err != nil {
				log.Fatal("Error starting listening for the webserver: " + err.Error())
			}
			err = srv.Serve(listenTLS)
			if err != nil {
				log.Fatal("Error starting the webserver: " + err.Error())
			}
		} else {
			err = http.ListenAndServe(cmd.IP+":"+bind, nil)
		}
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	} //else if cmd.Run {
	//if cmd.Binary != "" {
	//controllers.RunEmbeddedBinary(cmd.Binary, cmd.Args)
	//} else {
	//fmt.Println(controllers.EmbeddedFiles())
	//	}
	//}
}
