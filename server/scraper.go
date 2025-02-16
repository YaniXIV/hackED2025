package main 
import(
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/http/cookiejar"

	"context"
	"github.com/chromedp/chromedp"
	"strings"
  "time"

)

/* I think this is our best bet for the scraping.
my current issue is the inline javascript that is embedded into the html. It may be best to just keep it
and tell the llm to ignore any script.*/

func GetHtmlHybrid(url string)(string, string){
  bodyText, titleText :=  getHtmlFast(url)
  bodyText = CleanGoQueryContent(bodyText)

  if bodyText == "" || len(bodyText) < 100{
    fmt.Println("Fallng back to ChromeDp for :", url)
    return getHtmlFallback(url)
  }

  return bodyText, titleText
}

func getHtmlFast(l string) (string, string){
	//var link string
	fmt.Println("getting data.")

	//These are some testing links I was using.
	//link = "https://httpbin.org/#/HTTP_Methods/delete_delete"
	//link = "https://www.ign.com/articles/where-to-buy-midnight-black-ps5-playstation-portal-pulse-elite-dualsense-edge"
	//link = "https://scrapeops.io/web-scraping-playbook/403-forbidden-error-web-scraping/"
	//link = "https://www.forbes.com/sites/paultassi/2025/02/13/the-gta-6-release-date-window-narrowed-by-the-borderlands-4-release-date/"

	//enabling cookies
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}

	//creating the request
	req, err := http.NewRequest("GET", l, nil)
	if err != nil {
		log.Fatal(err)
	}

	//Settings headers so that we avoid 403 error
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 15_0 like Mac OS X) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Mobile Safari/537.36")
	req.Header.Set("Referer", "https://google.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "close")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	//executing our get request with special headers to avoid 403 status code.
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	//Error handling if body fails to close
	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil {
			log.Printf("error closing response body: %v", closeErr)
		}
	}()

	if res.StatusCode != http.StatusOK {
		log.Printf("HTTP Error: %d\n", res.StatusCode)
		return fmt.Sprintf("Error: recieved status code %d", res.StatusCode),""
	}

	//reading html data into go query
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	//removing the unnecessary html tags
	doc.Find(".header").Remove()
	doc.Find("script").Remove()
	doc.Find(".sidebar").Remove()
	doc.Find(".footer").Remove()

  title := doc.Find("title").Text()

	//parsing the different possible html tags that might contain our data
	selectors := []string{".article-body", ".post-content", "", "main"}
	var text string
	for _, sel := range selectors {
		textSelection := doc.Find(sel)
		if textSelection.Length() > 0 {
			text = textSelection.Text()
			break
		}
	}

	//strings.TrimSpace(text)

	//fmt.Println("The file size is :", len(text))
	//fmt.Println("The file content: ", text)
	if text == "" {
		log.Printf("No readable content found for URL: %s", l)
    return "Error: No readable context found", ""
	}
  bodyText := strings.TrimSpace(text)
  titleText := strings.TrimSpace(title)

	return bodyText, titleText

}

func CleanGoQueryContent(text string) string {
    // Remove unwanted notification or pop-up data
    if strings.Contains(text, "Click 'OK'") || strings.Contains(text, "Allow") {
        return ""
    }

    // Remove any other known unwanted content
    text = strings.Replace(text, "data", "", -1)

    // Trim any leading/trailing spaces or newlines
    return strings.TrimSpace(text)
}

// Add the http:// prefix to urls.
func fixUrl(url string) string {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return "http://" + url
	}
	return url
}

// Scrape Tried to run a headless Chrome browser.
// Ignore this function
func getHtmlFallback(url string)(string, string){
  // Time out to prevent potential hanging
  ctx, cancel := chromedp.NewContext(
    context.Background(),
    chromedp.WithLogf(log.Printf),
    )
	defer cancel()

  ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
  defer cancel()

  var bodyText string
  var titleText string

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("body", chromedp.ByQuery), // Ensures page loads
		//chromedp.OuterHTML("html", &htmlContent),

    chromedp.Text("title", &titleText),
    chromedp.Text("body", &bodyText),
	)
	

	if err != nil {
		log.Printf("Chromedp error: %v", err)
		return "Error fetching content", ""
	}
  body := strings.TrimSpace(bodyText)
  title := strings.TrimSpace(titleText)

	return body, title
}
