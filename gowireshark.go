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

func main()  {
	wireShark("eth0",uint16(12345))
}

func wireShark(deviceName string,port uint16)  {
	filter := getFilter(port)
	handle, err :=pcap.OpenLive(deviceName, int32(65535), true, pcap.BlockForever)
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

		//tcpLayer := packet.Layer(layers.LayerTypeTCP)
		//if tcpLayer != nil {
		//	tcp, _ := tcpLayer.(*layers.TCP)
		//	// TCP layer variables:
		//	// SrcPort, DstPort, Seq, Ack, DataOffset, Window, Checksum, Urgent
		//	// Bool flags: FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS
		//	fmt.Printf("From ip:port %d to %d\n", tcp.SrcPort, tcp.DstPort)
		//	fmt.Println("Sequence number: ", tcp.Seq)
		//	fmt.Println()
		//}
		//
		//applicationLayer := packet.ApplicationLayer()
		//if  applicationLayer!= nil{
		//	fmt.Printf("applicationLayer:%v\n",applicationLayer)
		//}
		fmt.Printf("packet:%v\n",packet)

		// tcp 层
		tcp := packet.TransportLayer().(*layers.TCP)
		fmt.Printf("tcp:%v\n", tcp)
		// tcp payload，也即是tcp传输的数据
		fmt.Printf("tcp payload:%v\n", tcp.Payload)
	}
}

//定义过滤器
func getFilter(port uint16) string {
	filter := fmt.Sprintf("tcp and ((src port %v) or (dst port %v))",  port, port)
	return filter
}