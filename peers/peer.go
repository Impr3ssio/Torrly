package peers

import "net"

type Peer struct {
	IP     net.IP // IP address of the peer in binary format.
	Port   int    // Port number of the peer to connect to.
	PeerId string // (Optional) ID of the Peer
}
