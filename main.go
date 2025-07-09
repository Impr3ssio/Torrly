package main

import "github.com/AcidOP/torrly/torrent"

func main() {
	t, err := torrent.NewTorrentFromFile("./debian.torrent")
	if err != nil {
		panic(err)
	}

	t.ViewTorrent()
	t.StartDownload()
}
