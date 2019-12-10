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
	"github.com/gomodule/redigo/redis"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func main() {
	wireShark("eth0", uint16(12345))
}
