package service

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

const (
	endpointYTS   = "https://yts.am/api/v2/"
	listMoviesYTS = "list_movies.json?"
)

type Yts struct {
}

type ytsMoviePage struct {
	Status string
	Data   struct {
		Movies []*YtsMovie
	}
}

type YtsMovie struct {
	Id               int
	Url              string
	ImdbCode         string `json:"imdb_code"`
	Title            string
	Year             int
	Rating           float32
	Runtime          int
	Genres           []string
	Summary          string
	YTCode           string `json:"yt_trailer_code"`
	Language         string
	SmallCoverImage  string `json:"small_cover_image"`
	MediumCoverImage string `json:"medium_cover_image"`
	LargeCoverImage  string `json:"large_cover_image"`
	Torrents         []*YtsTorrent
}

type YtsTorrent struct {
	Url     string
	Hash    string
	Quality string
	Seeds   int
	Peers   int
	Size    string
}

func NewClientYTS() *Yts {
	return &Yts{}
}

func (Yts) MoviePage(page string, query string, rating string) []*YtsMovie {
	v := url.Values{}
	if page != "" {
		v.Add("page", page)
	}
	if query != "" {
		v.Add("query_term", query)
	}
	if rating != "" {
		v.Add("minimum_rating", rating)
	}
	reqURL := endpointYTS + listMoviesYTS + v.Encode()
	res, err := http.Get(reqURL)
	log.WithField("req", reqURL).Debug("Page Request to YTS")
	if err != nil {
		log.Error(err)
		return nil
	}

	dec := json.NewDecoder(res.Body)
	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Error(err)
		}
	}()
	var response ytsMoviePage
	err = dec.Decode(&response)
	if err != nil {
		log.WithError(err).Error("error decoding json YTS response: ")
	}
	return response.Data.Movies
}
