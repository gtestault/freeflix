package torrent

import (
	"bytes"
	"fmt"
	"github.com/anacrolix/torrent"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

var trackers = [...]string{
	"udp://open.demonii.com:1337/announce",
	"udp://tracker.openbittorrent.com:80",
	"udp://tracker.coppersurfer.tk:6969",
	"udp://glotorrents.pw:6969/announce",
	"udp://tracker.opentrackr.org:1337/announce",
	"udp://torrent.gresille.org:80/announce",
	"udp://p4p.arenabg.com:1337",
	"udp://tracker.leechers-paradise.org:6969",
}

type Client struct {
	Client   *torrent.Client
	Torrents map[string]*Torrent
}

type Torrent struct {
	*torrent.Torrent
	Fetched bool
}

func NewClient() (client *Client, err error) {
	var c *torrent.Client
	client = &Client{}
	client.Torrents = make(map[string]*Torrent)

	//config
	torrentCfg := torrent.NewDefaultClientConfig()
	torrentCfg.Seed = false
	torrentCfg.DataDir = "./Movies"

	// Create client.
	c, err = torrent.NewClient(torrentCfg)
	if err != nil {
		return client, fmt.Errorf("creating torrent client failed: %v", err)
	}
	client.Client = c
	return
}

//Adds Torrent to the client. If the torrent is already added returns without error.
func (c *Client) AddTorrent(infoHash string) (err error) {
	//if torrent already registered in client return
	if _, ok := c.Torrents[infoHash]; ok {
		return nil
	}

	t, err := c.Client.AddMagnet(BuildMagnet(infoHash, infoHash))
	if err != nil {
		return fmt.Errorf("adding torrent failed: %v", err)
	}
	c.Torrents[infoHash] = &Torrent{Torrent: t}

	//wait for fetch to Download torrent
	go func() {
		<-t.GotInfo()
		t.DownloadAll()
		c.Torrents[infoHash].Fetched = true
	}()
	return
}

func (c *Client) getLargestFile(infoHash string) (*torrent.File, error) {
	var target *torrent.File
	var maxSize int64
	t, ok := c.Torrents[infoHash]
	if !ok {
		return nil, fmt.Errorf("error: unregistered infoHash")
	}
	for _, file := range t.Files() {
		if maxSize < file.Length() {
			maxSize = file.Length()
			target = file
		}
	}
	return target, nil
}

func (c *Client) MovieRequest(w http.ResponseWriter, r *http.Request) {
	infoHash, err := infoHashFromRequest(r)
	if err != nil {
		log.WithField("infoHash", infoHash).Warn("MovieRequest: Request without InfoHash")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.WithField("infoHash", infoHash).Debug("torrent request received")
	if err = c.AddTorrent(infoHash); err != nil {
		log.WithField("infoHash", infoHash).Error("MovieRequest adding torrent: %v", err)
	}
}

// GetFile is an http handler to serve the biggest file managed by the client.
func (c *Client) GetFile(w http.ResponseWriter, r *http.Request) {
	infoHash, err := infoHashFromRequest(r)
	if err != nil {
		log.WithField("infoHash", infoHash).Warn("GetFile: Request without InfoHash")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.WithField("infoHash", infoHash).Debug("movie file request received")

	target, err := c.getLargestFile(infoHash)
	if err != nil {
		log.WithField("infoHash", infoHash).WithError(err).Errorf("server: error getting file")
	}
	entry, err := NewFileReader(target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer func() {
		if err := entry.Close(); err != nil {
			log.Printf("Error closing file reader: %s\n", err)
		}
	}()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	log.Info(target.DisplayPath())
	http.ServeContent(w, r, target.DisplayPath(), time.Now(), entry)
}

func infoHashFromRequest(r *http.Request) (string, error) {
	packed, ok := r.URL.Query()["infoHash"]
	if !ok || len(packed) < 1 {
		return "", fmt.Errorf("infoHashFromRequest: no infoHash in Request")
	}
	return packed[0], nil
}

func BuildMagnet(infoHash string, title string) string {
	b := &bytes.Buffer{}
	b.WriteString("magnet:?xt=urn:btih:")
	b.WriteString(infoHash)
	b.WriteString("&dn=")
	b.WriteString(url.QueryEscape(title))
	for _, tracker := range trackers {
		b.WriteString("&tr=")
		b.WriteString(tracker)
	}
	return b.String()
}
