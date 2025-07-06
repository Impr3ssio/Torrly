package main

import "github.com/AcidOP/torrly/torrent"

func main() {
	encoded := "i108e"

	torrly := torrent.NewTorrly(encoded)
	torrly.Run()
}
