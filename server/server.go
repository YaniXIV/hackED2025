package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
  "github.com/gin-contrib/cors"
)

type responseData struct {
	Link string `json:"link"`
}

// locally storing the linkdata for now and using mutexes to ensure safe reads and writes.
var links []responseData
var mu sync.Mutex

func SetupRouter() {
	//initializing the gin engine
	r := gin.Default()
  r.Use(cors.Default())
	fmt.Println("Server starting~")
	/*
		  post request to handle links
		  the json structure should be in the format
		{
		  "link": "https:/sometihng.com"
		}

	*/
	r.POST("/api/link", handleLink)

	//starts server
  r.Run(":8080")
}

func handleLink(c *gin.Context) {
	var linkdata responseData

	if bindErr := c.ShouldBindJSON(&linkdata); bindErr != nil {
		c.AbortWithError(http.StatusBadRequest, bindErr)
		return
	}

	//Adding our link to slice stored in mem temporarily
	mu.Lock()
	links = append(links, linkdata)
	mu.Unlock()

	text, title := GetHtmlHybrid(linkdata.Link)
	c.JSON(http.StatusAccepted, gin.H{
		"message": "Link added successfully!",
		"data":    text,
    "title": title,
	})

	checkData()
}

// small function to check if links are storing correctly.
func checkData() {
	mu.Lock()
	for i := range links {
		fmt.Println(links[i])
	}
	mu.Unlock()
}
