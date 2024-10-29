package parsers

import (
	"genote-watcher/model"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

func ParseClasses(c *colly.Collector) []model.CourseRow {
	rows := []model.CourseRow{}
	c.OnHTML("table:nth-child(4) tbody", func(e *colly.HTMLElement) {
		cr := model.CourseRow{}
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			cr.CourseName = el.DOM.Find("td:nth-child(1)").Text()
			splitName := strings.Split(cr.CourseName, " ")
			courseCode := splitName[len(splitName)-2]

			cr.CourseCode = courseCode[1:]
			cr.EvaluationAmount, _ = strconv.Atoi(el.DOM.Find("td:nth-child(5)").Text())
			cr.CourseLink = el.DOM.Find("td:nth-child(6) a").AttrOr("href", "")

			rows = append(rows, cr)
		})
	})

	err := c.Visit("https://www.usherbrooke.ca/genote/application/etudiant/cours.php")
	if err != nil {
		log.Fatal(err)
	}
	c.Wait()

	return rows
}
