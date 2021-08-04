package main

import (
	"fmt"
	_ "math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CharSet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

var (
	count uint64 = 3932829
	urls = make(map[string]string, 0)
)

// Generate a 7 characters base62 representation for a given number
func Encode(n uint64) string {
	s := make([]byte, 7)
	v := n
	for i := 0; i < 7; i++ {
		r := v % 62
		s[6-i] = CharSet[r]
		v = v / 62
	}
	return string(s[:])
}

type Payload struct {
	Url string `json:url`
}

func getCount() uint64 {
	r := count
	count += 1
	return r
}

func getLongUrl(c *gin.Context) {
	tinyUrl := c.Param("tinyUrl")
	fmt.Println("Getting: ", tinyUrl)
	fmt.Println(urls)
	longUrl, prs := urls[tinyUrl]
	if !prs {
		longUrl = "http://google.com"
	}
	c.Redirect(http.StatusMovedPermanently, longUrl)
}

func createShortUrl(c *gin.Context) {
	var payload Payload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error})
		return
	}
	fmt.Println(payload)
	shortUrl := Encode(getCount())
	urls[shortUrl] = payload.Url
	c.JSON(http.StatusOK, gin.H{"shortUrl": shortUrl})
}

func main() {
	r := gin.Default()
	r.GET("/r/:tinyUrl", getLongUrl)
	r.POST("/tiny", createShortUrl)
	r.Run("localhost:8080")
}