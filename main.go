package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector()
	http.HandleFunc("/search/", func(w http.ResponseWriter, r *http.Request) {
		request_path := r.URL.Path
		url_params := strings.Split(request_path, "/")
		if len(url_params) != 3 {
			w.WriteHeader(http.StatusBadRequest)
			println("Bad url")
			return
		}
		term := url_params[2]
		url := fmt.Sprintf("https://techterms.com/definition/%s", term)
		var scraped_definition string

		c.OnHTML("article", func(h *colly.HTMLElement) {
      h.ForEach("p", func(i int, h *colly.HTMLElement) {
        scraped_definition += h.Text+"\n"
      })
		})

		err := c.Visit(url)
		if err != nil {
			error_code := err.Error()
			if error_code == "Not Found" {
				w.WriteHeader(http.StatusNotFound)
				println("not found")
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			println(err.Error())
			return
		}

		w.Write([]byte(scraped_definition))
		return
	})

	println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
