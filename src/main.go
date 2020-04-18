
package main

import (
  "github.com/PuerkitoBio/goquery"
  "github.com/joho/godotenv"
  "os"
  "fmt"
  "time"
  "strings"
  "strconv"
  "net/http"
  "net/url"
)

func main() {
  err := godotenv.Load("../.env")
  if err != nil {
  }

  accessToken := os.Getenv("TOKEN")

  ticker(5, accessToken)
}

func extract_price(text string) []string {
  return strings.Split(strings.TrimSpace(text), " ")
}

func scraping(accessToken string) {
  url := "https://www.amazon.co.jp/gp/offer-listing/B07X1QJRJJ/ref=olp_twister_all?ie=UTF8&mv_color_name=all&mv_pattern_name=all"
  doc, err := goquery.NewDocument(url)

  if err != nil {
    fmt.Print("url scarapping failed")
  }


  item_names := []string{}
  item_descriptions := []string{}
  item_prices := []int{}
  doc.Find("#olpOfferList > div > div").Each(func(_ int, s *goquery.Selection) {
    s.Find("#variationRow > ul >  li:nth-child(2) > span > span").Each(func(_ int, s *goquery.Selection) {
      item_name := strings.TrimSpace(s.Text())
      item_names = append(item_names, item_name)
    })
    s.Find("#variationRow > ul >  li:nth-child(3) > span > span").Each(func(_ int, s *goquery.Selection) {
      item_description := strings.TrimSpace(s.Text())
      item_descriptions = append(item_descriptions, item_description)
    })
    s.Find(".olpOffer > .olpPriceColumn > .olpOfferPrice").Each(func(_ int, s *goquery.Selection) {
      price_str := extract_price(s.Text())[1]
      price_str = strings.Replace(price_str, ",", "", -1)
      price_int, _ := strconv.Atoi(price_str)
      item_prices = append(item_prices, price_int)
    })
  })

  if len(item_names) == 0 {
    fmt.Print("Amazon error page...")
    return
  }

  for i := 0; i < len(item_names); i++ {

    output := "* " + item_names[i] + " / " + item_descriptions[i] + "\n  " + strconv.Itoa(item_prices[i])
    fmt.Printf(output)
    split := strings.Split(item_names[i], " ")
    if len(split) > 1 {
      if split[1] == "Lite" {
        if item_prices[i] < 30000 {
          fmt.Printf(" ! ")
          fmt.Printf("買い時!!( T_T)＼(^-^ )")
          notify(item_names[i] + ": " + strconv.Itoa(item_prices[i]) + " / url: " + url, accessToken)
        }
      }
    } else {
      if item_prices[i] < 40000 {
        fmt.Printf(" ! ")
        fmt.Printf("買い時!!( T_T)＼(^-^ )")
        notify(item_names[i] + ": " + strconv.Itoa(item_prices[i]) + " / url: " + url, accessToken)
      }
    }
    fmt.Println("")
  }
}

func notify(message string, accessToken string) {

  URL := "https://notify-api.line.me/api/notify"

  u, _ := url.ParseRequestURI(URL)
  c := &http.Client{}
  form := url.Values{}
  form.Add("message", message)
  body := strings.NewReader(form.Encode())
  req, _ := http.NewRequest("POST", u.String(), body)
  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
  req.Header.Set("Authorization", "Bearer "+accessToken)

  c.Do(req)
}


func ticker(sec int, accessToken string) {
  t := time.NewTicker(time.Duration(sec) * time.Second)
  defer t.Stop()

  for {
    select {
    case now := <-t.C:

      os.Stdout.Write([]byte{0x1B, 0x5B, 0x33, 0x3B, 0x4A, 0x1B, 0x5B, 0x48, 0x1B, 0x5B, 0x32, 0x4A})
      fmt.Println("", now.Format(time.RFC3339))
      fmt.Println("")
      scraping(accessToken)

    }
  }
}
