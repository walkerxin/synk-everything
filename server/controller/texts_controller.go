package controller

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TextsController(c *gin.Context) {
	var json struct { //1. 从body中获取用户上传的文本
		Raw string
	}
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		exe, err := os.Executable() //2. 获取可执行文件所在的目录
		if err != nil {
			log.Fatal(err)
		}
		exeDir := filepath.Dir(exe)
		uploadsDir := filepath.Join(exeDir, "uploads") //3. 上一步的目录拼接上 uploads，创建目录
		err = os.MkdirAll(uploadsDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("dir=" + uploadsDir)

		filename := uuid.New().String() + ".txt"                                            //4. 生成文件名
		err = ioutil.WriteFile(filepath.Join(uploadsDir, filename), []byte(json.Raw), 0644) //5. 写入文件
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"url": path.Join("/uploads", filename)})
	}
}
