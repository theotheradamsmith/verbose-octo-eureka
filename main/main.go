package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/theotheradamsmith/verbose-octo-eureka/image"
	"github.com/theotheradamsmith/verbose-octo-eureka/logic"
)

const indexHTML = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta http-equiv="X-UA-Compatible" content="ie=edge" />
    <title>File upload demo</title>
  </head>
  <body>
    <form
      id="form"
      enctype="multipart/form-data"
      action="/upload"
      method="POST"
    >
      <input class="input file-input" type="file" name="file" multiple />
      <button class="button" type="submit">Submit</button>
    </form>
  </body>
</html>
`

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	path = filepath.Join(h.staticPath, path)
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func uploadPost(w http.ResponseWriter, r *http.Request) {
	const maxUploadSize = 1024 * 1024
	w.WriteHeader(http.StatusOK)
	r.ParseMultipartForm(maxUploadSize)

	file, _, ok := r.FormFile("file")
	if ok != nil {
		fmt.Fprintf(w, "<h1>File Error</h1><p>%s</p>", ok)
		fmt.Fprintf(w, indexHTML)
		return
	}
	defer file.Close()
	object, ok := image.Decode(file)
	if ok != nil {
		fmt.Fprintf(w, "<h1>Decode Error</h1><p>%s</p>", ok)
		fmt.Fprintf(w, indexHTML)
		return
	}
	if _, ok = logic.Check(object); ok != nil {
		fmt.Fprintf(w, "<h1>Logic Check Error</h1><p>%s</p>", ok)
		fmt.Fprintf(w, indexHTML)
		return
	} else {
		fmt.Fprintf(w, "<h1>Congratulations!</h1><p>You have solved the puzzle!</p>")
	}
}

func main() {
	/*
		fmt.Println("Hello, CTF!")
		pFlag := flag.String("path", "", "path of the image to decode")
		flag.Parse()
		if *pFlag != "" {
			f, ok := os.Open(*pFlag)
			if ok != nil {
				fmt.Println(ok)
			}
			defer f.Close()
			object, ok := image.Decode(f)
			if ok != nil {
				fmt.Println(ok)
				return
			}
			if _, ok := logic.Check(object); ok != nil {
				fmt.Println(ok)
			} else {
				fmt.Println("Congratulations! You have solved the puzzle!")
			}
		} else {

		}
	*/
	// server mode!
	router := mux.NewRouter()

	router.HandleFunc("/upload", uploadPost).Methods(http.MethodPost)

	spa := spaHandler{staticPath: "/var/www/build", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)
	srv := &http.Server{
		Handler:      router,
		Addr:         ":80",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
	//log.Fatal(http.ListenAndServe(":80", router))
	// read from .config file
}
