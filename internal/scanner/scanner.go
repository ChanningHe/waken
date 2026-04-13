package scanner

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

// Host represents a discovered device on the local network.
type Host struct {
	IP       string `json:"ip"`
	MAC      string `json:"mac"`
	Hostname string `json:"hostname"`
}

// Scan reads the system ARP table and resolves hostnames via reverse DNS.
// On Linux, it reads /proc/net/arp. Falls back to empty result on other OS.
func Scan() ([]Host, error) {
	entries, err := readARPTable()
	if err != nil {
		return nil, fmt.Errorf("read ARP table: %w", err)
	}

	hosts := resolveHostnames(entries)
	return hosts, nil
}

func readARPTable() ([]Host, error) {
	f, err := os.Open("/proc/net/arp")
	if err != nil {
		return nil, fmt.Errorf("open /proc/net/arp: %w (only Linux with host network is supported)", err)
	}
	defer f.Close()

	var hosts []Host
	s := bufio.NewScanner(f)
	s.Scan() // skip header line

	for s.Scan() {
		fields := strings.Fields(s.Text())
		if len(fields) < 6 {
			continue
		}

		ip := fields[0]
		flags := fields[2]
		mac := fields[3]

		// Skip incomplete entries (flags 0x0 means incomplete)
		if flags == "0x0" || mac == "00:00:00:00:00:00" {
			continue
		}

		hosts = append(hosts, Host{
			IP:  ip,
			MAC: strings.ToUpper(mac),
		})
	}

	return hosts, s.Err()
}

func resolveHostnames(entries []Host) []Host {
	var wg sync.WaitGroup
	result := make([]Host, len(entries))
	copy(result, entries)

	for i := range result {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			result[idx].Hostname = lookupHost(result[idx].IP)
		}(i)
	}

	wg.Wait()
	return result
}

func lookupHost(ip string) string {
	done := make(chan []string, 1)
	go func() {
		names, err := net.LookupAddr(ip)
		if err != nil || len(names) == 0 {
			done <- nil
			return
		}
		done <- names
	}()

	select {
	case names := <-done:
		if len(names) > 0 {
			return strings.TrimSuffix(names[0], ".")
		}
		return ""
	case <-time.After(2 * time.Second):
		return ""
	}
}
