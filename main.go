package main

import (
	"fmt"
	"log"
	"net/http/cookiejar"
	"time"

	"genote-watcher/scrapers"
	"genote-watcher/utils"

	"github.com/gocolly/colly/v2"
)

const (
	LOGIN_URL = "https://cas.usherbrooke.ca/login?service=https://www.usherbrooke.ca/genote/public/index.php"
)

var config *utils.Config

func createCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.UserAgent(utils.GetRandomUserAgent()),
	)

	jar, _ := cookiejar.New(nil)
	c.SetCookieJar(jar)

	return c
}

func getLoginFields(c *colly.Collector) map[string]string {

	defer c.Visit(LOGIN_URL)

	fieldsData := map[string]string{
		"username": config.Username,
		"password": config.Password,
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
	config = utils.MustGetConfig()

	c := createCollector()
	login(c)

	rows := scrapers.ScrapeCourseRows(c.Clone())

	oldRows := utils.ReadResultFile()
	if oldRows == nil {
		utils.WriteResultFile(rows)
		return
	}

	diffRows := []string{}

	for index := range rows {
		if !rows[index].Equal(&oldRows[index]) {
			diffRows = append(diffRows, rows[index].CourseCode)
		}
	}

	now := time.Now()
	formattedDate := now.Format("2006/01/02 15:04:05")
	var changesDetected bool
	for _, courseCode := range diffRows {
		fmt.Printf("[%s] Nouvelle note en %s est disponible sur Genote!\n", formattedDate, courseCode)
		utils.NotifyUser(config.DiscordWebhook, courseCode)
		changesDetected = true
	}

	if !changesDetected {
		fmt.Printf("[%s] Aucun changement détecté\n", formattedDate)
	}

	utils.WriteResultFile(rows)
}
