package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"godlv/controller"
	"godlv/debug"
	"godlv/service"

	"github.com/gin-gonic/gin"
)

var (
	videoController controller.VideoController
	videoService    service.VideoService
)

func main() {
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

	server := gin.Default()
	server.GET("/videos", func(c *gin.Context) {
		videos := videoController.FindAll(c)
		c.JSON(http.StatusOK, videos)
	})
	server.POST("/videos", func(c *gin.Context) {
		err := videoController.Save(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"message": "Video created successfully"})
	})
	server.Run(":8080")
}

func run() error {
	debug.Logf("pid=%d", os.Getpid())
	return nil
}
