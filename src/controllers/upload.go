package controllers

import (
	"net/http"
	"fmt"
	"io"
	"path"
	"os"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
        fmt.Println(r.URL.Path)
        // Maximum upload of 10 MB files
        r.ParseMultipartForm(10 << 20)

        // Get handler for filename, size and headers
        file, handler, err := r.FormFile("file")
        if err != nil {
                fmt.Println("Error Retrieving the File")
                fmt.Println(err)
                return
        }

        defer file.Close()
        fmt.Printf("Uploaded File: %+v\n", handler.Filename)
        fmt.Printf("File Size: %+v\n", handler.Size)
    //fmt.Printf("MIME Header: %+v\n", handler.Header)

        // Create file
        dst, err := os.Create(*Root_folder+r.URL.Path+handler.Filename)
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
        filepath := path.Join((*Root_folder), path.Clean(r.URL.Path))
        f, err := os.Open(filepath)
	defer f.Close()
        if err != nil {
                http.Error(w, "404 Not Found : Error while opening the file.", 404)
                return
        }
		data := struct {
			Directory string
		}{
			r.URL.Path,
		}
		err = renderTemplate(w , "views/upload.tpl", data)
        if err != nil {
                fmt.Print("Error while uploading: ")
				fmt.Println(err)
        }

}
