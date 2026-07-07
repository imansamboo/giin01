package main

import (
	"flag"
	"fmt"
	"io/ioutil"
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

	//http://127.0.0.1:8080?iman=somename
	server.GET("/", getQueryParams)
	//http://127.0.0.1:8080/age/23/desc
	server.GET("/age/:age/:sort", getURLParams)
	// curl -X POST http://127.0.0.1:8080 -H 'accept: application/json' --raw-data '{"iman":"amir"}'
	server.POST("/", getRequestBody)
	server.Run(":8080")
}

func getQueryParams(c *gin.Context) {
	fmt.Println(c.Query("iman"))
	c.JSON(200, gin.H{"message": "OK"})
}

func getRequestBody(c *gin.Context) {
	reqBody := c.Request.Body
	body, _ := ioutil.ReadAll(reqBody)
	fmt.Println(string(body))
	c.JSON(200, gin.H{"iman": "req body ok"})

}

func getURLParams(c *gin.Context) {
	fmt.Println(c.Param("age"))
	fmt.Println(c.Param("sort"))
	c.JSON(200, gin.H{"message": "OK"})
}

func run() error {
	debug.Logf("pid=%d", os.Getpid())
	return nil
}
