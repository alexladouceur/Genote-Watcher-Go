package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http/cookiejar"
	"os"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"github.com/joho/godotenv"
)

func getEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func login(c *colly.Collector, fieldsData map[string]string) {
	loginUrl := "https://cas.usherbrooke.ca/login?service=https://www.usherbrooke.ca/genote/public/index.php"

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Error:", err)
		fmt.Println(r.Request.Headers)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(string(r.Body))

		err := c.Visit("https://www.usherbrooke.ca/genote/public/index.php")
		if err != nil {
			log.Fatal(err)
		} else {
			fmt.Println("Visited the page")
		}
	})

	err := c.Post(loginUrl, fieldsData)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Logged in")
	}

}

func main() {

	var userAgents = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:130.0) Gecko/20100101 Firefox/130.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36 Edg/128.0.0.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36 OPR/113.0.0.0",
	}
	// Create a new collector
	c := colly.NewCollector(
		// Attach a debugger to the collector
		colly.Debugger(&debug.LogDebugger{}),
		colly.UserAgent(userAgents[rand.Intn(len(userAgents))]),
	)
	jar, _ := cookiejar.New(nil)
	c.SetCookieJar(jar)

	// The URL for the login action
	loginPageURL := "https://cas.usherbrooke.ca/login?service=https%3A%2F%2Fwww.usherbrooke.ca%2Fgenote%2Fpublic%2Findex.php"

	fieldsData := map[string]string{
		"username": getEnvVariable("GENOTE_USER"),
		"password": getEnvVariable("GENOTE_PASSWORD"),
		"submit":   "",
	}

	// After logging in, visit the page you want to scrape
	c.OnHTML("input[type='hidden']", func(e *colly.HTMLElement) {
		// Scrape information
		// fmt.Printf("%s: %s \n", e.Attr("name"), e.Attr("value"))
		fieldsData[e.Attr("name")] = e.Attr("value")
	})

	c.OnScraped(func(r *colly.Response) {
		login(c, fieldsData)
	})

	// Visit a page that requires authentication
	c.Visit(loginPageURL)
}
