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
	log "github.com/sirupsen/logrus"
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
		log.Info("bindUdidAndFile,udid:%s,fileName:%s", req.Udid, req.FileName)
		pkg.BindUdidAndFile(req.Udid, req.FileName)
		c.JSON(200, gin.H{
			"Code":    0,
			"Message": "Udid与文件关联成功！",
		})
	})

	router.POST("/removeDownloading", func(c *gin.Context) {
		req := struct {
			Udid string `json:"udid"`
		}{}
		err := c.BindJSON(&req)
		if nil != err {
			c.JSON(200, gin.H{
				"Code":    -1,
				"Message": err.Error(),
			})
			return
		}
		if req.Udid == "" {
			c.JSON(200, gin.H{
				"Code":    -1,
				"Message": "udid不能为空！",
			})
			return
		}
		log.Info("removeDownloading,udid:%s", req.Udid)
		pkg.RemoveDownloading(req.Udid)
		c.JSON(200, gin.H{
			"Code":    0,
			"Message": "根据Udid删除下载进度成功！",
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
