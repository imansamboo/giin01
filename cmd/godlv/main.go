package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

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
	auth := gin.BasicAuth(gin.Accounts{
		"iman": "iman2iman2",
		"amir": "amir2amir2",
	})

	//http://127.0.0.1:8080?iman=somename
	admin := server.Group("/admin", auth)
	{
		admin.GET("/perms", GetPerms)
	}
	server.GET("/", getQueryParams)
	//http://127.0.0.1:8080/age/23/desc
	server.GET("/age/:age/:sort", getURLParams)
	// curl -X POST http://127.0.0.1:8080 -H 'accept: application/json' --raw-data '{"iman":"amir"}'
	server.POST("/", getRequestBody)
	customServer := &http.Server{
		Addr:         ":9090",
		Handler:      server,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	customServer.ListenAndServe()
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

func GetPerms(c *gin.Context) {
	c.JSON(http.StatusAccepted, gin.H{
		"perms": []string{"admin", "viewer"},
	})
}

func run() error {
	debug.Logf("pid=%d", os.Getpid())
	return nil
}
