package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	//"github.com/PuerkitoBio/goquery"
)


func GetHtml(){
  var link string

  fmt.Println("getting data.")

  link = "http://www.ign.com/articles/deadpool-kills-the-marvel-universe-one-last-time-caps-off-marvels-bloodiest-trilogy"
  //_, err := fmt.Scanln(&link)

    

  res, err := http.Get(fixUrl(link)) 
  if err != nil{
    log.Fatal(err)
  }

  defer res.Body.Close()


  if res.StatusCode != http.StatusOK {
    log.Fatalf("HTTP Error: %d\n", res.StatusCode)
  }
  

  content, err := io.ReadAll(res.Body)

  if err != nil{
    log.Fatal(err)
  }

  fmt.Println("The file size is :", len(content))
  



}


func fixUrl(url string)string{
  if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://"){
    return "http://" + url
  }
  return url
}

