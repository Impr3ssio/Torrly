package torrent

import (
	"fmt"

	"github.com/AcidOP/torrly/bencode"
)

type torrly struct {
	magnet string
}

func NewTorrly(mag string) *torrly {
	return &torrly{magnet: mag}
}

func (t *torrly) Run() {
	d, err := bencode.DecodeBencode(t.magnet)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(d)
}
