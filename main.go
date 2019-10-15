package main

import (
	"github.com/gin-gonic/gin"
	"go-object-server/commons"
	"go-object-server/objects"
	"log"
	"net/http"
)

var Config *commons.ConfigModel

//go对象存储服务器入口
func main() {
	log.Println(">>>>>>>>go-object-server start!>>>>>>>>")

	initConfig()

	Config.GetMqPath()

	//todo：数据服务心跳发送消息
	//todo：数据服务locate 定位对象和监听定位消息

	//配置路由监听
	r := gin.Default()
	//初始化路由配置
	initRoutes(r)

	err := r.Run(":" + Config.Server.Port)
	if err != nil {
		log.Fatalln(">>>>>>>>go-object-server error!>>>>>>>>>>")
	}

	log.Println(">>>>>>>>go-object-server end!>>>>>>>>>>")
}

func initConfig() {
	Config = commons.NewConfig()
}

func initRoutes(r *gin.Engine) {

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/ping")
	})
	//测试接口
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	//对象存储url path
	objectsPath := r.Group("/objects")
	{
		//增加存储对象
		objectsPath.POST("/:name", objects.AddObject)
		//删除存储对象
		objectsPath.DELETE("/:name", objects.DeleteObject)
		//修改存储对象
		objectsPath.PUT("/:name", objects.UpdateObject)
		//通过名称获取存储对象
		objectsPath.GET("/:name", objects.GetObjectByName)
	}
}
