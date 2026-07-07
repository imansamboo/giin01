package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"

	"godlv/controller"
	"godlv/debug"
	"godlv/middleware"
	"godlv/service"

	"github.com/gin-gonic/gin"
	gindump "github.com/tpkeeper/gin-dump"
)

var (
	videoController controller.VideoController
	videoService    service.VideoService
	logfile         *os.File
)

func setupLogOutput() {
	f, _ := os.Create("gin.log")
	writer := io.MultiWriter(f, os.Stdout)
	gin.DefaultWriter = writer
	log.SetOutput(writer)
}

func main() {
	setupLogOutput()
	videoService = service.New()
	videoController = controller.New(videoService)
	debugFlag := flag.Bool("debug", false, "enable debug output")
	pprofAddr := flag.String("pprof", "", "pprof listen address (e.g. :6060)")
	flag.Parse()

	if *debugFlag {
		debug.SetEnabled(true)
	}

	if *pprofAddr != "" {
		if _, err := debug.StartPprof(*pprofAddr); err != nil {
			log.Fatalf("failed to start pprof: %v", err)
		}
	}

	debug.Log("application started")
	log.Println("godlv running")

	if err := run(); err != nil {
		log.Fatal(err)
	}

	server := gin.New()
	server.Use(gin.Recovery(), middleware.Logger(), middleware.BasicAuth(), gindump.Dump())
	server.GET("/videos", func(c *gin.Context) {
		videos := videoController.FindAll(c)
		c.JSON(http.StatusOK, videos)
	})
	server.POST("/videos", func(c *gin.Context) {
		log.Printf("[save] request received from %s", c.ClientIP())
		err := videoController.Save(c)
		if err != nil {
			log.Printf("[save] failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		log.Println("[save] video created successfully")
		c.JSON(http.StatusAccepted, gin.H{"message": "Video created successfully"})
	})
	server.Run(":8080")
}

func run() error {
	debug.Logf("pid=%d", os.Getpid())
	return nil
}
