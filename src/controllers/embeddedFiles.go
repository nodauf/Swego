package controllers

import (
	"compress/gzip"
	"compress/zlib"
	"container/list"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"

	"SimpleHTTPServer-golang/src/utils"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gabriel-vasile/mimetype"
	"github.com/yeka/zip"
)

const pathEmbedded = "./assets/embedded/"

func EmbeddedRequest(w http.ResponseWriter, req *http.Request) {
	requestPath := strings.Split(req.RequestURI, "?")[0]

	if requestPath[len(requestPath)-1:] == "/" { // Directory listing if it's a folder (last character is /)
		handleEmbeddedDirectory(requestPath, w, req)
	} else { // It's a file, we serve the file
		fileName := requestPath[1:]
		serveEmbeddedFile(fileName, w, req)
	}
}

func serveEmbeddedFile(filePath string, w http.ResponseWriter, req *http.Request) {
	//Can't use variable, otherwise rice generate an error when rice embed-go ....
	templateBox, err := rice.FindBox("../assets/embedded/")

	if err != nil {
		log.Fatal(err)
	}
	// Opening the file handle
	//f, err := os.Open(filePath)
	f, err := templateBox.Open(filePath)
	// Content-Type handling
	query, errParseQuery := url.ParseQuery(req.URL.RawQuery)

	if err != nil {
		http.Error(w, "404 Not Found : Error while opening the file.", 404)
		log.Println("404 Not Found : Error while opening the file " + filePath)
		return
	}
	defer f.Close()

	// Checking if the opened handle is really a file
	statinfo, err := f.Stat()

	//buf := make([]byte, utils.Min(fs_maxbufsize, statinfo.Size()))
	//buf := make([]byte, statinfo.Size())
	//fmt.Println(len(buf))
	//f.Read(buf)

	if err != nil || errParseQuery != nil {
		http.Error(w, "500 Internal Error : stat() failure.", 500)
		log.Println("500 Internal Error : stat() failure for the file: " + filePath)
		return
	}
	if errParseQuery == nil && len(query["dl"]) > 0 { // The user explicitedly wanted to download the file (Dropbox style!)
		w.Header().Set("Content-Type", "application/octet-stream")
	} else if errParseQuery == nil && len(query["dlenc"]) > 0 { // Download the file as an encrypted zip

		// Absolute path to the file
		//filePathName := f.Name()
		// Create the zip file can't create in embedded path. Create temporarily in root of the webserver
		zipFile, err := os.Create(filePath + ".zip")
		if err != nil {
			log.Fatalln(err)
		}
		zipFilePath := zipFile.Name()
		zipw := zip.NewWriter(zipFile)

		// Add file f to the zip
		utils.AddRicefiletoZip(statinfo.Name(), f, filePath, zipw, true, "infected")

		// Manually close the zip
		zipw.Close()

		// Generate the request for the new file
		newFile := strings.Split(req.URL.String(), "?")
		fmt.Println(zipFilePath)
		newRequest, _ := http.NewRequest("GET", "http://"+req.Host+newFile[0], nil)

		// Serve the new file (encrypted zip)
		serveFile(zipFilePath, w, newRequest)
		os.Remove(zipFilePath)
		return
	} else {
		// Need its own rice.file otherwise it will miss the first chunck
		fileForMime, _ := templateBox.Open(filePath)
		defer fileForMime.Close()
		// Fetching file's mimetype and giving it to the browser
		if mimetype, _ := mimetype.DetectReader(fileForMime); mimetype.String() != "" {
			w.Header().Set("Content-Type", mimetype.String())
		} else {
			w.Header().Set("Content-Type", "application/octet-stream")
		}
	}

	// Manage gzip/zlib compression
	output_writer := w.(io.Writer)

	is_compressed_reply := false

	if (*Uses_gzip) == true && req.Header.Get("Accept-Encoding") != "" {
		encodings := utils.ParseCSV(req.Header.Get("Accept-Encoding"))

		for _, val := range encodings {
			if val == "gzip" {
				w.Header().Set("Content-Encoding", "gzip")
				output_writer = gzip.NewWriter(w)

				is_compressed_reply = true

				break
			} else if val == "deflate" {
				w.Header().Set("Content-Encoding", "deflate")
				output_writer = zlib.NewWriter(w)

				is_compressed_reply = true

				break
			}
		}
	}

	if !is_compressed_reply {
		// Add Content-Length
		w.Header().Set("Content-Length", strconv.FormatInt(statinfo.Size(), 10))
	}

	// Stream data out !
	buf := make([]byte, utils.Min(fs_maxbufsize, statinfo.Size()))
	n := 0

	for err == nil {
		n, err = f.Read(buf)
		buf = utils.SearchAndReplace(SearchAndReplaceMap, buf)
		output_writer.Write(buf[0:n])
	}
	// Closes current compressors
	switch output_writer.(type) {
	case *gzip.Writer:
		output_writer.(*gzip.Writer).Close()
	case *zlib.Writer:
		output_writer.(*zlib.Writer).Close()
	}
	//f.Close()
}

func listEmbeddedFiles() ([]string, []string) {
	//Can't use variable, otherwise rice generate an error when rice embed-go ....
	templateBox, err := rice.FindBox("../assets/embedded/")
	if err != nil {
		log.Fatal(err)
	}
	// Otherwise, generate folder content.
	children_dir_tmp := list.New()
	children_files_tmp := list.New()
	err = templateBox.Walk("/", func(path string, info os.FileInfo, err error) error {
		//fmt.Println(path)
		if info.IsDir() {
			children_dir_tmp.PushBack(info.Name())
		} else {
			children_files_tmp.PushBack(info.Name())
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// And transfer the content to the final array structure
	children_dir := utils.CopyToArray(children_dir_tmp)
	children_files := utils.CopyToArray(children_files_tmp)

	return children_dir, children_files
}

func handleEmbeddedDirectory(path string, w http.ResponseWriter, req *http.Request) {
	if !*DisableDirectoryListing {
		children_dir, children_files := listEmbeddedFiles()

		//Sort children_dir and children_files
		sort.Slice(children_dir, func(i, j int) bool { return children_dir[i] < children_dir[j] })

		//Sort children_dir and children_files
		sort.Slice(children_files, func(i, j int) bool { return children_files[i] < children_files[j] })

		data := utils.Dirlisting{Name: req.URL.Path,
			ServerUA:       serverUA,
			Children_dir:   children_dir,
			Children_files: children_files,
			Embedded:       true}
		err := renderTemplate(w, "directoryListing.tpl", data)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func readEmbeddedBinary(binary string) []byte {
	//Can't use variable, otherwise rice generate an error when rice embed-go ....
	templateBox, err := rice.FindBox("../assets/embedded/")
	if err != nil {
		log.Fatal(err)
	}
	binBytes, err := templateBox.Bytes(binary)
	if err != nil {
		fmt.Printf("[!] Error finding binary: %s\n", binary)
		log.Fatal(err)
	}
	return binBytes
}
