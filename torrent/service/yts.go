package service

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strconv"
)

const (
	endpointYTS = "https://yts.am/api/v2/"
)

type Yts struct {
}

type YtsMovie struct {
	Id              int
	Url             string
	ImdbCode        string `json:"imdb_code"`
	Title           string
	Year            int
	Rating          float32
	Runtime         int
	Genres          []string
	Summary         string
	YTCode          string `json:"yt_trailer_code"`
	Language        string
	SmallCoverImage string `json:"small_cover_image"`
	torrents        []*YtsTorrent
}

type YtsTorrent struct {
	Url     string
	Hash    string
	Quality string
	Seeds   int
	peers   int
	size    string
}

type movies []*YtsMovie

func (Yts) MoviePage(page int) {
	v := url.Values{}
	v.Add("page", strconv.Itoa(page))
	res, err := http.Get(endpointYTS + v.Encode())
	if err != nil {
		logrus.Error(err)
	}
	dec := json.NewDecoder(res.Body)
	defer func() {
		err := res.Body.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()
}
