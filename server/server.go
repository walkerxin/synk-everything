package server

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"

	ctr "walkerxin/synk-everything.git/server/controller"

	"github.com/gin-gonic/gin"
)

//go:embed frontend/dist/*
var FS embed.FS

func Run(port string) {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	staticFiles, _ := fs.Sub(FS, "frontend/dist")
	router.POST("api/v1/files", ctr.FilesController)
	router.GET("/uploads/:path", ctr.DownloadsController)
	router.GET("/api/v1/qrcodes", ctr.QrcodesController)
	router.GET("/api/v1/addresses", ctr.AddressesController)
	router.POST("/api/v1/texts", ctr.TextsController)
	router.StaticFS("/static", http.FS(staticFiles))
	router.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/static") {
			reader, err := staticFiles.Open("index.html") // 打开index文件，得到reader
			if err != nil {
				log.Fatal(err)
			}
			defer reader.Close() // 在当前函数退出前执行
			stat, err := reader.Stat()
			if err != nil {
				log.Fatal(err)
			}
			c.DataFromReader(http.StatusOK, stat.Size(), "text/html;", reader, nil)
		} else {
			c.Status(http.StatusNotFound)
		}
	})
	router.Run(":" + port)
}
