package main

import (
	"encoding/json"
	"fmt"
	"freeflix/service"
	"freeflix/torrent"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
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
	http.HandleFunc("/api/movie/delete", client.MovieDelete)
	//TODO: Access Control
	http.HandleFunc("/monitoring/status", client.Status)
	log.Debug("Listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func getYtsMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
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

func validateInt(a string, min, max int) (int, error) {
	i, err := strconv.Atoi(a)
	if err != nil || i < min || i > max {
		return 0, fmt.Errorf("validateInt: invalid")
	}
	return i, nil
}
