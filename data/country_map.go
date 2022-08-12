package data

import (
	"io"
	"net/http"

	"github.com/disgoorg/disgo/json"
	"github.com/disgoorg/log"
)

var (
	apiUrl = "https://restcountries.com/v3.1/all"
)

var CountryMap = make(map[string]Country)

func PopulateCountryMap() {
	var countries []Country
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
	err = json.Unmarshal(body, &countries)
	if err != nil {
		panic(err)
	}
	for _, country := range countries {
		CountryMap[country.Cca2] = country
	}
	log.Infof("loaded %d countries.", len(CountryMap))
}
