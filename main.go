package main

import "github.com/AcidOP/torrly/torrent"

func main() {
	t, err := torrent.NewTorrentFromFile("./test.torrent")
	if err != nil {
		panic(err)
	}

	t.ViewTorrent()
	t.StartDownload()
}
