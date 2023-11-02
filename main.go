package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"strings"
)

func main() {
	db, err := sqlx.Connect("sqlite3", ":memory:")
	_, err = db.Exec(`
     create table definitions(
       term text not null primary key,
       definition text not null
     );
  `, nil)
	if err != nil {
		panic(err)
	}
	c := colly.NewCollector()
	http.HandleFunc("/search/", func(w http.ResponseWriter, r *http.Request) {
		request_path := r.URL.Path
		url_params := strings.Split(request_path, "/")
		if len(url_params) != 3 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		term := url_params[2]
		url := fmt.Sprintf("https://techterms.com/definition/%s", term)
		var scraped_definition string

		c.OnHTML("article", func(h *colly.HTMLElement) {
			h.ForEach("p", func(i int, h *colly.HTMLElement) {
				scraped_definition += h.Text + "\n"
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
			if error_code == "URL already visited" {
        println("A request hit cache")
				t := new(struct {
					Term       string `db:"term"`
					Definition string `db:"definition"`
				})
				err = db.Get(t, "select * from definitions where term=$1", term)
				if err != nil {
					println(err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.Write([]byte(t.Definition))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			println(err.Error())
			return
		}
		_, err = db.NamedExec(
			"insert into definitions(term, definition) values(:term, :definition)",
			struct {
				Term       string `db:"term"`
				Definition string `db:"definition"`
			}{
				Term:       term,
				Definition: scraped_definition,
			},
		)
		if err != nil {
			println("Error while caching term definition ", err.Error())
		}
		w.Write([]byte(scraped_definition))
		return
	})

	println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
