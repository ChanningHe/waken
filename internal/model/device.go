package model

import (
	"fmt"
	"hash/crc32"
	"net"
	"strings"
	"time"
)

type Device struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	MAC           string    `json:"mac"`
	BroadcastAddr string    `json:"broadcast_addr"`
	Port          int       `json:"port"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateDeviceRequest struct {
	Name          string `json:"name"`
	MAC           string `json:"mac"`
	BroadcastAddr string `json:"broadcast_addr,omitempty"`
	Port          int    `json:"port,omitempty"`
}

type WakeRequest struct {
	MAC           string `json:"mac"`
	BroadcastAddr string `json:"broadcast_addr,omitempty"`
	Port          int    `json:"port,omitempty"`
}

// DeviceID computes a deterministic 8-char hex ID from a MAC address.
// Same MAC always produces the same ID.
func DeviceID(mac string) string {
	normalized := strings.ToUpper(
		strings.ReplaceAll(strings.ReplaceAll(mac, ":", ""), "-", ""),
	)
	h := crc32.ChecksumIEEE([]byte(normalized))
	return fmt.Sprintf("%08x", h)
}

func (r *CreateDeviceRequest) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(r.Name) > 255 {
		return fmt.Errorf("name must be 255 characters or less")
	}
	if _, err := net.ParseMAC(r.MAC); err != nil {
		return fmt.Errorf("invalid MAC address: %w", err)
	}
	if r.BroadcastAddr != "" {
		if ip := net.ParseIP(r.BroadcastAddr); ip == nil {
			return fmt.Errorf("invalid broadcast address")
		}
	}
	if r.Port != 0 && (r.Port < 1 || r.Port > 65535) {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	return nil
}

func (r *WakeRequest) Validate() error {
	if _, err := net.ParseMAC(r.MAC); err != nil {
		return fmt.Errorf("invalid MAC address: %w", err)
	}
	if r.BroadcastAddr != "" {
		if ip := net.ParseIP(r.BroadcastAddr); ip == nil {
			return fmt.Errorf("invalid broadcast address")
		}
	}
	if r.Port != 0 && (r.Port < 1 || r.Port > 65535) {
		return fmt.Errorf("port must be between 1 and 65535")
	}
	return nil
}
