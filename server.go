package main

import (
	"freeflix/torrent/service"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{})
}

func main() {
	service.NewClientYTS().MoviePage(1)
}
