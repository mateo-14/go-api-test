package main

import (
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	cities := loadCities()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			if r.URL.Path == "/search" {
				if r.URL.Query().Get("q") == "" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				q := r.URL.Query().Get("q")
				q = strings.ToLower(q)
				for _, city := range cities {
					searchLen := len(q)
					if len(city) >= searchLen && city[0:searchLen] == q {
						w.WriteHeader(http.StatusOK)
						f, err := os.Open("./cities/" + city)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							return
						}

						w.Header().Set("Content-Type", "application/json")
						defer f.Close()
						io.Copy(w, f)
						return
					}
				}

				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", nil)
}

func loadCities() []string {
	entries, _ := os.ReadDir("./cities")
	cities := make([]string, len(entries))
	for _, entry := range entries {
		cities = append(cities, entry.Name())
	}

	return cities
}
