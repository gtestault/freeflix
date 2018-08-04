package torrent

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var statusTemplate *template.Template

func init() {
	cwd, _ := os.Getwd()
	fp := filepath.Join(cwd, "./torrent/templates/status.html")
	statusTemplate = template.Must(
		template.
			New("status.html").
			Funcs(template.FuncMap{"progress": downloadProgressString}).
			ParseFiles(fp))
}

func (c *Client) Status(w http.ResponseWriter, r *http.Request) {
	if err := statusTemplate.Execute(w, c); err != nil {
		log.Error(fmt.Errorf("error while displaying status: %v", err))
	}
}

func downloadProgressString(completed int64, missing int64) string {
	mbCompleted := strconv.FormatInt(completed/1000000, 10)
	mbTotal := strconv.FormatInt((completed+missing)/1000000, 10)
	return fmt.Sprintf("%s / %s Mb", mbCompleted, mbTotal)
}
