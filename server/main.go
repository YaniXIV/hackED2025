package main

import (
	"fmt"
	"net/http"
  "sync"
	"github.com/gin-gonic/gin"
)

type responseData struct{
  Link string `json:"link"`
}

//locally storing the linkdata for now and using mutexes to ensure safe reads and writes.
var links []responseData
var mu sync.Mutex


func main(){
  //initializing the gin engine
  r := gin.Default()
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
  r.Run()
}




func handleLink(c *gin.Context){
    var linkdata responseData 

    if binderr := c.ShouldBindJSON(&linkdata); binderr != nil {
        c.AbortWithError(http.StatusBadRequest, binderr)
        return 
    }

  //Adding our link to slice stored in mem temporarily
    mu.Lock()
    links = append(links,linkdata)
    mu.Unlock()

    c.JSON(http.StatusAccepted, gin.H{
    "message": "Link added successfully!",
    "data": linkdata, 
  })

  checkData()
  }

//small function to check if links are storing correctly. 
func checkData(){
  mu.Lock()
  for i := range(links){
    fmt.Println(links[i])
  }
  mu.Unlock()
}






