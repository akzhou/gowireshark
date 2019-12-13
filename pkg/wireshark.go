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
	udidTimestampIPPortMap sync.Map
	udidTimestampFileMap   sync.Map
	iPPortFileMap          sync.Map
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
				log.Info(inputPayloadStr)
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
				m, err := url.ParseQuery(u.RawQuery)
				if nil != err {
					log.Error(err)
					continue
				}
				if 0 == len(m["udid"]) || 0 == len(m["timestamp"]) {
					log.Error(fmt.Errorf("udid and timestamp not nil"))
					continue
				}
				//udidTimestampIPPortMap.Store(m["udid"][0]+"_"+m["timestamp"][0], srcIP+"_"+srcPort)
				//udidTimestampFileMap.Store(m["udid"][0]+"_"+m["timestamp"][0], getFileSize(fileName))
				udidTimestampIPPortMap.Store(m["udid"][0], srcIP+"_"+srcPort)
				udidTimestampFileMap.Store(m["udid"][0], getFileSize(fileName))
				iPPortFileMap.Store(srcIP+"_"+srcPort, int64(0))
				log.Infof("iPPortFileMapInit(key:%v,value:%v)", srcIP+"_"+srcPort, 0)
			}
		}

		//出口流量
		key := dstIP + "_" + dstPort
		log.Infof("%v --->  %v", srcIP+"_"+srcPort, dstIP+"_"+dstPort)
		if v, ok := iPPortFileMap.Load(key); ok {
			if vv, ok := v.(int64); ok {
				iPPortFileMap.Store(key, vv+int64(len(applicationLayer.Payload())))
				log.Infof("iPPortFileMap(key:%v,value:%v)", key, vv+int64(len(applicationLayer.Payload())))
			}
		} else {
			iPPortFileMap.Store(key, int64(len(applicationLayer.Payload())))
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
func GetDownloading(udid, timestamp string) int {
	var fileSize, downloadSize int64
	if v, ok := udidTimestampFileMap.Load(udid); ok {
		//if v, ok := udidTimestampFileMap.Load(udid + "_" + timestamp); ok {
		if vv, ok := v.(int64); ok {
			fileSize = vv
		}
		//log.Infof("GetDownloading.udidTimestampFileMap(key:%v,value:%v)", udid+"_"+timestamp, v)
		log.Infof("GetDownloading.udidTimestampFileMap(key:%v,value:%v)", udid, v)
	}
	//if v, ok := udidTimestampIPPortMap.Load(udid + "_" + timestamp); ok {
	if v, ok := udidTimestampIPPortMap.Load(udid); ok {
		//log.Infof("GetDownloading.udidTimestampIPPortMap(key:%v,value:%v)", udid+"_"+timestamp, v)
		log.Infof("GetDownloading.udidTimestampIPPortMap(key:%v,value:%v)", udid, v)
		if vv, ok := v.(string); ok { //vv表示ip_port
			if vvv, ok := iPPortFileMap.Load(vv); ok { //vvv下载量v
				//log.Infof("GetDownloading.iPPortFileMap(key:%v,value:%v)", udid+"_"+timestamp, vvv)
				log.Infof("GetDownloading.iPPortFileMap(key:%v,value:%v)", udid, vvv)
				if vvvv, ok := vvv.(int64); ok {
					downloadSize = vvvv
				}
			}
		}
	}
	log.Infof("download size:%v,file size:%v", downloadSize, fileSize)
	if fileSize == 0 {
		return 0
	}
	return int(math.Min(float64(downloadSize)/float64(fileSize)*100, 100))
}
