package fanta_bigliettnonocelho_carichero

import (
	"fmt"
	colly "github.com/gocolly/colly/v2"
)

func scrapeRolse() {
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("div.col_full giocatore", func(e *colly.HTMLElement) {
		e.
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://www.fantacalciopedia.com/lista-calciatori-serie-a/Difensori/")
}
