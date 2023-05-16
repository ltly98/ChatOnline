package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"math"
	"os"
	"time"
)

/*
	编辑者：F,Z
	功能：使用logrus进行日志管理
	说明：此处使用洋葱模型中间件,此处WithLinkName(软链接)在windows下调试应具有管理员权限
         日志分割方式记录：
         1.运维日志分割
         2.使用go的包，这里推荐file-rotatelogs,附带使用hook的包lfshook
*/

func Logger() gin.HandlerFunc {
	filePath := "log/log"
	linkName := "Latest_log.log"
	src, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("err:", err)
	}
	logger := logrus.New()
	logger.Out = src
	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)
	logWriter, _ := rotatelogs.New(filePath+"%Y%m.log",
		rotatelogs.WithMaxAge(7*24+time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithLinkName(linkName))
	//用于将日志级别映射到 io 的映射
	writeMap := lfshook.WriterMap{
		logrus.InfoLevel:  logWriter,
		logrus.FatalLevel: logWriter,
		logrus.DebugLevel: logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.PanicLevel: logWriter,
		logrus.TraceLevel: logWriter,
	}
	//使用lfshook创建hook，并设置输出格式
	Hook := lfshook.NewHook(writeMap, &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	//添加hook
	logger.AddHook(Hook)

	return func(c *gin.Context) {
		//开始时间
		startTime := time.Now()
		//洋葱模型中间件
		c.Next()
		//终止时间
		stopTime := time.Since(startTime)
		//时间花销.ceil向上取整，除以一百万
		spendTime := fmt.Sprintf("%d ns", int(math.Ceil(float64(stopTime.Nanoseconds()))/1000000.0))
		//主机名，出错设置默认值
		hostName, err := os.Hostname()
		if err != nil {
			hostName = "unknown"
		}
		//状态码
		statusCode := c.Writer.Status()
		//客户端ip
		clientIp := c.ClientIP()
		//请求终端
		userAgent := c.Request.UserAgent()
		//数据大小
		dataSize := c.Writer.Size()
		if dataSize < 0 {
			dataSize = 0
		}
		//请求方式
		method := c.Request.Method
		//请求地址
		path := c.Request.RequestURI

		entry := logger.WithFields(logrus.Fields{
			"HostName":  hostName,
			"Status":    statusCode,
			"SpendTime": spendTime,
			"IP":        clientIp,
			"Method":    method,
			"Path":      path,
			"DataSize":  dataSize,
			"Agent":     userAgent,
		})
		//如果gin内部错误，传入错误
		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		}
		if statusCode >= 500 {
			entry.Error()
		} else if statusCode >= 400 {
			entry.Warn()
		} else {
			entry.Info()
		}

	}
}
