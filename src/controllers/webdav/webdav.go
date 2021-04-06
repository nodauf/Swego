package webdav

import (
	"Swego/src/cmd"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/webdav"
)

func New() {
	srv := &webdav.Handler{
		FileSystem: webdav.Dir(cmd.RootFolder),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				log.Printf("WEBDAV [%s]: %s, ERROR: %s\n", r.Method, r.URL, err)
			} else {
				log.Printf("WEBDAV [%s]: %s \n", r.Method, r.URL)
			}
		},
	}
	if cmd.TLS == true {

		http.ListenAndServeTLS(fmt.Sprintf(":%d", cmd.WebdavPort), cmd.TLSCertificate, cmd.TLSKey, srv)
	}
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cmd.WebdavPort), srv); err != nil {
		log.Fatalf("Error with WebDAV server: %v", err)
	}

}
