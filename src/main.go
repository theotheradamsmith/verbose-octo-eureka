package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func handleUploadPost(w http.ResponseWriter, r *http.Request) {
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
		if object, ok := Decode(file); ok != nil {
			err = ok.Error()
		} else {
			//fmt.Fprintf(w, object)
			if _, ok = Check(object); ok != nil {
				err = ok.Error()
			} else {
				fmt.Fprintf(w, "<h1>Congratulations!</h1><p>You have solved the puzzle!</p>")
			}
		}
	}
	fmt.Fprintf(w, err)
}

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME/.appname")
	viper.AddConfigPath("./src")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w\n", err))
	}

	port := fmt.Sprintf(":%s", viper.GetString("port"))
	if port == ":" {
		log.Print("Error in configuration file: 'port' not found. Defaulting to 8000.")
		port = ":8000"
	}

	router := mux.NewRouter()
	router.HandleFunc("/", handleUploadPost).Methods(http.MethodPost)
	http.Handle("/", router)
	log.Print("Listening on port ", port)
	log.Fatal(http.ListenAndServe(port, router))
}
