package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/gocolly/colly/v2"
)

type FilmsDispo struct {
	Films []Film
}

type Film struct {
	Title     string
	Duration  string
	Director  string
	Genre     string
	Showtimes []Showtime
}
type Showtime struct {
	Day      string
	Horaires string
}

var date = time.Now()
var date2 = time.Now()

func truncateString(str string, num int) string {
	if len(str) > num {
		return str[:num]
	}
	return str
}

func main() {
	test := scrape2()

	for _, film := range test.Films {
		fmt.Printf("Titre: %s\nDurée: %s\nRéalisateur: %s\nGenre: %s\n", film.Title, film.Duration, film.Director, film.Genre)
		for _, showtime := range film.Showtimes {
			fmt.Printf("Jour: %s, Séances: %s\n", showtime.Day, showtime.Horaires)
		}
		fmt.Println()
	}

}

func scrape2() FilmsDispo {
	c := colly.NewCollector()
	var filmsDispo FilmsDispo

	c.OnHTML(".hr_film.aff_l", func(e *colly.HTMLElement) {
		film := Film{
			Title:    e.ChildText("h2 a"),
			Duration: truncateString(e.ChildText(".hr_dur strong"), 4),
			Director: e.ChildText(".hr_real strong"),
			Genre:    e.ChildText(".genre strong"),
		}

		e.ForEach(".tab_seances", func(i int, s *colly.HTMLElement) {
			jour := s.Attr("class")
			re := regexp.MustCompile(`\d+`)
			jour = re.FindString(jour)

			horaires := s.ChildText(".hor")
			if horaires != "" {
				jourInt, err := strconv.Atoi(jour)
				if err != nil {
					fmt.Printf("Invalid day format: %s\n", jour)
					return
				}

				date = date.Add(time.Hour * time.Duration(jourInt) * 24).Add(-time.Hour * 24)
				reHoraires := regexp.MustCompile(`(\d{2}h\d{2})`)
				horairesCorrige := reHoraires.ReplaceAllString(horaires, "$1 ")

				showtime := Showtime{
					Day:      date.Format("2006-01-02"),
					Horaires: horairesCorrige,
				}
				film.Showtimes = append(film.Showtimes, showtime)
				date = date2
			}
		})
		filmsDispo.Films = append(filmsDispo.Films, film)
	})

	c.Visit("https://www.cine-aire.com/horaires/")

	/*
		exemples :
			https://www.cinemas-legrandclub.fr/dax/horaires/
			https://www.cinemas-legrandclub.fr/mont-de-marsan/horaires/
			https://www.cinemas-legrandclub.fr/cap-breton/horaires/
			https://www.cinemas-legrandclub.fr/hossegor/horaires/
			https://www.cine-aire.com/horaires/
	*/

	return filmsDispo
}
