package data

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"sort"

	"github.com/disgoorg/log"
)

const (
	dataURL = "https://raw.githubusercontent.com/sebm253/FlagGuessr/main/data.json"
)

var (
	popBoundaries = []int{10000000, 5000000, 1000000, 500000, 250000, 100000, 50000, 10000}
)

type CountryData struct {
	Countries []*Country

	indexBoundaries map[int]int
}

func (d *CountryData) Populate() {
	rs, err := http.Get(dataURL)
	if err != nil {
		panic(err)
	}
	defer rs.Body.Close()
	if err := json.NewDecoder(rs.Body).Decode(&d.Countries); err != nil {
		panic(err)
	}
	sort.Slice(d.Countries, func(i, j int) bool {
		return d.Countries[i].Population > d.Countries[j].Population
	})

	d.indexBoundaries = make(map[int]int)

	curBoundaryIndex := 0
	for i, country := range d.Countries {
		if curBoundaryIndex == len(popBoundaries) {
			break
		}
		popBoundary := popBoundaries[curBoundaryIndex]
		if popBoundary > country.Population {
			d.indexBoundaries[popBoundary] = i - 1
			curBoundaryIndex++
		}
	}
	d.indexBoundaries[0] = len(d.Countries)
	log.Infof("loaded %d countries.", len(d.Countries))
}

func (d *CountryData) GetRandomCountry(minPopulation int) (int, *Country) {
	i := rand.Intn(d.indexBoundaries[minPopulation])
	return i, d.Countries[i]
}
