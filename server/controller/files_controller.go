package controller

import (
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func FilesController(c *gin.Context) {
	file, err := c.FormFile("raw") // 获取单个表单文件
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exeDir := filepath.Dir(exe)
	uploadsDir := filepath.Join(exeDir, "uploads")
	err = os.MkdirAll(uploadsDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	filename := uuid.New().String() + filepath.Ext(file.Filename)
	fileErr := c.SaveUploadedFile(file, filepath.Join(uploadsDir, filename))
	if fileErr != nil {
		log.Fatal(fileErr)
	}
	c.JSON(http.StatusOK, gin.H{"url": path.Join("/uploads", filename)})
}
