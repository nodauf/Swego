package controllers

import (
	"Swego/src/cmd"
	"Swego/src/utils"
	"Swego/src/views"
	"compress/gzip"
	"compress/zlib"
	"container/list"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/yeka/zip"
)

const serverUA = ""
const fsMaxbufsize = 4096 // 4096 bits = default page size on OSX

// HandleFile is the entrypoint to manage the web request
func HandleFile(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Server", serverUA)

	path := filepath.Join((cmd.RootFolder), filepath.Clean(req.URL.Path))

	if strings.Contains(req.URL.Path, "/private/") {
		req.URL.Path = strings.Replace(req.URL.Path, "/private/", "", 1)
		path = filepath.Join((cmd.PrivateFolder), filepath.Clean(req.URL.Path))
	}
	serveFile(path, w, req)

}

func serveFile(filePath string, w http.ResponseWriter, req *http.Request) {
	// Opening the file handle
	f, err := os.Open(filePath)

	// Content-Type handling
	query, errParseQuery := url.ParseQuery(req.URL.RawQuery)

	//        if errParseQuery == nil && len(query["embedded"]) > 0{ // Manage embedded files
	//            embeddedRequest(w, req)
	//            return
	//        }

	if err != nil {
		http.Error(w, "404 Not Found : Error while opening the file.", 404)
		log.Println("404 Not Found : Error while opening the file " + filePath)
		return
	}

	defer f.Close()

	// Checking if the opened handle is really a file
	statinfo, err := f.Stat()
	if err != nil || errParseQuery != nil {
		http.Error(w, "500 Internal Error : stat() failure.", 500)
		log.Println("500 Internal Error : stat() failure for the file: " + filePath)
		return
	}

	if statinfo.IsDir() { // If it's a directory, open it !
		if errParseQuery == nil && len(query["dl"]) > 0 {
			zipFilePath := utils.ZipDirectory(f, false)

			// Generate the request for the new file - remove ?dl to download the file
			newFile := strings.Split(req.URL.String(), "?")
			newRequest, _ := http.NewRequest("GET", "http://"+req.Host+newFile[0], nil)

			// Serve the new file (encrypted zip)
			serveFile(zipFilePath, w, newRequest)

			// Remove the zip file
			os.Remove(zipFilePath)

			return

		} else if errParseQuery == nil && len(query["dlenc"]) > 0 {
			zipFilePath := utils.ZipDirectory(f, true)
			// Generate the request for the new file - remove ?dl to download the file

			newFile := strings.Split(req.URL.String(), "?")
			newRequest, _ := http.NewRequest("GET", "http://"+req.Host+newFile[0], nil)

			// Serve the new file (encrypted zip)
			serveFile(zipFilePath, w, newRequest)

			// Remove the zip file
			os.Remove(zipFilePath)
			return
		} else {
			handleDirectory(f, w, req)
		}
		return
	}

	if (statinfo.Mode() &^ 07777) == os.ModeSocket { // If it's a socket, forbid it !
		http.Error(w, "403 Forbidden : you can't access this resource.", 403)
		return
	}

	// Manages If-Modified-Since and add Last-Modified (taken from Golang code)
	if t, err := time.Parse(http.TimeFormat, req.Header.Get("If-Modified-Since")); err == nil && statinfo.ModTime().Unix() <= t.Unix() {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	w.Header().Set("Last-Modified", statinfo.ModTime().Format(http.TimeFormat))

	if errParseQuery == nil && len(query["dl"]) > 0 { // The user explicitedly wanted to download the file (Dropbox style!)
		w.Header().Set("Content-Type", "application/octet-stream")
	} else if errParseQuery == nil && len(query["dlenc"]) > 0 { // Download the file as an encrypted zip

		// Absolute path to the file
		filePathName := f.Name()
		// Create the zip file
		zipFile, err := os.Create(filePathName + ".zip")
		if err != nil {
			log.Fatalln(err)
		}
		zipFilePath := zipFile.Name()
		zipw := zip.NewWriter(zipFile)

		// Add file f to the zip
		utils.AddfiletoZip(statinfo.Name(), f, zipw, true, "infected")

		// Manually close the zip
		zipw.Close()

		// Generate the request for the new file
		newFile := strings.Split(req.URL.String(), "?")
		newRequest, _ := http.NewRequest("GET", "http://"+req.Host+newFile[0], nil)

		// Serve the new file (encrypted zip)
		serveFile(zipFilePath, w, newRequest)
		os.Remove(zipFilePath)
		return
	} else {
		// Need its own rice.file otherwise it will miss the first chunck
		fileForMime, _ := os.Open(filePath)
		defer fileForMime.Close()
		// Fetching file's mimetype and giving it to the browser
		if mimetype, _ := mimetype.DetectReader(fileForMime); mimetype.String() != "" {
			w.Header().Set("Content-Type", mimetype.String())
		} else {
			w.Header().Set("Content-Type", "application/octet-stream")
		}
	}
	w.Header().Set("Cache-Control", "store, public, min-age=5, max-age=120")
	// Manage Content-Range (TODO: Manage end byte and multiple Content-Range)
	if req.Header.Get("Range") != "" {
		startByte := utils.ParseRange(req.Header.Get("Range"))

		if startByte < statinfo.Size() {
			f.Seek(startByte, 0)
		} else {
			startByte = 0
		}

		w.Header().Set("Content-Range",
			fmt.Sprintf("bytes %d-%d/%d", startByte, statinfo.Size()-1, statinfo.Size()))
	}

	// Manage gzip/zlib compression
	outputWriter := w.(io.Writer)

	isCompressedReply := false

	if (cmd.Gzip) == true && req.Header.Get("Accept-Encoding") != "" {
		encodings := utils.ParseCSV(req.Header.Get("Accept-Encoding"))

		for _, val := range encodings {
			if val == "gzip" {
				w.Header().Set("Content-Encoding", "gzip")
				outputWriter = gzip.NewWriter(w)

				isCompressedReply = true

				break
			} else if val == "deflate" {
				w.Header().Set("Content-Encoding", "deflate")
				outputWriter = zlib.NewWriter(w)

				isCompressedReply = true

				break
			}
		}
	}

	if !isCompressedReply {
		// Add Content-Length
		w.Header().Set("Content-Length", strconv.FormatInt(statinfo.Size(), 10))
	}

	// Stream data out !
	buf := make([]byte, utils.Min(fsMaxbufsize, statinfo.Size()))
	n := 0
	for err == nil {
		n, err = f.Read(buf)
		buf = utils.SearchAndReplace(cmd.SearchAndReplaceMap, buf)
		outputWriter.Write(buf[0:n])
	}

	// Closes current compressors
	switch outputWriter.(type) {
	case *gzip.Writer:
		outputWriter.(*gzip.Writer).Close()
	case *zlib.Writer:
		outputWriter.(*zlib.Writer).Close()
	}

	//f.Close()
}

func handleDirectory(f *os.File, w http.ResponseWriter, req *http.Request) {
	if !cmd.DisableListing {
		names, _ := f.Readdir(-1)

		// First, check if there is any index in this folder.
		for _, val := range names {
			if val.Name() == "index.html" {
				serveFile(filepath.Join(f.Name(), "index.html"), w, req)
				return
			}
		}

		// Otherwise, generate folder content.
		childrenDirTmp := list.New()
		childrenFilesTmp := list.New()

		for _, val := range names {
			//if val.Name()[0] == '.' {
			//        continue
			//} // Remove hidden files from listing

			if val.IsDir() {
				childrenDirTmp.PushBack(val.Name())
			} else {
				childrenFilesTmp.PushBack(val.Name())
			}
		}

		// And transfer the content to the final array structure
		childrenDir := utils.CopyToArray(childrenDirTmp)
		childrenFiles := utils.CopyToArray(childrenFilesTmp)
		//Sort children_dir and children_files
		sort.Slice(childrenDir, func(i, j int) bool { return childrenDir[i] < childrenDir[j] })

		//Sort children_dir and children_files
		sort.Slice(childrenFiles, func(i, j int) bool { return childrenFiles[i] < childrenFiles[j] })

		data := utils.Dirlisting{Name: req.URL.Path, ServerUA: serverUA,
			ChildrenDir: childrenDir, ChildrenFiles: childrenFiles}
		err := renderTemplate(w, "directoryListing.tpl", data)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func renderTemplate(w http.ResponseWriter, view string, data interface{}) error {
	//templateBox, err := rice.FindBox("../views/")
	//if err != nil {
	//	log.Fatal(err)
	//}
	// get file contents as string
	//templateString, err := templateBox.String(view)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//tpl := template.Must(template.Parse(templateString))
	content, err := views.GetViews(view)
	if err != nil {
		log.Fatal(err)
	}
	tpl, err := template.New("tpl").Parse(string(content))
	if err != nil {
		http.Error(w, "500 Internal Error : Error while generating directory listing.", 500)
		fmt.Println(err)
		return err
	}

	err = tpl.Execute(w, data)
	if err != nil {
		return err
	}
	return nil

}

func DeleteRequest(w http.ResponseWriter, req *http.Request) {
	path := filepath.Join((cmd.RootFolder), filepath.Clean(req.URL.Path))
	if strings.Contains(req.URL.Path, "/private/") {
		req.URL.Path = strings.Replace(req.URL.Path, "/private/", "", 1)
		path = filepath.Join((cmd.PrivateFolder), filepath.Clean(req.URL.Path))
	}

	f, err := os.Open(path)
	if err != nil {
		http.Error(w, "404 Not Found : Can't delete the file.", 404)
		log.Println("404 Not Found : Can't delete the file " + path)
		return
	}
	statinfo, err := f.Stat()
	if err != nil {
		http.Error(w, "500 Internal Error : stat() failure.", 500)
		log.Println("500 Internal Error : stat() failure for the file: " + path)
		return
	}

	if statinfo.IsDir() { // If it's a directory, open it !
		err := os.RemoveAll(path) // delete an entire directory
		if err != nil {
			http.Error(w, "500 Internal Error : Can't delete the folder", 500)
			log.Println("500 Internal Error : Can't delete the folder:  " + path)
		}
		http.Redirect(w, req, filepath.Join(req.URL.Path+"/../"), 302)
	} else {
		err := os.Remove(path) // remove a single file
		if err != nil {
			http.Error(w, "500 Internal Error : Can't delete the file", 500)
			log.Println("500 Internal Error : Can't delete the file:  " + path)
		}

		http.Redirect(w, req, filepath.Join(req.URL.Path+"/../"), 302)
	}
}
