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
	"math"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var (
	snapshotLen = int32(65536)
	promiscuous = true
	timeout     = pcap.BlockForever
)

var (
	udidAndFileMap   sync.Map
	fileAndIPPortMap sync.Map
	ipPortTrafficMap sync.Map
)

func WireShark(deviceName string) {
	filter := getFilter(wireSharkCfg.FileServerPort)
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

		//入口流量
		if !strings.Contains(srcPort, strconv.Itoa(int(wireSharkCfg.FileServerPort))) {
			inputPayloadStr := string(applicationLayer.Payload())
			if strings.Contains(inputPayloadStr, wireSharkCfg.UrlFlag) {
				requests := strings.Split(inputPayloadStr, " ")
				if len(requests) < 2 {
					continue
				}
				u, err := url.Parse(requests[1])
				if nil != err {
					log.Error(err)
					continue
				}
				paths := strings.Split(u.Path, "/")
				fileName := paths[len(paths)-1]
				if "" == fileName {
					log.Errorf("未获取到文件名")
					continue
				}
				fileAndIPPortMap.Store(fileName, srcIP+"_"+srcPort)
				ipPortTrafficMap.Store(srcIP+"_"+srcPort, int64(0))
			}
		}

		//出口流量
		//log.Infof("%v --->  %v", srcIP+"_"+srcPort, dstIP+"_"+dstPort)
		key := dstIP + "_" + dstPort
		if v, ok := ipPortTrafficMap.Load(key); ok {
			if vv, ok := v.(int64); ok {
				ipPortTrafficMap.Store(key, vv+int64(len(applicationLayer.Payload())))
				log.Infof("iPPortFileMap(key:%v,value:%v)", key, vv+int64(len(applicationLayer.Payload())))
			}
		} else {
			ipPortTrafficMap.Store(key, int64(len(applicationLayer.Payload())))
			log.Infof("iPPortFileMap(key:%v,value:%v)", key, len(applicationLayer.Payload()))
		}
	}
}

//定义过滤器
func getFilter(port uint16) string {
	filter := fmt.Sprintf("tcp and ((src port %v) or (dst port %v))", port, port)
	return filter
}

func getFileSize(fileName string) int64 {
	fileName = wireSharkCfg.UrlPath + fileName
	var result int64
	filepath.Walk(fileName, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	log.Infof("getFileSize--->%s ：%d", fileName, result)
	return result
}

//TODO:获取下载进度
func GetDownloading(udid string) int {
	var fileSize, downloadSize int64

	//step1:根据udid获取文件名
	iFileName, ok := udidAndFileMap.Load(udid)
	if !ok {
		log.Warningf("未获取到Udid(%s)对应的文件名称", udid)
		return 0
	}
	fileName, ok := iFileName.(string)
	if !ok {
		log.Warningf("Udid(%s)对应的文件(%v)类型断言失败", udid, iFileName)
		return 0
	}
	fileSize = getFileSize(fileName)

	//step2:根据文件名获取ip:port
	iIPPort, ok := fileAndIPPortMap.Load(fileName)
	if !ok {
		log.Warningf("未获取到Udid(%s)->文件(%s)的IP:Port", udid, fileName)
		return 0
	}
	ipPort, ok := iIPPort.(string)
	if !ok {
		log.Warningf("Udid(%s)->文件(%v)所对应的IPPort(%v)类型断言失败", udid, fileName, iIPPort)
		return 0
	}

	//step3:根据ip:port获取流量
	iTraffic, ok := ipPortTrafficMap.Load(ipPort)
	if !ok {
		log.Warningf("未获取到Udid(%s)->文件(%s)->IP:Port(%s)对应的下载流量", udid, fileName, ipPort)
		return 0
	}
	traffic, ok := iTraffic.(int64)
	if !ok {
		log.Warningf("Udid(%s)->文件(%s)->IPPort(%s)类型所对应的流量(%v)类型断言失败", udid, fileName, ipPort, iTraffic)
		return 0
	}

	//step4:流量统计
	downloadSize = traffic
	log.Infof("download size:%v,file size:%v", downloadSize, fileSize)
	if fileSize == 0 {
		return 0
	}
	return int(math.Min(float64(downloadSize)/float64(fileSize)*100, 100))
}

func BindUdidAndFile(udid, file string) {
	udidAndFileMap.Store(udid, file)
}
