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
	go pkg.WireShark("eth0")

	router := gin.Default()

	router.POST("/bindUdidAndFile", func(c *gin.Context) {
		req := struct {
			Udid     string `json:"udid"`
			FileName string `json:"fileName"`
		}{}
		err := c.BindJSON(&req)
		if nil != err {
			c.JSON(200, gin.H{
				"Code":    -1,
				"Message": err.Error(),
			})
			return
		}
		if req.Udid == "" || req.FileName == "" {
			c.JSON(200, gin.H{
				"Code":    -1,
				"Message": "udid或fileName不能为空！",
			})
			return
		}
		pkg.BindUdidAndFile(req.Udid, req.FileName)
		c.JSON(200, gin.H{
			"Code":    0,
			"Message": "Udid与文件关联成功！",
		})
	})

	//TODO:获取下载进度
	router.GET("/getDownloading", func(c *gin.Context) {
		udid := c.Query("udid")
		if udid == "" {
			c.JSON(200, gin.H{
				"code":    -1,
				"message": "udid不为空！",
				"data":    nil,
			})
			return
		}
		downloading := pkg.GetDownloading(udid)
		c.JSON(200, gin.H{
			"code":    0,
			"message": "success",
			"data":    downloading,
		})
	})

	router.Run(":8181")
}
