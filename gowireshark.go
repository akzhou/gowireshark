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
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gowireshark/pkg"
	"net/http"
	"strings"
)

func main() {
	go pkg.WireShark("eth0")

	router := gin.Default()

	router.Use(cors())

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

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		origin := c.Request.Header.Get("Origin")
		var headerKeys []string
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			//下面的都是乱添加的-_-~
			// c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", headerStr)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			// c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			// c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}

		c.Next()
	}
}
