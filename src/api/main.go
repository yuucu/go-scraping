package main

import (
  "fmt"

  "github.com/gin-gonic/gin"
  "github.com/PuerkitoBio/goquery"

)

func main() {

  r := gin.Default()

  r.GET("/", func(c *gin.Context) {
    c.JSON(200, gin.H{
      "message": "Hello",
    })
  })

  r.GET("/scraping", func(c *gin.Context) {

    doc, err := goquery.NewDocument("https://github.com/PuerkitoBio/goquery")
    if err != nil {
      fmt.Print("url scarapping failed")
    }

    doc.Find("a").Each(func(_ int, s *goquery.Selection) {
      url, _ := s.Attr("href")
      fmt.Println(url)
    })

    c.JSON(200, gin.H{
      "message": "switch-scraping",
    })
  })

  // ポートを設定しています。
  r.Run(":3001")
}
