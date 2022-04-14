package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/theotheradamsmith/verbose-octo-eureka/image"
	"github.com/theotheradamsmith/verbose-octo-eureka/logic"
)

func handleUploadPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Handling POST")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	const maxUploadSize = 1024 * 1024
	r.ParseMultipartForm(maxUploadSize)
	var err string = ""

	file, _, ok := r.FormFile("file")
	if ok != nil {
		fmt.Fprintf(w, "<h1>File Error</h1><p>%s</p>", ok)
		err = ok.Error()
	} else {
		defer file.Close()
		if object, ok := image.Decode(file); ok != nil {
			err = ok.Error()
		} else {
			fmt.Fprintf(w, object)
			if _, ok = logic.Check(object); ok != nil {
				err = ok.Error()
			} else {
				fmt.Fprintf(w, "<h1>Congratulations!</h1><p>You have solved the puzzle!</p>")
			}
		}
	}
	fmt.Fprintf(w, err)
	//render(w, r, homepageTpl, "homepage_view", data)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", handleUploadPost).Methods(http.MethodPost)
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(":8000", router))
}
