package api

import (
	"encoding/json"
	"freeflix/service"
	"freeflix/torrent"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var yts *service.Yts
var client torrent.Client

func init() {
	yts = service.NewClientYTS()
	var cfg = torrent.NewClientConfig()
	cfg.TorrentPath = torrent.BuildMagnet("0B2A8EAC63A94CEDF31118400F3F4BCB08B72D1A", "ReadyPlayerOne")
	var err error
	client, err = torrent.NewClient(cfg)
	if err != nil {
		log.Fatalf(err.Error())
		os.Exit(1)
	}
}

func StartServer() {
	http.HandleFunc("/api/yts", getYtsMovies)
	http.HandleFunc("/api/movie.mp4", client.GetFile)
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
