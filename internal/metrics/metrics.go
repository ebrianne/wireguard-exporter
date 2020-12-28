package metrics

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	labels = []string{"device", "public_key"}

	// WireguardDeviceInfo - Metadata about a device
	WireguardDeviceInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "device_info",
			Namespace: "wireguard",
			Help:      "Metadata about a device",
		},
		labels,
	)

	// WireguardPeerInfo - Metadata about a peer. The public_key label on peer metrics refers to the peer's public key; not the device's public key
	WireguardPeerInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "peer_info",
			Namespace: "wireguard",
			Help:      "Metadata about a peer. The public_key label on peer metrics refers to the peer's public key; not the device's public key",
		},
		append(labels, []string{"endpoint", "name"}...),
	)

	// WireguardPeerReceiveBytes - Number of bytes received from a given peer
	WireguardPeerReceiveBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "peer_receive_bytes",
			Namespace: "wireguard",
			Help:      "Number of bytes received from a given peer",
		},
		append(labels, []string{"endpoint", "name"}...),
	)

	// WireguardPeerTransmitBytes - Number of bytes transmitted to a given peer
	WireguardPeerTransmitBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "peer_transmit_bytes",
			Namespace: "wireguard",
			Help:      "Number of bytes transmitted to a given peer",
		},
		append(labels, []string{"endpoint", "name"}...),
	)

	// WireguardPeerLastHandshake - UNIX timestamp for the last handshake with a given peer
	WireguardPeerLastHandshake = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "peer_last_handshake_seconds",
			Namespace: "wireguard",
			Help:      "UNIX timestamp for the last handshake with a given peer",
		},
		append(labels, []string{"endpoint", "name"}...),
	)
)

// Init initializes all Prometheus metrics made available by wireguard exporter.
func Init() {
	initMetric("device_info", WireguardDeviceInfo)
	initMetric("peer_info", WireguardPeerInfo)
	initMetric("peer_receive_bytes", WireguardPeerReceiveBytes)
	initMetric("peer_transmit_bytes", WireguardPeerTransmitBytes)
	initMetric("peer_last_handshake_seconds", WireguardPeerLastHandshake)
}

func initMetric(name string, metric *prometheus.GaugeVec) {
	prometheus.MustRegister(metric)
	log.Printf("New Prometheus metric registered: %s", name)
}
