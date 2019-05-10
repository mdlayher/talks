package main

import (
	"log"
	"net"

	"github.com/mdlayher/ethernet"
)

func main() {
	// START 1 OMIT
	// The frame to be sent over the network.
	f := &ethernet.Frame{
		// Broadcast frame to all machines on same network segment.
		// Destination ff:ff:ff:ff:ff:ff
		Destination: ethernet.Broadcast,

		// Identify our machine as the sender.
		Source: net.HardwareAddr{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad},

		// Identify frame with an unused EtherType.
		EtherType: 0xcccc,

		// Send a simple message.
		Payload: []byte("hello world"),
	}
	// END 1 OMIT

	// START 2 OMIT
	// Marshal to wire format.
	fb, err := f.MarshalBinary()
	if err != nil {
		log.Fatalf("failed to marshal: %v", err)
	}

	// Send frame over some network interface.
	if err := sendEthernetFrame(fb); err != nil {
		log.Fatalf("failed to send frame: %v", err)
	}
	// END 2 OMIT
}

func sendEthernetFrame(b []byte) error {
	return nil
}
