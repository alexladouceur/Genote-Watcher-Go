package main

import (
	"fmt"
	"log"
	"net/http/cookiejar"

	"genote-watcher/utils"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
)

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
		}
	})

	err := c.Post(loginUrl, fieldsData)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {

	c := colly.NewCollector(
		colly.Debugger(&debug.LogDebugger{}),
		colly.UserAgent(utils.GetRandomUserAgent()),
	)

	jar, _ := cookiejar.New(nil)
	c.SetCookieJar(jar)

	loginPageURL := "https://cas.usherbrooke.ca/login?service=https%3A%2F%2Fwww.usherbrooke.ca%2Fgenote%2Fpublic%2Findex.php"

	fieldsData := map[string]string{
		"username": utils.GetEnvVariable("GENOTE_USER"),
		"password": utils.GetEnvVariable("GENOTE_PASSWORD"),
		"submit":   "",
	}

	c.OnHTML("input[type='hidden']", func(e *colly.HTMLElement) {
		fieldsData[e.Attr("name")] = e.Attr("value")
	})

	c.OnScraped(func(r *colly.Response) {
		login(c, fieldsData)
	})

	c.Visit(loginPageURL)
}
