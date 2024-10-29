package main

import (
	"fmt"
	"log"
	"net/http/cookiejar"

	"genote-watcher/parsers"
	"genote-watcher/utils"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
)

const (
	LOGIN_URL = "https://cas.usherbrooke.ca/login?service=https://www.usherbrooke.ca/genote/public/index.php"
)

func createCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.Debugger(&debug.LogDebugger{}),
		colly.UserAgent(utils.GetRandomUserAgent()),
	)

	jar, _ := cookiejar.New(nil)
	c.SetCookieJar(jar)

	return c
}

func getLoginFields(c *colly.Collector) map[string]string {

	defer c.Visit(LOGIN_URL)

	fieldsData := map[string]string{
		"username": utils.GetEnvVariable("GENOTE_USER"),
		"password": utils.GetEnvVariable("GENOTE_PASSWORD"),
		"submit":   "",
	}

	c.OnHTML("input[type='hidden']", func(e *colly.HTMLElement) {
		fieldsData[e.Attr("name")] = e.Attr("value")
	})

	return fieldsData
}

func login(c *colly.Collector) {
	fieldsData := getLoginFields(c)

	err := c.Post(LOGIN_URL, fieldsData)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	c := createCollector()
	login(c)

	rows := parsers.ParseClasses(c.Clone())

	fmt.Println(rows)
}
