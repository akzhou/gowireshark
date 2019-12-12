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
)

var (
	snapshotLen = int32(65535)
	promiscuous = false
	timeout     = pcap.BlockForever
)

func WireShark(deviceName string, port uint16) {
	filter := getFilter(port)
	handle, err := pcap.OpenLive(deviceName, snapshotLen, promiscuous, timeout)
	if err != nil {
		log.Error(err)
		return
	}
	if err := handle.SetBPFFilter(filter); err != nil {
		log.Error(err)
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

		var srcIP, srcPort, dstIP, dstPort string

		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer != nil {
			ip, _ := ipLayer.(*layers.IPv4)
			srcIP = ip.SrcIP.String()
			dstIP = ip.DstIP.String()
		}

		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer != nil {
			tcp, _ := tcpLayer.(*layers.TCP)
			srcPort = tcp.SrcPort.String()
			dstPort = tcp.DstPort.String()
		}

		applicationLayer := packet.ApplicationLayer()
		if applicationLayer == nil {
			continue
		}
		//log.Infof("%s:%s  ->  %s:%s", srcIP, srcPort, dstIP, dstPort)
		//出口流量
		if strings.Contains(srcPort, strconv.Itoa(int(port))) {
			key := fmt.Sprintf("%s_%s", dstIP, dstPort)
			IncrBy(key, len(applicationLayer.Payload()))
			continue
		}

		//入口流量统计

		log.Infof("InPayload:%s", applicationLayer.Payload())
	}
}

//定义过滤器
func getFilter(port uint16) string {
	filter := fmt.Sprintf("tcp and ((src port %v) or (dst port %v))", port, port)
	return filter
}
