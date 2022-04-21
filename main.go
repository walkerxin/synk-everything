package main

import (
	"embed"
	"io/fs"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"github.com/zserge/lorca"
)

//go:embed frontend/dist/*
var FS embed.FS

func main() {
	port := "27149"
	go func() {
		gin.SetMode(gin.DebugMode)
		router := gin.Default()
		staticFiles, _ := fs.Sub(FS, "frontend/dist")
		router.POST("api/v1/files", FilesController)
		router.GET("/api/v1/qrcodes", QrcodesController)
		router.GET("/api/v1/addresses", AddressesController)
		router.POST("/api/v1/texts", TextsController)
		// 静态路由
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
	}()

	ui, _ := lorca.New("http://localhost:"+port+"/static", "", 1000, 600)

	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-chSignal:
	case <-ui.Done():
	}
	ui.Close()
}

func FilesController(c *gin.Context) {
	file, err := c.FormFile("raw") // 获取单个表单文件
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	exe, err := os.Executable() // 创建目录
	if err != nil {
		log.Fatal(err)
	}
	exeDir := filepath.Dir(exe)
	uploadsDir := filepath.Join(exeDir, "uploads")
	err = os.MkdirAll(uploadsDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	filename := uuid.New().String() // 存
	filename += filepath.Ext(file.Filename)
	fileErr := c.SaveUploadedFile(file, filepath.Join(uploadsDir, filename))
	if fileErr != nil {
		log.Fatal(fileErr)
	}
	c.JSON(http.StatusOK, gin.H{"url": "/" + path.Join("uploads", filename)})
}

// "qrcodes?content=" 获取查询字符串
func QrcodesController(c *gin.Context) {
	if content := c.Query("content"); content != "" {
		png, err := qrcode.Encode(content, qrcode.Medium, 256)
		if err != nil {
			log.Fatal(err)
		}
		c.Data(http.StatusOK, "image/png", png) // 为了在前端展示（不设置未二进制流）
	} else {
		c.Status(http.StatusBadRequest)
	}
}

func AddressesController(c *gin.Context) {
	addrs, _ := net.InterfaceAddrs()
	var result []string
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				result = append(result, ipnet.IP.To4().String())
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"addresses": result})
}

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

		filename := uuid.New().String() //4. 生成文件名
		fullpath := path.Join("uploads", filename+".txt")
		err = ioutil.WriteFile(filepath.Join(exeDir, fullpath), []byte(json.Raw), 0644) //5. 写入文件
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"url": "/" + fullpath})
	}
}
