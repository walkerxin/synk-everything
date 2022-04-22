package controller

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func GetFilePath(name string) string {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exeDir := filepath.Dir(exe)
	return filepath.Join(exeDir, "uploads", name)
}

// "/uploads/xxx"
func DownloadsController(c *gin.Context) {
	name := c.Param("path")
	c.Header("Content-Description", "download file")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+name)
	c.Header("Content-Type", "application/octet-stream")
	c.File(GetFilePath(name))
}
