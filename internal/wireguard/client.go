package wireguard

import (
	"log"
	"time"

	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"github.com/ebrianne/wireguard-exporter/internal/metrics"
)

type Client struct {
	Devices []*wgtypes.Device
	WgPeerFile string
	interval time.Duration
}

// NewClient method initializes a new AdGuard  client.
func NewClient(file string, interval time.Duration) *Client {

	client, err := wgctrl.New()
	if err != nil {
		log.Fatalf("failed to open WireGuard control client: %v", err)
	}

	devices, err := client.Devices()
	if err != nil {
		log.Fatalf("failed to fetch WireGuard devices: %v", err)
	}

	return &Client {
		Devices: devices,
		WgPeerFile: file,
		interval: interval,
	}
}

// Scrape method authenticates and retrieves statistics from AdGuard  JSON API
// and then pass them as Prometheus metrics.
func (c *Client) Scrape() {
	for range time.Tick(c.interval) {

		stats := c.getStatistics()
		//Set the metrics
		c.setMetrics(stats)

		log.Printf("New tick of statistics")
	}
}

func (c *Client) setMetrics(stats *Stats) {
	//Stats
	metrics.WireguardDeviceInfo.WithLabelValues(stats.WireguardDeviceInfo[0], stats.WireguardDeviceInfo[1]).Set(float64(1))
}

func (c *Client) getStatistics() *Stats {

	var stats Stats
	for _, d := range c.Devices {
		stats.WireguardDeviceInfo = []string{d.Name, d.PublicKey.String()}
	}

	return &stats
}