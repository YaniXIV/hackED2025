package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/http/cookiejar"

	"context"
	"github.com/chromedp/chromedp"
	//"github.com/PuerkitoBio/goquery"
	//"log"
	//"net/http"
	//"net/http/cookiejar"
	"strings"
	"time"
)

// I think this is our best bet for the scraping.
// my current issue is the inline javascript that is embessed into the html. It may be best to just keep it
// and tell the llm to ignore any script.
func GetHtml() {
	var link string
	fmt.Println("getting data.")
	//link = "https://httpbin.org/#/HTTP_Methods/delete_delete"
	//link = "https://www.ign.com/articles/where-to-buy-midnight-black-ps5-playstation-portal-pulse-elite-dualsense-edge"
	//link = "https://scrapeops.io/web-scraping-playbook/403-forbidden-error-web-scraping/"
	link = "https://www.forbes.com/sites/paultassi/2025/02/13/the-gta-6-release-date-window-narrowed-by-the-borderlands-4-release-date/"
	//enabling cookies
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	//creating the request
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		log.Fatal(err)
	}

	//emulating browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Mobile Safari/537.36")
	req.Header.Set("Referer", "https://google.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "close")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		log.Fatalf("HTTP Error: %d\n", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".header").Remove()
	doc.Find("script").Remove()
	doc.Find(".sidebar").Remove()
	doc.Find(".footer").Remove()
	text := doc.Find(".article-body").Text()
	//strings.TrimSpace(text)

	//fmt.Println("The file size is :", len(text))
	fmt.Println("The file content: ", text)

}

// Add the http:// prefix to urls.
func fixUrl(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "http://" + url
	}
	return url
}

// Tried to run a headless chrome browser.
func Scrape() {
	// Create a Chrome instance
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Set timeout to avoid hanging
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var pageContent string

	// Navigate and extract text
	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.forbes.com/sites/paultassi/2025/02/13/the-gta-6-release-date-window-narrowed-by-the-borderlands-4-release-date/"), // Replace with target URL
		chromedp.WaitVisible("body", chromedp.ByQuery),        // Wait for page to load
		chromedp.Text("body", &pageContent, chromedp.ByQuery), // Extract body text
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Clean and print text
	cleanedText := strings.TrimSpace(pageContent)
	fmt.Println(cleanedText)
}
