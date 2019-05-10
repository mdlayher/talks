package main

import (
	"net"

	"github.com/mdlayher/ethernet"
)

func main() {
	// START OMIT
	// The frame to be sent over the network.
	f := &ethernet.Frame{
		// Broadcast frame to all machines on same network segment.
		// Destination ff:ff:ff:ff:ff:ff
		Destination: ethernet.Broadcast,

		// Identify our machine as the sender.
		Source: net.HardwareAddr{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad},

		// Tag traffic to VLAN 10.
		VLAN: &ethernet.VLAN{
			ID: 10,
		},

		// Identify frame with an unused EtherType.
		EtherType: 0xcccc,

		// Send a simple message.
		Payload: []byte("hello world"),
	}
	// END OMIT

	_ = f
}
