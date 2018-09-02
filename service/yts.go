package service

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
)

const (
	endpointYTS   = "https://yts.am/api/v2/"
	listMoviesYTS = "list_movies.json?"
)

//Yts service from website https://yts.am/
type Yts struct {
}

type ytsMoviePage struct {
	Status string
	Data   struct {
		Movies []*YtsMovie
	}
}

//YtsMovie stores a Movie Object returned from the YTS API.
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

//YtsTorrent stores torrent information of a specific YTS movie.
type YtsTorrent struct {
	Url     string
	Hash    string
	Quality string
	Seeds   int
	Peers   int
	Size    string
}

//NewClientYTS creates a new YTS Service instance.
func NewClientYTS() *Yts {
	return &Yts{}
}

//MoviePage gets a page of movies from the YTS API.
func (Yts) MoviePage(page, query, rating, sortBy, orderBy string) ([]*YtsMovie, error) {
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
	if sortBy != "" {
		v.Add("sort_by", sortBy)
	}
	if orderBy != "" {
		v.Add("order_by", orderBy)
	}
	reqURL := endpointYTS + listMoviesYTS + v.Encode()
	res, err := http.Get(reqURL)
	log.WithField("req", reqURL).Debug("Page Request to YTS")
	if err != nil {
		log.Error(err)
		return nil, err
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
		return nil, fmt.Errorf("MoviePage: error decoding json YTS response: %v", err)
	}
	return response.Data.Movies, nil
}
