package main

import (
	"encoding/json"
	"fmt"
	"freeflix/service"
	"freeflix/torrent"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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
	r.PathPrefix("/assets").Handler(http.FileServer(http.Dir("./public/freeflix")))
	r.PathPrefix("/dist").Handler(http.FileServer(http.Dir("./public/freeflix")))

	//redirect unmatched paths for processing by the angular router
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/freeflix/index.html")
	})
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

	err = json.NewEncoder(w).Encode(yts.MoviePage(page, query, rating))
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
