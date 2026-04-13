package wol

import (
	"net"
	"testing"
)

func TestNewMagicPacket(t *testing.T) {
	mac, _ := net.ParseMAC("AA:BB:CC:DD:EE:FF")
	packet, err := NewMagicPacket(mac)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(packet) != 102 {
		t.Fatalf("packet length = %d, want 102", len(packet))
	}

	// First 6 bytes must be 0xFF
	for i := 0; i < 6; i++ {
		if packet[i] != 0xFF {
			t.Errorf("packet[%d] = %02x, want 0xFF", i, packet[i])
		}
	}

	// Next 96 bytes = 16 repetitions of the MAC
	for i := 0; i < 16; i++ {
		offset := 6 + i*6
		for j := 0; j < 6; j++ {
			if packet[offset+j] != mac[j] {
				t.Errorf("packet[%d] = %02x, want %02x", offset+j, packet[offset+j], mac[j])
			}
		}
	}
}

func TestNewMagicPacketInvalidMAC(t *testing.T) {
	short := net.HardwareAddr{0x01, 0x02, 0x03}
	_, err := NewMagicPacket(short)
	if err == nil {
		t.Fatal("expected error for short MAC")
	}
}

func TestNewMagicPacketFormats(t *testing.T) {
	formats := []string{
		"AA:BB:CC:DD:EE:FF",
		"AA-BB-CC-DD-EE-FF",
		"aa:bb:cc:dd:ee:ff",
	}

	for _, f := range formats {
		mac, err := net.ParseMAC(f)
		if err != nil {
			t.Fatalf("failed to parse %s: %v", f, err)
		}
		packet, err := NewMagicPacket(mac)
		if err != nil {
			t.Fatalf("failed for %s: %v", f, err)
		}
		if len(packet) != 102 {
			t.Errorf("packet length = %d for %s, want 102", len(packet), f)
		}
	}
}
