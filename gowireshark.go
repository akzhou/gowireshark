/**
 * @Author: Administrator
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2019/12/10 9:35
 */

package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func main() {
	wireShark("eth0", uint16(12345))
}

var (
	device      = "eth0"
	snapshotLen = int32(65535)
	promiscuous = false
	err         error
	timeout     = pcap.BlockForever
	handle      *pcap.Handle
	ethLayer    layers.Ethernet
	ipLayer     layers.IPv4
	tcpLayer    layers.TCP
)

func wireShark(deviceName string, port uint16) {
	filter := getFilter(port)
	handle, err := pcap.OpenLive(deviceName, snapshotLen, promiscuous, timeout)
	if err != nil {
		fmt.Printf("pcap open live failed: %v", err)
		return
	}
	if err := handle.SetBPFFilter(filter); err != nil {
		fmt.Printf("set bpf filter failed: %v", err)
		return
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetSource.NoCopy = true

	for packet := range packetSource.Packets() {
		if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
			fmt.Println("unexpected packet")
			continue
		}

		parser := gopacket.NewDecodingLayerParser(
			layers.LayerTypeEthernet,
			&ethLayer,
			&ipLayer,
			&tcpLayer,
		)

		var foundLayerTypes []gopacket.LayerType
		err := parser.DecodeLayers(packet.Data(), &foundLayerTypes)

		if err != nil {
			fmt.Println("Trouble decoding layers: ", err)
		}

		for _, layerType := range foundLayerTypes {
			var srcIP, srcPort, dstIP, dstPort string
			if layerType == layers.LayerTypeIPv4 {
				srcIP = ipLayer.SrcIP.String()
				dstIP = ipLayer.DstIP.String()
			}
			if layerType == layers.LayerTypeTCP {
				srcPort = tcpLayer.SrcPort.String()
				dstPort = tcpLayer.DstPort.String()
				//fmt.Println("TCP SYN:", tcpLayer.SYN, " | ACK:", tcpLayer.ACK)
			}

			fmt.Printf("%s:%s -> %s:%s\n", srcIP, srcPort, dstIP, dstPort)
		}

	}
}

//定义过滤器
func getFilter(port uint16) string {
	filter := fmt.Sprintf("tcp and ((src port %v) or (dst port %v))", port, port)
	return filter
}
