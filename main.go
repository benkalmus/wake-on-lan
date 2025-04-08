package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

// Map of device names to MAC addresses
var deviceNameToMAC map[string]string

type WakeHandler struct {
	wakeOnLANPort int
}

// Load JSON file into memory
func loadMACAddresses(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &deviceNameToMAC)
	if err != nil {
		return err
	}
	return nil
}

// Function to send a Wake-on-LAN magic packet
func sendWakeOnLan(macAddr string, wakeOnLANPort int) error {
	macAddr = strings.ToUpper(macAddr)
	hwAddr, err := net.ParseMAC(macAddr)
	if err != nil {
		return fmt.Errorf("invalid MAC address: %s", err)
	}

	// Create the WoL magic packet
	packet := make([]byte, 102)
	for i := 0; i < 6; i++ {
		packet[i] = 0xFF
	}
	for i := 6; i < 102; i += 6 {
		copy(packet[i:i+6], hwAddr)
	}
	broadcastAddrStr := fmt.Sprintf("255.255.255.255:%d", wakeOnLANPort)

	conn, err := net.Dial("udp", broadcastAddrStr)
	if err != nil {
		return fmt.Errorf("failed to dial UDP: %s", err)
	}
	defer conn.Close()

	n, err := conn.Write(packet)
	if err != nil && n != 102 {
		return fmt.Errorf("failed to send magic packet: %s", err)
	}

	return nil
}

// HTTP handler for /wake/{pc_name or mac}
func (h WakeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 || parts[1] != "wake" {
		http.Error(w, "Invalid request: expect /wake/<mac | name>", http.StatusBadRequest)
		return
	}
	deviceName := parts[2]

	if deviceName == "all" {
		// wake up all devices
		for _, addr := range deviceNameToMAC {
			_ = sendWakeOnLan(addr, h.wakeOnLANPort)
		}
		w.WriteHeader(http.StatusAccepted)
		return
	}

	macAddr, found := deviceNameToMAC[deviceName]
	if !found {
		// If not found, assume identifier is a MAC address
		macAddr = deviceName
	}

	err := sendWakeOnLan(macAddr, h.wakeOnLANPort)
	if err != nil {
		http.Error(w, "Failed to send WoL packet: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Sent Wake-on-LAN packet to %s : %s", deviceName, macAddr)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func start(httpPort, wakeOnLanPort int, configFile string) {
	// Load the mappings from file
	deviceNameToMAC = make(map[string]string)
	err := loadMACAddresses(configFile)
	if err != nil {
		log.Printf("Error loading MAC addresses: %s", err)
	}

	wakeHandler := WakeHandler{wakeOnLanPort}
	// Setup the HTTP server
	http.Handle("/wake/", wakeHandler)
	http.HandleFunc("GET /health", healthHandler)

	log.Printf("Starting WoL server on port %d...", httpPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil))
}

func main() {
	// Define flags
	httpPort := flag.Int("http-port", 8000, "Port for the HTTP server [8000]")
	wolPort := flag.Int("wol-port", 9, "Port for Wake-on-LAN UDP packets [9]")
	configFile := flag.String("config", "config.json", "Port for Wake-on-LAN UDP packets [9]")
	flag.Parse()
	start(*httpPort, *wolPort, *configFile)
}
