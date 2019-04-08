package main

import (
	"log"
	"net"

	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
)

func main() {
	// START 1 OMIT
	// Select the eth0 interface to use for Ethernet traffic.
	ifi, err := net.InterfaceByName("eth0")
	if err != nil {
		log.Fatalf("failed to open interface: %v", err)
	}

	// Open an Ethernet socket using same EtherType as our frame.
	c, err := raw.ListenPacket(ifi, 0xcccc, nil)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer c.Close()
	// END 1 OMIT

	// START 2 OMIT
	// Craft the same frame as our previous "hello world" example.
	fb := ethernetFrameBytes("hello world")

	// Broadcast the frame to all devices on our network segment.
	addr := &raw.Addr{HardwareAddr: ethernet.Broadcast}
	if _, err := c.WriteTo(fb, addr); err != nil {
		log.Fatalf("failed to write frame: %v", err)
	}

	// Listen for incoming frames with messages up to our MTU size.
	b := make([]byte, ifi.MTU)
	n, _, err := c.ReadFrom(b)
	if err != nil {
		log.Fatalf("failed to read frame: %v", err)
	}

	// Unmarshal the newly received frame.
	var f ethernet.Frame
	if err := (&f).UnmarshalBinary(b[:n]); err != nil {
		log.Fatalf("failed to unmarshal frame: %v", err)
	}

	// END 2 OMIT
}

func ethernetFrameBytes(s string) []byte {
	return nil
}
