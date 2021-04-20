package webdav

import (
	"Swego/src/cmd"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/webdav"
)

func New() {
	srv := &CustomHandler{
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
		bind := strconv.Itoa(cmd.WebdavPort)
		srv := &http.Server{
			Handler:      srv,
			Addr:         cmd.IP + ":" + bind,
			WriteTimeout: 0,
			ReadTimeout:  0,
			IdleTimeout:  5 * time.Second,
			TLSConfig:    &cmd.TLSConfig,
		}
		listenTLS, _ := tls.Listen("tcp", cmd.IP+":"+bind, &cmd.TLSConfig)
		_ = srv.Serve(listenTLS)
		//http.ListenAndServeTLS(fmt.Sprintf(":%d", cmd.WebdavPort), cmd.TLSCertificate, cmd.TLSKey, srv)
	}

	fmt.Printf("Sharing with webdav %s on %s:%d ...\n", cmd.RootFolder, cmd.IP, cmd.WebdavPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cmd.WebdavPort), srv); err != nil {
		log.Fatalf("Error with WebDAV server: %v", err)
	}
}
