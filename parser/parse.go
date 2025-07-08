package parser

import (
	"crypto/sha1"
	"errors"

	"github.com/AcidOP/torrly/bencode"
	"github.com/AcidOP/torrly/types"
)

// Parse bencoded torrent metadata into structured data
func ParseTorrentMeta(bcode bencode.BValue) (*types.TorrentMeta, error) {
	// Dictionary to hold the torrent metadata
	// Torrent root must be a dictionary
	rootMap, ok := bcode.(map[string]bencode.BValue)
	if !ok {
		return nil, errors.New("torrent root is not a dictionary")
	}

	announce, ok := rootMap["announce"].(string)
	if !ok {
		return nil, errors.New("missing or invalid announce field")
	}

	// Info is another dictionary within the torrent metadata
	infoDict, ok := rootMap["info"].(map[string]bencode.BValue)
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

	// SHA-1 hash to verify the integrity of the torrent metadata
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
