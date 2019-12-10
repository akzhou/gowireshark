/**
 * @Author: Administrator
 * @Description:
 * @File:  wireshark
 * @Version: 1.0.0
 * @Date: 2019/12/10 19:47
 */

package pkg

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"sync"
)

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
	portTraffic sync.Map
)

func WireShark(deviceName string, port uint16) {
	filter := getFilter(port)
	handle, err := pcap.OpenLive(deviceName, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Error("pcap open live failed: %v", err)
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
		var srcIP, srcPort, dstIP, dstPort string
		for _, layerType := range foundLayerTypes {
			if layerType == layers.LayerTypeIPv4 {
				srcIP = ipLayer.SrcIP.String()
				dstIP = ipLayer.DstIP.String()
			}
			if layerType == layers.LayerTypeTCP {
				srcPort = tcpLayer.SrcPort.String()
				dstPort = tcpLayer.DstPort.String()
			}
		}
		log.Infof("%s:%s  ->  %s:%s", srcIP, srcPort, dstIP, dstPort)
		if !strings.Contains(srcPort, strconv.Itoa(int(port))) {
			continue
		}
		if v, ok := portTraffic.Load(fmt.Sprintf("%s:%s", dstIP, dstPort)); ok {
			fmt.Println(v)
		}
	}
}

//定义过滤器
func getFilter(port uint16) string {
	filter := fmt.Sprintf("tcp and ((src port %v) or (dst port %v))", port, port)
	return filter
}
