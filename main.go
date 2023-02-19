package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type ProxyRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

type ProxyResponse struct {
	ID      string            `json:"id"`
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Length  int               `json:"length"`
}

var requests = make(map[string]ProxyRequest)
var responses = make(map[string]ProxyResponse)

func main() {
	r := gin.Default()
	r.POST("/", proxyHandler)
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, "hey")
	})
	r.Run()
}

func proxyHandler(c *gin.Context) {
	w := c.Writer
	r := c.Request
	//w http.ResponseWriter, r *http.Request
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var proxyReq ProxyRequest
	err = json.Unmarshal(body, &proxyReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(proxyReq.Method, proxyReq.URL, nil)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for key, value := range proxyReq.Headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	proxyRes := ProxyResponse{
		ID:      "requestId",
		Status:  resp.StatusCode,
		Headers: make(map[string]string),
		Length:  len(body),
	}

	for key, value := range resp.Header {
		proxyRes.Headers[key] = value[0]
	}

	responseBody, err := json.Marshal(proxyRes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)

	requests["requestId"] = proxyReq
	responses["requestId"] = proxyRes
	fmt.Printf("Request: %+v\n", proxyReq)
	fmt.Printf("Response: %+v\n", proxyRes)
}
