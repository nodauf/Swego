package controllers

import (
	"Swego/src/cmd"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
)

// CreateFolder will create a folder from the web request
func CreateFolder(w http.ResponseWriter, req *http.Request) {
	dirPath := path.Join((cmd.RootFolder), path.Clean(req.URL.Path))
	query, errParseQuery := url.ParseQuery(req.URL.RawQuery)
	folder := query["newFolder"][0]
	// Check for directory traversal
	if strings.Contains(folder, "..") {
		http.Error(w, "500 Internal Error : Invalid character", 500)
		log.Println("Invalid character on " + folder)
		return
	}

	folderToCreate := dirPath + "/" + folder

	f, err := os.Open(dirPath)
	if err != nil {
		http.Error(w, "404 Not Found : Error while opening the directory.", 404)
		log.Println("404 Not Found : Error while opening the directory " + dirPath)
		return
	}

	statinfo, err := f.Stat()
	if errParseQuery != nil || err != nil {
		http.Error(w, "500 Internal Error : stat() failure.", 500)
		log.Println("Failed to create folder on: : " + dirPath)
		return
	}

	if !statinfo.IsDir() {
		http.Error(w, path.Clean(req.URL.Path)+" is not a folder", 500)
		log.Println("Failed to create folder on: : " + dirPath)
		return
	}

	// If there is no error it's means the directory (or file) exists
	_, err = os.Open(folderToCreate)
	if err == nil {
		http.Error(w, "500 Internal Error : Directory already exists.", 500)
		log.Println("500 Internal Error : Directory already exists." + folderToCreate)
		return
	}
	// All test passed - create the folder
	err = os.Mkdir(folderToCreate, 0755)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Folder " + folderToCreate + " created")
	http.Redirect(w, req, path.Clean(req.URL.Path), 302)
}
