//package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/theotheradamsmith/verbose-octo-eureka/image"
	"github.com/theotheradamsmith/verbose-octo-eureka/logic"
)

type config struct {
	host         string
	readTimeout  time.Duration
	writeTimeout time.Duration
}

type htmlServer struct {
	server *http.Server
	wg     sync.WaitGroup
}

type spaHandler struct {
	staticPath string
}

var homepageTpl *template.Template

func init_assets() {
	if t, err := ioutil.ReadFile("/var/www/build/templateIndexHtml"); err != nil {
		panic(err)
	} else {
		homepageTpl = template.Must(template.New("homepage_view").Parse(string(t)))
	}
}

func render(w http.ResponseWriter, r *http.Request, tpl *template.Template, name string, data interface{}) {
	buf := new(bytes.Buffer)
	if err := tpl.ExecuteTemplate(buf, name, data); err != nil {
		fmt.Printf("\nRender Error: %v\n", err)
		return
	}
	w.Write(buf.Bytes())
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	joined_path := filepath.Join(h.staticPath, path)
	_, err = os.Stat(joined_path)
	if os.IsNotExist(err) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		render(w, r, homepageTpl, "homepage_view", map[string]interface{}{})
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if path != "/" {
		http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
	}
	render(w, r, homepageTpl, "homepage_view", map[string]interface{}{})
}

func handleUploadPost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	const maxUploadSize = 1024 * 1024
	r.ParseMultipartForm(maxUploadSize)
	var err string = ""

	file, _, ok := r.FormFile("file")
	if ok != nil {
		fmt.Fprintf(w, "<h1>File Error</h1><p>%s</p>", ok)
		err = ok.Error()
	}
	defer file.Close()
	if object, ok := image.Decode(file); ok != nil {
		err = ok.Error()
	} else {
		if _, ok = logic.Check(object); ok != nil {
			err = ok.Error()
		} else {
			fmt.Fprintf(w, "<h1>Congratulations!</h1><p>You have solved the puzzle!</p>")
		}
	}

	data := map[string]interface{}{
		"Error": err,
	}
	render(w, r, homepageTpl, "homepage_view", data)
}

func start(cfg config) *htmlServer {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	router := mux.NewRouter()
	router.HandleFunc("/upload", handleUploadPost).Methods(http.MethodPost)
	spa := spaHandler{staticPath: "/var/www/"}
	router.PathPrefix("/").Handler(spa)
	server := htmlServer{
		server: &http.Server{
			Addr:           cfg.host,
			Handler:        router,
			ReadTimeout:    cfg.readTimeout,
			WriteTimeout:   cfg.writeTimeout,
			MaxHeaderBytes: 1 << 20,
		},
	}

	server.wg.Add(1)

	go func() {
		fmt.Printf("\nServer: Service started: Host=%v\n", cfg.host)
		log.Fatal(server.server.ListenAndServe())
		server.wg.Done()
	}()

	return &server
}

func (server *htmlServer) stop() error {
	// create context to attempt graceful 5 second shutdown
	const timeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	fmt.Printf("\nServer: Service stopping\n")

	if err := server.server.Shutdown(ctx); err != nil {
		// timeout on graceful shutdown? force close
		if err := server.server.Close(); err != nil {
			fmt.Printf("\nServer: Service stopping: Error=%v\n", err)
			return err
		}
	}

	// wait for the listener to report that it is closed
	server.wg.Wait()
	fmt.Printf("\nServer: Stopped\n")
	return nil
}

func main() {
	init_assets()
	serverCfg := config{
		host:         ":80",
		readTimeout:  15 * time.Second,
		writeTimeout: 15 * time.Second,
	}

	server := start(serverCfg)
	defer server.stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
	fmt.Println("main: shutting down")
}
