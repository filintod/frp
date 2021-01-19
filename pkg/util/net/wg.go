package net

import (
	"context"
	"fmt"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun"
	"log"
	"net"
)

// ConnectWireguardServer dials wireguard address (addr: domain:port)
func ConnectWireguardServer(localTunnelAddr string, addr string, dnsServer string, mtu int, privateKey string, wgPublicKey string, wgAddress string) (net.Conn, error) {
	tun, tnet, err := tun.CreateNetTUN(
		[]net.IP{net.ParseIP(localTunnelAddr)},
		[]net.IP{net.ParseIP(dnsServer)},
		mtu)
	if err != nil {
		log.Panic(err)
	}
	l := device.NewLogger(device.LogLevelInfo, "wg")
	dev := device.NewDevice(tun, l)
	if err = dev.IpcSet(fmt.Sprintf(
		"private_key=%s\n"+
			"public_key=%s\n"+
			"endpoint=%s\n"+
			"allowed_ip=0.0.0.0/0", privateKey, wgPublicKey, wgAddress)); err != nil {
		log.Panic(err)
	}
	dev.Up()

	return tnet.DialContext(context.Background(), "tcp", addr)
}
