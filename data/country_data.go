package data

import (
	"encoding/json"
	"github.com/disgoorg/log"
	"net/http"
	"sort"
)

const (
	apiUrl = "https://restcountries.com/v3.1/all"
)

var IndexBoundaries = make(map[int]int)
var CountrySlice = make([]Country, 0)
var popBoundaries = []int{10000000, 5000000, 1000000, 500000, 250000, 100000, 50000, 10000}

func PopulateCountries() {
	response, err := http.Get(apiUrl)
	if err != nil {
		panic(err)
	}
	body := response.Body
	defer body.Close()
	if err = json.NewDecoder(body).Decode(&CountrySlice); err != nil {
		panic(err)
	}
	sort.Slice(CountrySlice, func(i, j int) bool {
		return CountrySlice[i].Population > CountrySlice[j].Population
	})
	currentBoundaryIndex := 0
	for i, country := range CountrySlice {
		if currentBoundaryIndex == len(popBoundaries) {
			break
		}
		popBoundary := popBoundaries[currentBoundaryIndex]
		if popBoundary > country.Population {
			IndexBoundaries[popBoundary] = i - 1
			currentBoundaryIndex++
		}
	}
	l := len(CountrySlice)
	IndexBoundaries[0] = l
	log.Infof("loaded %d countries.", l)
}
