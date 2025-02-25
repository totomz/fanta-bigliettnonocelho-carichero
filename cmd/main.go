package main

import (
	"log"
	"strings"

	"github.com/gocolly/colly/v2"
)

// Player represents the scraped data for one player.
type Player struct {
	Title  string // from h3.tit_calc
	Team   string // from p > small (e.g. Napoli)
	ALG    string // from span.punt_calc data-original-title attribute
	Pres   string // numeric value from first span.stats_calc (value near PRES.)
	FMedia string // numeric value from second span.stats_calc (value near F.MEDIA)
	Trend  string
}

func scrapePlayers(url string) ([]Player, error) {
	var players []Player

	// Instantiate default collector.
	c := colly.NewCollector(
		// If you know the domain, restrict visits.
		colly.AllowedDomains("www.fantacalciopedia.com"),
	)

	// Called for each HTML element which matches the selector.
	// We select div elements that have both classes "col_full" and "giocatore".
	c.OnHTML("div.col_full.giocatore", func(e *colly.HTMLElement) {
		player := Player{}

		// Get the player's title (e.g., from <h3 class="tit_calc">)
		player.Title = strings.TrimSpace(e.ChildText("h3.tit_calc"))

		// Get the team name (e.g., from <p><small> Napoli</small></p>)
		player.Team = strings.TrimSpace(e.ChildText("p small"))

		// Get the ALG tooltip from span.punt_calc attribute "data-original-title"
		selection := e.DOM.Find(`span.punt_calc`)
		if selection.Length() > 0 {
			player.ALG = selection.Text()
		}
		// There are two span.stats_calc elements. We assume the first one is PRES and the second is F.MEDIA.
		e.ForEach("span.stats_calc", func(index int, el *colly.HTMLElement) {
			// The content of the span includes both the number and a <small> label.
			// We use the underlying DOM to filter out the <small> element.
			// This uses goquery (which Colly uses under the hood).
			text := el.DOM.Contents().Not("small").Text()
			text = strings.TrimSpace(text)
			if index == 0 {
				player.Pres = text
			} else if index == 1 {
				player.FMedia = text
			} else if index == 2 {
				player.Trend = text
			}
		})

		players = append(players, player)
	})

	// Log any error that occurred during the request
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error: %s", err)
	})

	// Start scraping on the provided URL.
	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	// Wait until threads are finished.
	c.Wait()

	return players, nil
}

func main() {
	role := "Attaccanti"

	url := "https://www.fantacalciopedia.com/lista-calciatori-serie-a/" + role + "/"

	players, err := scrapePlayers(url)
	if err != nil {
		log.Fatalf("Failed to scrape: %v", err)
	}

	// Print the scraped player data.
	println("Nome,Role,Team,ALG,Pres,FMedia,Trend")
	for _, p := range players {
		println(strings.Join([]string{p.Title, role, p.Team, p.ALG, p.Pres, p.FMedia, p.Trend}, ","))
		//
		// fmt.Printf("Title:  %s\n", p.Title)
		// fmt.Printf("Team:   %s\n", p.Team)
		// fmt.Printf("ALG:    %s\n", p.ALG)
		// fmt.Printf("PRES:   %s\n", p.Pres)
		// fmt.Printf("F.MEDIA:%s\n", p.FMedia)
		// fmt.Println("---------")
	}
}
