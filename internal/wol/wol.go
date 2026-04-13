package wol

import (
	"fmt"
	"net"
)

const packetSize = 102

// NewMagicPacket builds a 102-byte WOL magic packet:
// 6 bytes of 0xFF followed by the MAC address repeated 16 times.
func NewMagicPacket(mac net.HardwareAddr) ([]byte, error) {
	if len(mac) != 6 {
		return nil, fmt.Errorf("invalid MAC address length: %d", len(mac))
	}

	packet := make([]byte, packetSize)
	for i := 0; i < 6; i++ {
		packet[i] = 0xFF
	}
	for i := 0; i < 16; i++ {
		copy(packet[6+i*6:], mac)
	}
	return packet, nil
}

// Send constructs and sends a magic packet to the given broadcast address and port.
func Send(macStr, broadcastAddr string, port int) error {
	mac, err := net.ParseMAC(macStr)
	if err != nil {
		return fmt.Errorf("parse MAC: %w", err)
	}

	packet, err := NewMagicPacket(mac)
	if err != nil {
		return fmt.Errorf("build packet: %w", err)
	}

	addr := &net.UDPAddr{
		IP:   net.ParseIP(broadcastAddr),
		Port: port,
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return fmt.Errorf("dial UDP: %w", err)
	}
	defer conn.Close()

	if _, err := conn.Write(packet); err != nil {
		return fmt.Errorf("send packet: %w", err)
	}

	return nil
}
