package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ebrianne/wireguard-exporter/config"
	"github.com/ebrianne/wireguard-exporter/internal/metrics"
	"github.com/ebrianne/wireguard-exporter/internal/server"
	"github.com/ebrianne/wireguard-exporter/internal/wireguard"
)

const (
	name = "wireguard-exporter"
)

var (
	s *server.Server
)

func main() {
	conf := config.Load()

	metrics.Init()
	go wireguard.Scrape(conf.Interval)
	initHttpServer(conf.ServerPort)
	handleExitSignal()
}

func initHttpServer(port string) {
	s = server.NewServer(port)
	go s.ListenAndServe()
}

func handleExitSignal() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	s.Stop()
	fmt.Println(fmt.Sprintf("\n%s HTTP server stopped", name))
}
