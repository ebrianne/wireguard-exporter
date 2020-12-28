package wireguard

import (
	"fmt"
	"log"
	"time"

	"github.com/ebrianne/wireguard-exporter/internal/metrics"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"

	"code.cloudfoundry.org/bytefmt"
)

// Scrape method authenticates and retrieves statistics from AdGuard  JSON API
// and then pass them as Prometheus metrics.
func Scrape(interval time.Duration) {

	for range time.Tick(interval) {

		client, err := wgctrl.New()
		if err != nil {
			log.Fatalf("failed to open WireGuard control client: %v", err)
		}
		defer client.Close()

		devices, err := client.Devices()
		if err != nil {
			log.Fatalf("failed to fetch WireGuard devices: %v", err)
		}

		for _, d := range devices {
			log.Printf(fmt.Sprintf("Found device %s with public key %s", d.Name, d.PublicKey.String()))
		}

		log.Printf("New tick of statistics")

		stats := getStatistics(devices)
		//Set the metrics
		setMetrics(stats)
	}
}

func setMetrics(stats *Stats) {

	//Stats
	for l := range stats.WireguardDeviceInfo {
		for device, pubkey := range stats.WireguardDeviceInfo[l] {

			metrics.WireguardDeviceInfo.WithLabelValues(device, pubkey).Set(float64(1))

			for m := range stats.WireguardPeerInfo {
				for endpoint, peerpub := range stats.WireguardPeerInfo[m] {

					if endpoint != "" {
						metrics.WireguardPeerInfo.WithLabelValues(device, pubkey, endpoint, peerpub).Set(float64(1))
					} else {
						metrics.WireguardPeerInfo.WithLabelValues(device, pubkey, endpoint, peerpub).Set(float64(0))
					}

					metrics.WireguardPeerReceiveBytes.WithLabelValues(device, pubkey, endpoint, peerpub).Set(float64(stats.WireguardPeerReceiveBytes[m]))
					metrics.WireguardPeerTransmitBytes.WithLabelValues(device, pubkey, endpoint, peerpub).Set(float64(stats.WireguardPeerTransmitBytes[m]))
					metrics.WireguardPeerLastHandshake.WithLabelValues(device, pubkey, endpoint, peerpub).Set(stats.WireguardPeerLastHandshake[m])
				}
			}
		}
	}
}

func getStatistics(devices []*wgtypes.Device) *Stats {

	var stats Stats
	var idev = 0

	stats.WireguardDeviceInfo = make([]map[string]string, len(devices))

	for _, d := range devices {

		log.Printf("Getting stats for device %s, pub key %s", d.Name, d.PublicKey.String())

		stats.WireguardDeviceInfo[idev] = make(map[string]string)
		stats.WireguardDeviceInfo[idev][d.Name] = d.PublicKey.String()

		var ipeer = 0

		stats.WireguardPeerInfo = make([]map[string]string, len(d.Peers))
		stats.WireguardPeerReceiveBytes = make([]float64, len(d.Peers))
		stats.WireguardPeerTransmitBytes = make([]float64, len(d.Peers))
		stats.WireguardPeerLastHandshake = make([]float64, len(d.Peers))

		for _, p := range d.Peers {
			peerPub := p.PublicKey.String()

			// Use empty string instead of special Go <nil> syntax for no endpoint.
			var endpoint string
			if p.Endpoint != nil {
				endpoint = p.Endpoint.String()
			}

			stats.WireguardPeerInfo[ipeer] = make(map[string]string)
			stats.WireguardPeerInfo[ipeer][endpoint] = peerPub
			stats.WireguardPeerReceiveBytes[ipeer] = float64(p.ReceiveBytes)
			stats.WireguardPeerTransmitBytes[ipeer] = float64(p.TransmitBytes)

			// Expose last handshake of 0 unless a last handshake time is set.
			var last float64
			if !p.LastHandshakeTime.IsZero() {
				last = float64(p.LastHandshakeTime.Unix())
			}

			stats.WireguardPeerLastHandshake[ipeer] = last

			if endpoint != "" {
				log.Printf("Peer %s, Received %v, Sent %v, Last Handshake %s", peerPub, bytefmt.ByteSize(uint64(p.ReceiveBytes)), bytefmt.ByteSize(uint64(p.TransmitBytes)), time.Unix(p.LastHandshakeTime.Unix(), 0))
			}

			ipeer++
		}

		idev++
	}

	return &stats
}
