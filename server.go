package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"freeflix/service"
	"freeflix/torrent"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

var yts *service.Yts
var client *torrent.Client

func init() {
	yts = service.NewClientYTS()
	var err error
	client, err = torrent.NewClient()
	if err != nil {
		log.Fatalf(err.Error())
		os.Exit(1)
	}
}

//hookedResponseWriter hijacks the behavior of the file server if it tries to return 404
//we serve the index file instead so that the angular router handles routing.
type hookedResponseWriter struct {
	http.ResponseWriter
	ignore bool
}

func (hrw *hookedResponseWriter) WriteHeader(status int) {
	if status == 404 {
		hrw.ignore = true
		indexFile, err := ioutil.ReadFile("./public/index.html")
		if err != nil {
			log.Error("HookedResponseWriter: couldn't read index.html %v", err)
		}
		b := bytes.NewBuffer(indexFile)
		hrw.ResponseWriter.Header().Set("Content-type", "text/html")
		hrw.ResponseWriter.WriteHeader(http.StatusOK)
		if _, err := b.WriteTo(hrw.ResponseWriter); err != nil {
			log.Error("HookedResponseWriter: couldn't send index.html to client %v", err)
		}
	}
	hrw.ResponseWriter.WriteHeader(status)
}

func (hrw *hookedResponseWriter) Write(p []byte) (int, error) {
	if hrw.ignore {
		return len(p), nil
	}
	return hrw.ResponseWriter.Write(p)
}

//NotFoundHook is a special HTTP handler using a hooked ResponseWriter instead of the default ResponseWriter.
type NotFoundHook struct {
	h http.Handler
}

func (nfh NotFoundHook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	nfh.h.ServeHTTP(&hookedResponseWriter{ResponseWriter: w}, r)
}

//StartServer starts the freeflix Server. Serving static content and API.
func StartServer() {
	r := mux.NewRouter()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS", "DELETE"})
	r.HandleFunc("/api/yts", getYtsMovies).Methods("GET")
	r.HandleFunc("/api/movie/watch", client.GetFile)
	r.HandleFunc("/api/movie/request", client.MovieRequest).Methods("GET")
	r.HandleFunc("/api/movie/status", client.TorrentStatus)
	r.HandleFunc("/api/movie/delete", client.MovieDelete).Methods("DELETE")
	//TODO: Access Control
	r.HandleFunc("/monitoring/status", client.Status)
	r.PathPrefix("/").Handler(NotFoundHook{http.FileServer(http.Dir("./public/"))})
	log.Debug("Listening on port 8080")
	if err := http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(r)); err != nil {
		panic(err)
	}
}

func getYtsMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//query is search term for movies
	query, err := getParam(r, "query")
	//rating is minimum imdb
	rating, err := getParam(r, "rating")
	page, err := getParam(r, "page")
	sortBy, err := getParam(r, "sort_by")
	orderBy, err := getParam(r, "order_by")

	moviePage, err := yts.MoviePage(page, query, rating, sortBy, orderBy)
	if err != nil {
		http.Error(w, "yts service offline", http.StatusServiceUnavailable)
		log.Error(err)
	}
	err = json.NewEncoder(w).Encode(moviePage)
	if err != nil {
		log.WithError(err).Error("encoding YtsPage failed")
		http.Error(w, ":whale:", http.StatusInternalServerError)
	}
}

func getParam(r *http.Request, param string) (string, error) {
	packed, ok := r.URL.Query()[param]
	if !ok || len(packed) < 1 {
		return "", fmt.Errorf("getParam(%s): no infoHash in Request", param)
	}
	return packed[0], nil
}
