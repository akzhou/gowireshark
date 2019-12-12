/**
 * @Author: Administrator
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2019/12/10 9:35
 */

package main

import (
	"github.com/gin-gonic/gin"
	"gowireshark/pkg"
)

func main() {
	go pkg.WireShark("eth0", 80)

	router := gin.Default()
	//TODO:获取下载进度
	router.GET("/getDownloading", func(c *gin.Context) {
		udid := c.Query("udid")
		timestamp := c.Query("timestamp")
		if udid == "" || timestamp == "" {
			c.JSON(200, gin.H{
				"code":    -1,
				"message": "udid及timestamp不为空！",
				"data":    nil,
			})
			return
		}
		downloading := pkg.GetDownloading(udid, timestamp)
		c.JSON(200, gin.H{
			"code":    0,
			"message": "success",
			"data":    downloading,
		})
	})

	router.Run(":8081")
}
