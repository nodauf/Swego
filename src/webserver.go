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
	"Swego/src/controllers/webdav"
	"Swego/src/routers"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	cmd.Execute()
	if cmd.Webdav {
		go webdav.New()
	}
	if cmd.Web {
		if cmd.Oneliners {
			go controllers.CliOnelinersMenu()
		}

		bind := strconv.Itoa(cmd.Bind)
		http.Handle("/", routers.Use(routers.Router))
		fmt.Printf("Sharing %s on %s:%s ...\n", cmd.RootFolder, cmd.IP, bind)
		if cmd.PrivateFolder != "" {
			http.Handle("/private/", routers.Use(routers.Router, controllers.BasicAuth))
			fmt.Printf("Sharing private %s on %s:%s ...\n", cmd.PrivateFolder, cmd.IP, bind)
		}
		var err error
		// Check if HTTPS or not
		if cmd.TLS {
			err = http.ListenAndServeTLS(cmd.IP+":"+bind, cmd.TLSCertificate, cmd.TLSKey, nil)
		} else {
			err = http.ListenAndServe(cmd.IP+":"+bind, nil)
		}
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	} else if cmd.Run {
		if cmd.Binary != "" {
			controllers.RunEmbeddedBinary(cmd.Binary, cmd.Args)
		} else {
			fmt.Println(controllers.EmbeddedFiles())
		}
	}
}
