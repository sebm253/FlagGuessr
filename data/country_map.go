package data

import (
	"encoding/json"
	"github.com/disgoorg/log"
	"io"
	"net/http"
	"sort"
)

var (
	apiUrl = "https://restcountries.com/v3.1/all"
)

var IndexBoundaries = make(map[int]int)
var CountrySlice = make([]Country, 0)
var boundaries = []int{10000000, 5000000, 1000000, 500000, 250000, 100000, 50000, 10000}

func PopulateCountryMap() {
	response, err := http.Get(apiUrl)
	if err != nil {
		panic(err)
	}
	closer := response.Body
	body, err := io.ReadAll(closer)
	err = closer.Close()
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &CountrySlice)
	if err != nil {
		panic(err)
	}
	sort.Slice(CountrySlice, func(i, j int) bool {
		return CountrySlice[i].Population > CountrySlice[j].Population
	})
	for i, country := range CountrySlice {
		for _, boundary := range boundaries {
			if country.Population > boundary {
				IndexBoundaries[boundary] = i
			}
		}
	}
	log.Infof("loaded %d countries.", len(CountrySlice))
}
