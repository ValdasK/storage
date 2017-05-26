package main

import (
	"errors"
	"flag"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"log"

	"strconv"

	fr "github.com/DATA-DOG/fastroute"
)

var routes = map[string]fr.Router{
	"GET":    fr.New("/:filename", serveStaticFile),
	"POST":   fr.New("/:filename", uploadFile),
	"DELETE": fr.New("/:filename", deleteHandler),
}

var router = fr.RouterFunc(func(req *http.Request) http.Handler {
	return routes[req.Method]
})

var app = fr.RouterFunc(func(req *http.Request) http.Handler {
	if h := router.Route(req); h != nil {
		return h // routed and can be served
	}

	var allows []string
	for method, routes := range routes {
		if h := routes.Route(req); h != nil {
			allows = append(allows, method)
			fr.Recycle(req) // we will not serve it, need to recycle
		}
	}

	if len(allows) == 0 {
		return nil
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		err := errors.New("Method not allowed")
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
	})
})

const _64M = (1 << 10) * 64
const _defaultPath = "/tmp"

var storagePath = _defaultPath
var maxFileSize = int64(_64M)

func main() {

	port := flag.Int("port", 8222, "Port for connecting to application")
	flag.Int64Var(&maxFileSize, "max", _64M, "Max uploaded file size in bytes")
	flag.StringVar(&storagePath, "path", _defaultPath, "Storage path for files")

	flag.Parse()

	storageStat, err := os.Stat(storagePath)
	if err != nil {
		log.Fatal(err)
	}

	if !storageStat.IsDir() {
		log.Fatal("Storage path is not a folder")
	}

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), app))
}

func deleteHandler(res http.ResponseWriter, req *http.Request) {
	var filename = fr.Parameters(req).ByName("filename")

	if filename == "" {
		http.Error(res, "Not yet implemented", http.StatusNotImplemented)
		return
	}

	err := os.Remove(storagePath + filename)

	if err != nil {
		if os.IsNotExist(err) {
			http.Error(res, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	successResponse(res)
}

func serveStaticFile(res http.ResponseWriter, req *http.Request) {
	http.ServeFile(res, req, storagePath+fr.Parameters(req).ByName("filename"))
}

func uploadFile(res http.ResponseWriter, req *http.Request) {

	var err error

	if err = req.ParseMultipartForm(maxFileSize); nil != err {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, fheaders := range req.MultipartForm.File {
		for _, hdr := range fheaders {
			// open uploaded
			var infile multipart.File
			if infile, err = hdr.Open(); nil != err {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			// open destination
			var outfile *os.File
			if outfile, err = os.Create(storagePath + hdr.Filename); nil != err {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			if _, err = io.Copy(outfile, infile); nil != err {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}

			successResponse(res)
		}
	}
}

func successResponse(res http.ResponseWriter) {
	res.Write([]byte("ok"))
}
