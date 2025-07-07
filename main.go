package main

import "github.com/AcidOP/torrly/torrent"

func main() {
	encoded := "izze"

	torrly := torrent.NewTorrly(encoded)
	torrly.Run()
}
