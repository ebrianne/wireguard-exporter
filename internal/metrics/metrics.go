package metrics

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// WireguardDeviceInfo - Metadata about a device
	WireguardDeviceInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "wireguard_device_info",
			Namespace: "wireguard",
			Help:      "Metadata about a device",
		},
		[]string{"device", "public_key"},
	)
)

// Init initializes all Prometheus metrics made available by AdGuard  exporter.
func Init() {
	initMetric("wireguard_device_info", WireguardDeviceInfo)
}

func initMetric(name string, metric *prometheus.GaugeVec) {
	prometheus.MustRegister(metric)
	log.Printf("New Prometheus metric registered: %s", name)
}
