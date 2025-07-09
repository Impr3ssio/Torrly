package torrent

import (
	"fmt"
	"os"
	"strings"

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

	tf, err := metaInfoFromFile(f)
	if err != nil {
		return nil, err
	}

	return &Torrent{
		TorrentMeta: tf,
	}, nil
}

// Visualize information about the torrent.
// (e.g. announce URL, file name, size, piece length, number of pieces, info hash)
func (t *Torrent) ViewTorrent() {
	var displaySize string
	s := t.Info.Length / (1024 * 1024)

	// Format the size in GB, MB, or KB
	if s >= 1024 {
		displaySize = fmt.Sprintf("%.2f GB", float64(s)/1024)
	} else if s >= 1 {
		displaySize = fmt.Sprintf("%.2f MB", float64(s))
	} else {
		displaySize = fmt.Sprintf("%d KB", t.Info.Length/1024)
	}

	fmt.Printf("\n\nAnnounce: %s\n", t.Announce)
	fmt.Printf("File name: %s\n", t.Info.Name)
	fmt.Printf("File size: %s\n", displaySize)
	fmt.Printf("Piece length: %d KB\n", t.Info.PieceLength/1024)
	fmt.Printf("Num pieces: %d\n", t.Info.PieceLength/20)
	fmt.Printf("Info Hash: %x\n", t.Info.InfoHash)
}

func (t *Torrent) StartDownload() {
	d, err := t.RequestPeers()
	if err != nil {
		panic(err)
	}

	peers, err := parser.ParsePeers(d)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nFound %d peers\n", len(peers))
}

// Takes a path as an argument and checks if the file is a .torrent file
// If torrent, reads the file and returns the contents of the file
func parseTorrentFromPath(path string) (string, error) {
	f, err := os.Stat(path)
	if os.IsNotExist(err) {
		return "", err
	}

	if strings.Split(f.Name(), ".")[1] != "torrent" {
		return "", fmt.Errorf("file %s is not a .torrent file", f.Name())
	}

	file, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(file), nil
}

// Takes a raw string (NOT decoded) and parses it into a valid torrent metadata
// Returns a TorrentMeta struct with the parsed data
func metaInfoFromFile(b string) (*types.TorrentMeta, error) {
	// Decode the bencoded string
	// Returns map[string]interface{}
	bcode, err := bencode.DecodeBencode(b)
	if err != nil {
		panic(err)
	}

	// Parse into valid torrent metadata
	meta, err := parser.ParseTorrentMeta(bcode)
	if err != nil {
		return nil, err
	}

	return meta, nil
}
