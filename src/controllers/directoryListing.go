package controllers

import (
    "net/http"
    "net/url"
    "path"
    "os"
    "fmt"
    "io"
    "time"
    "mime"
    "compress/gzip"
    "compress/zlib"
    "strconv"
    "SimpleHTTPServer-golang/src/utils"
    "container/list"
    "html/template"
    "strings"
    "log"
    "sort"

    "github.com/GeertJohan/go.rice"
    "github.com/yeka/zip"
)

const serverUA = ""
const fs_maxbufsize = 4096 // 4096 bits = default page size on OSX

func handleFile(w http.ResponseWriter, req *http.Request) {
        w.Header().Set("Server", serverUA)

        filepath := path.Join((*Root_folder), path.Clean(req.URL.Path))
        serveFile(filepath, w, req)

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
                log.Println("404 Not Found : Error while opening the file "+ filePath)
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
                    newFile := strings.Split(req.URL.String(),"?")
                    newRequest, _ := http.NewRequest("GET", "http://"+req.Host+newFile[0], nil)

                    // Serve the new file (encrypted zip)
                    serveFile(zipFilePath, w, newRequest)

                    // Remove the zip file
                    os.Remove(zipFilePath)

                    return

                }else if errParseQuery == nil && len(query["dlenc"]) > 0{
                    zipFilePath := utils.ZipDirectory(f, true)
                    // Generate the request for the new file - remove ?dl to download the file

                    newFile := strings.Split(req.URL.String(),"?")
                    newRequest, _ := http.NewRequest("GET", "http://"+req.Host+newFile[0], nil)

                    // Serve the new file (encrypted zip)
                    serveFile(zipFilePath, w, newRequest)

                    // Remove the zip file
                    os.Remove(zipFilePath)
                    return
                }else{
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
        }else if errParseQuery == nil && len(query["dlenc"]) > 0{ // Download the file as an encrypted zip

                // Absolute path to the file
                filePathName := f.Name()
                // Create the zip file
                zipFile, err := os.Create(filePathName+".zip")
                if err != nil {
                    log.Fatalln(err)
                }
                zipFilePath := zipFile.Name()
                zipw := zip.NewWriter(zipFile)

                // Add file f to the zip 
                utils.AddfiletoZip(statinfo.Name(),f, zipw, true, "infected")

                // Manually close the zip
                zipw.Close()

                // Generate the request for the new file
                newFile := strings.Split(req.URL.String(),"?")
                newRequest, _ := http.NewRequest("GET", "http://"+req.Host+newFile[0], nil)

                // Serve the new file (encrypted zip)
                serveFile(zipFilePath ,w , newRequest)
                os.Remove(zipFilePath)
                return
        } else {
                // Fetching file's mimetype and giving it to the browser
                if mimetype := mime.TypeByExtension(path.Ext(filePath)); mimetype != "" {
                        w.Header().Set("Content-Type", mimetype)
                } else {
                        w.Header().Set("Content-Type", "application/octet-stream")
                }
        }
        w.Header().Set("Cache-Control", "store, public, min-age=5, max-age=120")
        // Manage Content-Range (TODO: Manage end byte and multiple Content-Range)
        if req.Header.Get("Range") != "" {
                start_byte := utils.ParseRange(req.Header.Get("Range"))

                if start_byte < statinfo.Size() {
                        f.Seek(start_byte, 0)
                } else {
                        start_byte = 0
                }

                w.Header().Set("Content-Range",
                        fmt.Sprintf("bytes %d-%d/%d", start_byte, statinfo.Size()-1, statinfo.Size()))
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

func handleDirectory(f *os.File, w http.ResponseWriter, req *http.Request) {
        names, _ := f.Readdir(-1)

        // First, check if there is any index in this folder.
        for _, val := range names {
                if val.Name() == "index.html" {
                        serveFile(path.Join(f.Name(), "index.html"), w, req)
                        return
                }
        }

        // Otherwise, generate folder content.
        children_dir_tmp := list.New()
        children_files_tmp := list.New()

        for _, val := range names {
                //if val.Name()[0] == '.' {
                //        continue
                //} // Remove hidden files from listing

                if val.IsDir() {
                        children_dir_tmp.PushBack(val.Name())
                } else {
                        children_files_tmp.PushBack(val.Name())
                }
        }

        // And transfer the content to the final array structure
        children_dir := utils.CopyToArray(children_dir_tmp)
        children_files := utils.CopyToArray(children_files_tmp)
        //Sort children_dir and children_files
        sort.Slice(children_dir, func(i, j int) bool { return children_dir[i] < children_dir[j] })

        //Sort children_dir and children_files
        sort.Slice(children_files, func(i, j int) bool { return children_files[i] < children_files[j] })

        data := utils.Dirlisting{Name: req.URL.Path, ServerUA: serverUA,
                Children_dir: children_dir, Children_files: children_files}
	    err := renderTemplate(w,"directoryListing.tpl",data)
	    if err != nil {
		    fmt.Println(err)
	}
}

func renderTemplate(w http.ResponseWriter, view string,  data interface{}) error{
        templateBox, err := rice.FindBox("../views/")
        if err != nil {
            log.Fatal(err)
        }
        // get file contents as string
        templateString, err := templateBox.String(view)
        if err != nil {
            log.Fatal(err)
        }
        //tpl := template.Must(template.Parse(templateString))
        tpl, err := template.New("tpl").Parse(templateString)
        if err != nil {
                http.Error(w, "500 Internal Error : Error while generating directory listing.", 500)
                fmt.Println(err)
                return err
        }

        err = tpl.Execute(w,data)
        if err != nil {
                return err
        }
	return nil

}
