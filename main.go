package main

import (
	"fmt"

	"github.com/AcidOP/torrly/torrent"
)

func main() {
	t, err := torrent.NewTorrentFromFile("./test.torrent")
	if err != nil {
		panic(err)
	}

	t.ViewTorrent()

	peers, err := t.GetPeerList()
	if err != nil {
		panic(err)
	}

	fmt.Println("Peers:", peers)
}
