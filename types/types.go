package types

// TorrentMeta represents the metadata of a torrent file.
// All the properties are required to be present in the torrent file.
// Taken from https://wiki.theory.org/BitTorrentSpecification#Bencoding

type TorrentMeta struct {
	Info     Info   // The dictionary which describes the files of the torrent.
	Announce string // The URL of the tracker where peers can be asked for files.
}

type Info struct {
	Name        string     // The name of the file or directory.
	Length      int        // length of the file in bytes (integer)
	InfoHash    []byte     // SHA1 hash of the info dictionary, used to verify the integrity of the torrent metadata.
	PieceHashes [][20]byte // Array of SHA1 hashes for each piece of the file
	PieceLength int        // The length of each piece in bytes (integer). Commonly 256 KiB.
}
