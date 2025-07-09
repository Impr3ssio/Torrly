package torrent

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/AcidOP/torrly/bencode"
)

// Create a URL to request to the tracker for peer information
// Must be a GET request with the following:
// https://wiki.theory.org/BitTorrentSpecification#Tracker_Request_Parameters
func buildTrackerURL(t Torrent) (string, error) {
	base, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}

	params := url.Values{
		"info_hash":  []string{string(t.Info.InfoHash)},
		"peer_id":    []string{"-TR0001-123456789012"},
		"port":       []string{"6881"},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"left":       []string{strconv.Itoa(t.Info.Length)},
		"compact":    []string{"1"},
	}

	base.RawQuery = params.Encode()
	return base.String(), nil
}

// Announce to the tracker to get a list of peers
// Returns a map of peers with their IP addresses and ports
func getTrackerResponse(t Torrent) ([]byte, error) {
	trackerURL, err := buildTrackerURL(t)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(trackerURL)
	if err != nil {
		return nil, errors.New("failed to connect to tracker")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("tracker returned non-200 status code: " + resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	return data, nil
}

// Returns a list of peers from the tracker
func (t *Torrent) RequestPeers() (bencode.BValue, error) {
	res, err := getTrackerResponse(*t)
	if err != nil {
		return nil, err
	}

	decoded, err := bencode.DecodeBencode(string(res))
	if err != nil {
		return nil, fmt.Errorf("failed to decode tracker response: %v", err)
	}
	return decoded, nil
}
