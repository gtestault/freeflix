package main

import (
	"encoding/json"
	"freeflix/service"
	"freeflix/torrent"
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
	http.HandleFunc("/api/yts", getYtsMovies)
	http.HandleFunc("/api/movie/watch", client.GetFile)
	http.HandleFunc("/api/movie/request", client.MovieRequest)
	//TODO: Acess Control
	http.HandleFunc("/monitoring/status", client.Status)
	log.Debug("Listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func getYtsMovies(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	err := json.NewEncoder(w).Encode(yts.MoviePage(1))
	if err != nil {
		log.WithError(err).Error("encoding YtsPage failed")
		http.Error(w, ":whale:", http.StatusInternalServerError)
	}
}
