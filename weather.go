package main

import (
	"fmt"
	"http"
	"io/ioutil"
	"os"
	"xml"
	"io"
)

const (
	WeatherUrl = "http://api.wunderground.com/auto/wui/geo/ForecastXML/index.xml?query=55116"
)

type Result struct {
	XMLName xml.Name "forecast"
	Forecastday []Forecastday "txt_forecast>forecastday"
}

type Forecastday struct {
	Fcttext string
	Title string
}

func printResp(resp *http.Response) {
	resp, err := http.Get(WeatherUrl)
	s, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fatal("Error reading xml data", err)
	}
	fmt.Printf(string(s))
}

func ParseWeatherResponse(r io.Reader) Result {
	var result Result
	xml.Unmarshal(r, &result)
	return result
}

func main() {
	resp, err := http.Get(WeatherUrl)
	if err != nil {
		fatal("Error pulling down weather data.", err)
	}
	result := ParseWeatherResponse(resp.Body)
	for i := 0; i < len(result.Forecastday); i++ {
		var day = result.Forecastday[i]
		fmt.Printf("%s:\n*******************\n%s\n\n", day.Title, day.Fcttext)
	}	
}

func fatal(s string, a ... interface{}) { 
        fmt.Fprintf(os.Stderr, "netfwd: %s\n", fmt.Sprintf(s, a)) 
        os.Exit(2) 
}
 
