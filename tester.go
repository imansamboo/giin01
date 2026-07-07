package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	url := "http://127.0.0.1:8080/videos"
	header := "Basic " + base64.StdEncoding.EncodeToString([]byte("iman:123456"))
	method := "POST"
	reqBody := `{"title": "iman test", "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ", "description": "Test Description"}`
	req, err := http.NewRequest(method, url, strings.NewReader(reqBody))
	req.Header.Add("Authorization", header)
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer response.Body.Close()
	responseBody, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(responseBody))
	fmt.Println(response.StatusCode)
	if response.StatusCode >= 400 {
		os.Exit(0)
	}
	req.Header.Add("Authorization", header)
	req.Header.Add("Content-Type", "application/json")

	req, err = http.NewRequest(method, url, nil)
	url = "http://127.0.0.1:8080/videos"
	header = "Basic " + base64.StdEncoding.EncodeToString([]byte("iman:123456"))
	method = "GET"
	req, err = http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	req.Header.Add("Authorization", header)
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
	fmt.Println(res.StatusCode)

}
