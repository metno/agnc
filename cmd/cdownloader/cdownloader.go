package main

import (
	"context"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/metno/agnc/pkg/downloader"
)

func main() {
	callstring := os.Args[1]
	log.Printf(callstring)

	ctx := context.Background()

	url, err := url.Parse(callstring)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(url.Host)
	log.Printf(filepath.Base(url.Path))

	downloader.DownloadTimestep(ctx, url.Host, url.Path, ".")

	return
}
