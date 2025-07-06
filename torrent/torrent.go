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

func (t *torrly) Run() error {
	d, err := bencode.DecodeBencode(t.magnet)
	if err != nil {
		return err
	}

	fmt.Println(d)
	return nil
}
