package controllers

import (
	"Swego/src/cmd"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

// UploadFile will receive and create the uploaded file
func UploadFile(w http.ResponseWriter, r *http.Request) {
	filepath := path.Join((cmd.RootFolder), path.Clean(r.URL.Path))
	if strings.Contains(r.URL.Path, "/private/") {
		r.URL.Path = strings.Replace(r.URL.Path, "/private/", "", 1)
		filepath = path.Join((cmd.PrivateFolder), path.Clean(r.URL.Path))
	}
	fmt.Println(r.URL.Path)
	// Maximum upload of 1000 MB files
	r.ParseMultipartForm(10000000 << 20)
	// If the variable files exists
	if r.MultipartForm != nil {

		// Get handler for filename, size and headers
		filesHandler := r.MultipartForm.File["files"]
		for _, handler := range filesHandler {
			file, err := handler.Open()
			if err != nil {
				fmt.Println("Error Retrieving the File")
				fmt.Println(err)
				return
			}

			defer file.Close()

			// Create file
			dst, err := os.Create(filepath + "/" + handler.Filename)
			defer dst.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Copy the uploaded file to the created file on the filesystem
			if _, err := io.Copy(dst, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			f, err := os.Open(filepath)
			defer f.Close()
			if err != nil {
				http.Error(w, "404 Not Found : Error while opening the file.", 404)
				return
			}
		}
		data := struct {
			Directory string
			ServerUA  string
		}{
			r.URL.Path,
			serverUA,
		}
		err := renderTemplate(w, "upload.tpl", data)
		if err != nil {
			fmt.Print("Error while uploading: ")
			fmt.Println(err)
		}
	}

}
