# Wireguard Prometheus Exporter

![Build/Push (master)](https://github.com/ebrianne/wireguard-exporter/workflows/Build/Push%20(master)/badge.svg?branch=master)
[![GoDoc](https://godoc.org/github.com/ebrianne/wireguard-exporter?status.png)](https://godoc.org/github.com/ebrianne/wireguard-exporter)
[![GoReportCard](https://goreportcard.com/badge/github.com/ebrianne/wireguard-exporter)](https://goreportcard.com/report/github.com/ebrianne/wireguard-exporter)

This is a Prometheus exporter for [Wireguard](https://www.wireguard.com).

<!-- ![Grafana dashboard](https://raw.githubusercontent.com/ebrianne/adguard-exporter/master/grafana/dashboard.png)

Grafana dashboard is [available here](https://grafana.com/dashboards/13330) on the Grafana dashboard website and also [here](https://raw.githubusercontent.com/ebrianne/adguard-exporter/master/grafana/dashboard.json) on the GitHub repository. -->

## Prerequisites

* [Go](https://golang.org/doc/)
* [Wireguard](https://www.wireguard.com)
* Sudo rights

## Installation

### For linux based systems

The wireguard package can be installed easily via the package repository. See [here](https://www.wireguard.com/install/) for how to install the package for your linux distro.

### For MacOS

Wireguard provides an app for MacOS. This exporter is not compatible with it. You can install wireguard CLI with [brew](https://brew.sh).

```bash
brew install wireguard-go wireguard-tools
```

### Download binary

You can download the latest version of the binary built for your architecture here:

* Architecture **i386** [
[Darwin](https://github.com/ebrianne/wireguard-exporter/releases/latest/download/wireguard_exporter-darwin-386) /
[FreeBSD](https://github.com/ebrianne/wireguard-exporter/releases/latest/download/wireguard_exporter-freebsd-386) /
[Linux](https://github.com/ebrianne/wireguard-exporter/releases/latest/download/wireguard_exporter-linux-386) /
[Windows](https://github.com/ebrianne/wireguard-exporter/releases/latest/download/wireguard_exporter-windows-386.exe)
]
* Architecture **amd64** [
[Darwin](https://github.com/ebrianne/wireguard-exporter/releases/latest/download/wireguard_exporter-darwin-amd64) /
[FreeBSD](https://github.com/ebrianne/wireguard-exporter/releases/latest/download/wireguard_exporter-freebsd-amd64) /
[Linux](https://github.com/ebrianne/wireguard-exporter/releases/latest/download/wireguard_exporter-linux-amd64) /
[Windows](https://github.com/ebrianne/wireguard-exporter/releases/latest/download/wireguard_exporter-windows-amd64.exe)
]
* Architecture **arm** [
[Linux](https://github.com/ebrianne/wireguard-exporter/releases/latest/download/wireguard_exporter-linux-arm)
]
* Architecture **arm64** [
[Linux](https://github.com/ebrianne/wireguard-exporter/releases/latest/download/wireguard_exporter-linux-arm64)
]

### From sources

Optionally, you can download and build it from the sources. You have to retrieve the project sources by using one of the following way:
```bash
$ go get -u github.com/ebrianne/wireguard-exporter
# or
$ git clone https://github.com/ebrianne/wireguard-exporter.git
```

Install the needed vendors:

```
$ GO111MODULE=on go mod vendor
```

Then, build the binary (here, an example to run on Raspberry PI ARM architecture):
```bash
$ GOOS=linux GOARCH=arm GOARM=7 go build -o wireguard_exporter .
```

## Usage

In order to run the exporter, type the following command (arguments are optional):

```bash
$ sudo ./wireguard_exporter -server_port 9586 -interval 10s
```

```bash
2020/12/27 14:26:26 ---------------------------------------
2020/12/27 14:26:26 - Wireguard exporter configuration -
2020/12/27 14:26:26 ---------------------------------------
2020/12/27 14:26:26 ServerPort : 9586
2020/12/27 14:26:26 Interval : 10s
2020/12/27 14:26:26 ---------------------------------------
2020/12/27 14:26:26 New Prometheus metric registered: wireguard_device_info
2020/12/27 14:26:26 New Prometheus metric registered: wireguard_peer_info
2020/12/27 14:26:26 New Prometheus metric registered: wireguard_peer_receive_bytes
2020/12/27 14:26:26 New Prometheus metric registered: wireguard_peer_transmit_bytes
2020/12/27 14:26:26 New Prometheus metric registered: wireguard_peer_last_handshake_seconds
2020/12/27 14:26:26 Starting HTTP server
2020/12/27 14:26:36 Found device wg0 with public key ********
2020/12/27 14:26:36 New tick of statistics
2020/12/27 14:26:36 Getting stats for device wg0, pub key ********
2020/12/27 14:26:36 Peer ********, Received 65.1M, Sent 410.2M, Last Handshake 2020-12-27 14:26:13 +0000 UTC
```

Once the exporter is running, you also have to update your `prometheus.yml` configuration to let it scrape the exporter:

```yaml
scrape_configs:
  - job_name: 'wireguard'
  static_configs:
  - targets: ['localhost:9586']
```

## Available CLI options
```bash
# Interval of time the exporter will fetch data from Adguard
-interval duration (optional) (default 10s)

# Port to be used for the exporter
-server_port string (optional) (default "9586")
```

## Available Prometheus metrics

| Metric name                           | Description                                                 |
|:-------------------------------------:|-------------------------------------------------------------|
| wireguard_device_info                 | Metadata about a device                                     |
| wireguard_peer_info                   | Metadata about a peer                                       |
| wireguard_peer_receive_bytes          | Number of bytes received from a given peer                  |
| wireguard_peer_transmit_bytes         | Number of bytes transmitted to a given peer                 |
| wireguard_peer_last_handshake_seconds | UNIX timestamp for the last handshake with a given peer     |


## Systemd file 

### Ubuntu

One can enable the program to work at startup by writing a systemd file. You can put this file in /etc/systemd/system/wireguard-exporter.service

```
[Unit]
Description=Wireguard-Exporter
After=syslog.target network-online.target

[Service]
ExecStart=/opt/wireguard_exporter/wireguard_exporter-linux-arm
Restart=on-failure
RestartSec=10s

[Install]
WantedBy=multi-user.target
```

Then do this command to start the service:
```
$ sudo systemctl start wireguard-exporter.service
```
To enable the service at startup:
```
$ sudo systemctl enable wireguard-exporter.service
```
