package parser

import (
	"crypto/sha1"
	"errors"
	"net"

	"github.com/AcidOP/torrly/bencode"
	"github.com/AcidOP/torrly/types"
)

// the map kept repeating, so defined a new type for my sanity
type BDict = map[string]bencode.BValue

// Parse bencoded torrent metadata into structured data
func ParseTorrentMeta(bcode bencode.BValue) (*types.TorrentMeta, error) {
	// Dictionary to hold the torrent metadata
	// Torrent root must be a dictionary
	root, ok := bcode.(BDict)
	if !ok {
		return nil, errors.New("torrent root is not a dictionary")
	}

	announce, ok := root["announce"].(string)
	if !ok {
		return nil, errors.New("missing or invalid announce field")
	}

	// Info is another dictionary within the torrent metadata
	infoDict, ok := root["info"].(BDict)
	if !ok {
		return nil, errors.New("missing or invalid info dictionary")
	}

	// Need to re-encode the info dictionary to get raw bytes
	// This is to find out the SHA-1 hash of the info dict
	// which is used to verify the integrity of the torrent metadata
	infoRawBytes, err := bencode.Encode(infoDict)
	if err != nil {
		return nil, errors.New("failed to re-encode info dict to bencode")
	}

	// SHA-1 hash the `info` to verify the integrity of the torrent metadata
	infoHash := sha1.Sum(infoRawBytes)

	name, _ := infoDict["name"].(string)
	length, _ := infoDict["length"].(int)
	pieceLength, _ := infoDict["piece length"].(int)
	piecesStr, _ := infoDict["pieces"].(string)

	// Chunks of pieces are stored as a string of concatenated SHA1 hashes
	pieces, err := splitChunks([]byte(piecesStr))
	if err != nil {
		return nil, err
	}

	// Convert each piece hash from byte slice to [20]byte
	// This is necessary because the torrent protocol
	// expects piece hashes to be in this format
	var pieceHashes [][20]byte
	for _, p := range pieces {
		var hash [20]byte
		copy(hash[:], p)
		pieceHashes = append(pieceHashes, hash)
	}

	return &types.TorrentMeta{
		Announce: announce,
		Info: types.Info{
			Name:        name,
			Length:      length,
			PieceLength: pieceLength,
			InfoHash:    infoHash[:],
			PieceHashes: pieceHashes,
		},
	}, nil
}

// Take a tracker response and return a dictionary of peers
func ParsePeers(tResponse bencode.BValue) ([]types.Peer, error) {
	// Root dictionary of peers
	root, ok := tResponse.(BDict)
	if !ok {
		return nil, errors.New("tracker response is not a valid dictionary")
	}

	// Extract the `peers` key from the response
	peerListRaw, ok := root["peers"]
	if !ok {
		return nil, errors.New("missing `peers` key")
	}

	peerList, ok := peerListRaw.([]bencode.BValue)
	if !ok {
		return nil, errors.New("`peers` is not a list")
	}

	var peers []types.Peer

	// Convert from BValue to Peer types
	for _, raw := range peerList {
		d, ok := raw.(BDict)
		if !ok {
			return nil, errors.New("peer entry is not a dictionary")
		}

		ip, ok := d["ip"].(string)
		if !ok {
			return nil, errors.New("invalid ip address")
		}

		port, ok := d["port"].(int)
		if !ok {
			return nil, errors.New("peer missing or invalid 'port'")
		}

		// Some clients do not have a peer ID (Optionally ignore this)
		peerId, _ := d["peer id"].(string)

		peers = append(peers, types.Peer{
			IP:     net.ParseIP(ip),
			Port:   port,
			PeerID: peerId,
		})
	}
	return peers, nil
}

func splitChunks(data []byte) ([][]byte, error) {
	const hashLength = 20
	var dataLength = len(data)

	var chunks [][]byte
	if dataLength%hashLength != 0 {
		return nil, errors.New("invalid piece length")
	}

	for i := 0; i < dataLength; i += hashLength {
		hash := data[i : i+hashLength]
		chunks = append(chunks, hash)
	}

	return chunks, nil
}
