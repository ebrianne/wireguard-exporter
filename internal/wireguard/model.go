package wireguard

// Stats struct containing all wireguard statistics structs
type Stats struct {
	WireguardDeviceInfo          []map[string]string
	WireguardPeerInfo            []map[string]string
	WireguardPeerReceiveBytes    []float64
	WireguardPeerTransmitBytes   []float64
	WireguardPeerLastHandshake   []float64
}