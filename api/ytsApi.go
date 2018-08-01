package api

import (
	"encoding/json"
	"freeflix/service"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var yts *service.Yts

func init() {
	yts = service.NewClientYTS()
}

func StartServer() {
	http.HandleFunc("/api/yts", getYtsMovies)
	log.Info("Listening on port 8080")
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
