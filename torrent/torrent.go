package torrent

import (
	"fmt"
	"os"

	"github.com/AcidOP/torrly/bencode"
	"github.com/AcidOP/torrly/parser"
	"github.com/AcidOP/torrly/types"
)

type Torrent struct {
	*types.TorrentMeta
}

func NewTorrentFromFile(path string) (*Torrent, error) {
	f, err := parseTorrentFromPath(path)
	if err != nil {
		return nil, err
	}

	t, err := metaInfoFromFile(f)
	if err != nil {
		return nil, err
	}

	return &Torrent{
		TorrentMeta: t,
	}, nil
}

func parseTorrentFromPath(path string) (string, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", err
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(file), nil
}

func metaInfoFromFile(b string) (*types.TorrentMeta, error) {
	bcode, err := bencode.DecodeBencode(b)
	if err != nil {
		panic(err)
	}

	meta, err := parser.ParseTorrentMeta(bcode)
	if err != nil {
		return nil, err
	}

	return meta, nil
}

func (t *Torrent) ViewTorrent() {
	fmt.Println()
	fmt.Println()

	fmt.Printf("Announce: %s\n", t.Announce)
	fmt.Printf("File name: %s\n", t.Info.Name)
	fmt.Printf("File size: %d MB\n", t.Info.Length/(1024*1024))
	fmt.Printf("Piece length: %d KB\n", t.Info.PieceLength/1024)
	fmt.Printf("Num pieces: %d\n", t.Info.PieceLength/20)
	fmt.Printf("Info Hash: %x\n", t.Info.InfoHash)

	fmt.Println()
	fmt.Println()
}
