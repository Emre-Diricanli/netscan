package net

//  package netutil

import (
	"errors"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type ARPEntry struct {
	IP  net.IP
	MAC net.HardwareAddr
}

// Pick the NIC that routes to the target subnet
func FindInterfaceFor(ipnet *net.IPNet) (*net.Interface, net.IP, error) {
	ifs, _ := net.Interfaces()
	for _, ifc := range ifs {
		if (ifc.Flags&net.FlagUp) == 0 || (ifc.Flags&net.FlagLoopback) != 0 || len(ifc.HardwareAddr) == 0 {
			continue
		}
		addrs, _ := ifc.Addrs()
		for _, a := range addrs {
			ipa, ok := a.(*net.IPNet)
			if !ok || ipa.IP.To4() == nil {
				continue
			}
			if ipnet.Contains(ipa.IP) {
				return &ifc, ipa.IP, nil
			}
		}
	}
	return nil, nil, errors.New("no interface for subnet")
}

// ARPSweep sends who-has for each IP and collects replies; fast LAN discovery
func ARPSweep(cidr string, timeout time.Duration) (map[string]ARPEntry, error) {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil { // single IP case -> infer a /24 just to select an interface
		ip := net.ParseIP(cidr)
		if ip == nil {
			return nil, err
		}
		ipnet = &net.IPNet{IP: ip.Mask(net.CIDRMask(24, 32)), Mask: net.CIDRMask(24, 32)}
	}

	ifc, srcIP, err := FindInterfaceFor(ipnet)
	if err != nil {
		return nil, err
	}

	handle, err := pcap.OpenLive(ifc.Name, 65536, false, pcap.BlockForever)
	if err != nil {
		return nil, err
	}
	defer handle.Close()

	srcMAC := ifc.HardwareAddr
	broadcast := net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}

	eth := layers.Ethernet{
		SrcMAC:       srcMAC,
		DstMAC:       broadcast,
		EthernetType: layers.EthernetTypeARP,
	}
	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   []byte(srcMAC),
		SourceProtAddress: []byte(srcIP.To4()),
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
		// DstProtAddress will be set per target below
	}

	// build target list
	var ips []net.IP
	for ip := ipnet.IP.Mask(ipnet.Mask); ipnet.Contains(ip); incIP(ip) {
		ip4 := make(net.IP, len(ip))
		copy(ip4, ip)
		ips = append(ips, ip4)
	}
	if len(ips) >= 2 {
		ips = ips[1 : len(ips)-1] // drop network & broadcast
	}

	// send ARP requests
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{FixLengths: true}
	for _, target := range ips {
		arp.DstProtAddress = []byte(target.To4())
		buf.Clear()
		_ = gopacket.SerializeLayers(buf, opts, &eth, &arp)
		_ = handle.WritePacketData(buf.Bytes())
	}

	// collect replies
	out := make(map[string]ARPEntry)
	ps := gopacket.NewPacketSource(handle, handle.LinkType())
	deadline := time.After(timeout)

	for {
		select {
		case pkt := <-ps.Packets():
			if pkt == nil {
				continue
			}
			if al := pkt.Layer(layers.LayerTypeARP); al != nil {
				reply := al.(*layers.ARP)
				if reply.Operation == layers.ARPReply {
					ipStr := net.IP(reply.SourceProtAddress).String()
					out[ipStr] = ARPEntry{
						IP:  net.IP(reply.SourceProtAddress),
						MAC: net.HardwareAddr(reply.SourceHwAddress),
					}
				}
			}
		case <-deadline:
			return out, nil
		}
	}
}
