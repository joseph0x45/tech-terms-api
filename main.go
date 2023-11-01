package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gocolly/colly/v2"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	c := colly.NewCollector()

	r.Get("/search/{term}", func(w http.ResponseWriter, r *http.Request) {
		term := chi.URLParam(r, "term")
		if term == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		c.OnHTML("article", func(h *colly.HTMLElement) {
			w.Write([]byte(h.Text))
			w.WriteHeader(http.StatusOK)
			return
		})

		c.OnError(func(r *colly.Response, err error) {
			error_code := err.Error()
			if error_code == "Not Found" {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		})

		url := fmt.Sprintf("https://techterms.com/definition/%s", term)
		c.Visit(url)
	})

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}

}