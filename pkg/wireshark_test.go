/**
 * @Author: Administrator
 * @Description:
 * @File:  wireshark_test
 * @Version: 1.0.0
 * @Date: 2019/12/12 10:08
 */

package pkg

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
	"testing"
)

func TestWireShark(t *testing.T) {
	s := "/storage/5de0bde5bcd9c.ipa?udid=233d3433a7475bda14bd9aa6a9491c31186db935&timestamp=12222222"
	//解析这个 URL 并确保解析没有出错。
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(u.Path)

	paths := strings.Split(u.Path, "/")
	fileName := paths[len(paths)-1]
	fmt.Println(fileName)
	m, _ := url.ParseQuery(u.RawQuery)
	fmt.Println(m["udid"][0])
	fmt.Println(m["timestamp"][0])
}

func TestWireShark2(t *testing.T) {
	getFileSize("e:/xinxinserver/config/gowireshark.toml")
}

func TestWireShark3(t *testing.T) {
	var temp sync.Map
	temp.Store("temp", 0)
	fmt.Println(temp.Load("temp"))

	temp.Store("temp", 1)
	fmt.Println(temp.Load("temp"))

}
