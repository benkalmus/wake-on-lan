package main

import (
	"bytes"
	"net"
	"strconv"
	"testing"
	"time"
)

const testPort = 9999

// Helper function to generate the WoL magic packet
func createWoLPacket(mac string) ([]byte, error) {
	hwAddr, err := net.ParseMAC(mac)
	if err != nil {
		return nil, err
	}

	// WoL packet structure: 6 bytes of 0xFF followed by MAC address repeated 16 times
	packet := make([]byte, 102)
	for i := 0; i < 6; i++ {
		packet[i] = 0xFF
	}
	for i := 6; i < 102; i += 6 {
		copy(packet[i:i+6], hwAddr)
	}
	return packet, nil
}

// Test function to verify WoL packet send/receive
func TestWakeOnLan(t *testing.T) {
	macAddr := "AA:BB:CC:DD:EE:FF"
	expectedPacket, err := createWoLPacket(macAddr)
	if err != nil {
		t.Fatalf("Failed to create test WoL packet: %v", err)
	}

	// Start a UDP listener to catch our WoL packet
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(testPort))
	if err != nil {
		t.Fatalf("Failed to resolve UDP address: %v", err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		t.Fatalf("Failed to listen on UDP port 9: %v", err)
	}
	defer conn.Close()

	// Goroutine to receive the packet
	packetChan := make(chan []byte, 1)
	go func() {
		buf := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			t.Errorf("Failed to read from UDP: %v", err)
		}
		packetChan <- buf[:n]
	}()

	// Give some time for the listener to start
	time.Sleep(500 * time.Millisecond)

	// Send the WoL packet
	err = sendWakeOnLan(macAddr, testPort)
	if err != nil {
		t.Fatalf("sendWakeOnLan failed: %v", err)
	}

	// Wait for the packet
	select {
	case receivedPacket := <-packetChan:
		if !bytes.Equal(receivedPacket, expectedPacket) {
			t.Errorf("Received incorrect packet:\nExpected: %v\nGot: %v", expectedPacket, receivedPacket)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout: Did not receive Wake-on-LAN packet")
	}
}
