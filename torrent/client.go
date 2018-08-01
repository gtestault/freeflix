package torrent

import (
	"bytes"
	"fmt"
	"github.com/anacrolix/torrent"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
	"time"
)

//modified MIT licensed code from https://github.com/Sioro-Neoku/go-peerflix/blob/master/client.go

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
	Client  *torrent.Client
	Torrent *torrent.Torrent
	Config  ClientConfig
}

type ClientConfig struct {
	TorrentPort    int
	TorrentPath    string
	Seed           bool
	TCP            bool
	MaxConnections int
}

// NewClientConfig creates a new default configuration.
func NewClientConfig() ClientConfig {
	return ClientConfig{
		Seed:           false,
		TCP:            true,
		MaxConnections: 200,
	}
}

// NewClient creates a new torrent client based on a magnet or a torrent file.
// If the torrent file is on http, we try downloading it.
func NewClient(cfg ClientConfig) (client Client, err error) {
	var t *torrent.Torrent
	var c *torrent.Client

	client.Config = cfg

	// Create client.
	c, err = torrent.NewClient(&torrent.ClientConfig{
		DataDir:    os.TempDir(),
		NoUpload:   !cfg.Seed,
		Seed:       cfg.Seed,
		DisableTCP: !cfg.TCP,
	})

	if err != nil {
		return client, fmt.Errorf("creating torrent client failed: %v", err)
	}

	client.Client = c

	// Add torrent.

	// Add as magnet url.
	if t, err = c.AddMagnet(cfg.TorrentPath); err != nil {
		return client, fmt.Errorf("adding torrent failed: %v", err)
	}
	//// Add torrent file
	//
	//// If it's online, we try downloading the file.
	//		if isHTTP.MatchString(cfg.TorrentPath) {
	//			if cfg.TorrentPath, err = downloadFile(cfg.TorrentPath); err != nil {
	//				return client, ClientError{Type: "downloading torrent file", Origin: err}
	//			}
	//		}
	//
	//		if t, err = c.AddTorrentFromFile(cfg.TorrentPath); err != nil {
	//			return client, ClientError{Type: "adding torrent to the client", Origin: err}
	//		}

	client.Torrent = t
	client.Torrent.SetMaxEstablishedConns(cfg.MaxConnections)

	go func() {
		<-t.GotInfo()
		t.DownloadAll()

		// Prioritize first 5% of the file.
		client.getLargestFile().Torrent().DownloadPieces(0, t.NumPieces()/100*5)
	}()
	return
}

func (c Client) getLargestFile() *torrent.File {
	var target torrent.File
	var maxSize int64

	for _, file := range c.Torrent.Files() {
		if maxSize < file.Length() {
			maxSize = file.Length()
			target = file
		}
	}

	return &target
}

// GetFile is an http handler to serve the biggest file managed by the client.
func (c Client) GetFile(w http.ResponseWriter, r *http.Request) {
	target := c.getLargestFile()
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

	w.Header().Set("Content-Disposition", "attachment; filename=\""+c.Torrent.Info().Name+"\"")
	http.ServeContent(w, r, target.DisplayPath(), time.Now(), entry)
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
