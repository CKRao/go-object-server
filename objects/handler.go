package objects

import (
	"github.com/gin-gonic/gin"
	"go-object-server/util"
	"io"
	"log"
	"net/http"
	"os"
)

const StorageRoot = "F:/objectsstorage"

//通过名称获取存储对象
func GetObjectByName(c *gin.Context) {
	c.File(getFilePath(c))
}

//修改存储对象
func UpdateObject(c *gin.Context) {
	filePath := getFilePath(c)
	//判断文件是否存在
	if b, err := util.PathExists(filePath); !b || err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "objects is not exists",
		})
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "update object error",
		})
		return
	}

	defer file.Close()

	_, err = io.Copy(file, c.Request.Body)

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "update object error",
		})
	}
}

//增加存储对象
func AddObject(c *gin.Context) {
	filePath := getFilePath(c)
	//判断文件是否存在
	if b, err := util.PathExists(filePath); b && err == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "objects is exists",
		})
		return
	}

	//从file中获取file对象
	header, err := c.FormFile("file")

	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "get file error",
		})
	}

	//保存文件
	if err := c.SaveUploadedFile(header, filePath); err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "add object error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

//删除存储对象
func DeleteObject(c *gin.Context) {
	filePath := getFilePath(c)

	err := os.Remove(filePath)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "delete error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "delete success",
	})
}

func getFilePath(c *gin.Context) string {
	return StorageRoot + "/objects/" + c.Param("name")
}
