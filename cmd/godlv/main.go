package main

import (
	"flag"
	"log"
	"os"

	"godlv/debug"

	"github.com/gin-gonic/gin"
)

func main() {
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
	server.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "OK"})
	})
	server.Run(":8080")
}

func run() error {
	debug.Logf("pid=%d", os.Getpid())
	return nil
}
