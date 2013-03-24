package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	WeatherUrl = "http://api.wunderground.com/api/573832cedbb28381/conditions/forecast/q/%s.xml"
)

type Result struct {
	XMLName    xml.Name         `xml:"response"`
	TextCast   []TextForecast   `xml:"forecast>txt_forecast>forecastdays>forecastday"`
	SimpleCast []SimpleForecast `xml:"forecast>simpleforecast>forecastdays>forecastday"`
	Location   string           `xml:"current_observation>display_location>full"`
	Temp       string           `xml:"current_observation>temp_f"`
	Weather    string           `xml:"current_observation>weather"`
}

type SimpleForecast struct {
	High       int32  `xml:"high>fahrenheit"`
	Low        int32  `xml:"low>fahrenheit"`
	Year       int32  `xml:"date>year"`
	Day        int32  `xml:"date>day"`
	Month      int32  `xml:"date>month"`
	Weekday    string `xml:"date>weekday"`
	Conditions string `xml:"conditions"`
}

type TextForecast struct {
	Fcttext string `xml:"fcttext"`
	Title   string `xml:"title"`
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

func ParseWeatherResponse(r io.Reader) (*Result, error) {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var result Result
	err = xml.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func main() {
	location := flag.String("l", "autoip", "Weather location to query. Zip code or city,state.")
	flag.Parse()
	url := fmt.Sprintf(WeatherUrl, *location)

	resp, err := http.Get(url)
	if err != nil {
		fatal("Error pulling down weather data.", err)
	}
	result, err := ParseWeatherResponse(resp.Body)
	if err != nil {
		fatal("Error reading weather data.", err)
	}
	fmt.Print("++++++++++++++++++++++++++++++++++++++\n")
	fmt.Printf("Weather for %s\n", result.Location)
	fmt.Print("++++++++++++++++++++++++++++++++++++++\n")
	// printTextForecast(*result)
	printSimpleForcast(*result)
	fmt.Print("\n")
	printCurrentForcast(*result)
}

func printTextForecast(result Result) {
	for i := 0; i < len(result.TextCast); i++ {
		var day = result.TextCast[i]
		fmt.Printf("%s:\n*******************\n%s\n\n", day.Title, day.Fcttext)
	}
}

func printCurrentForcast(result Result) {
	if len(result.Temp) > 0 && len(result.Weather) > 0 {
		fmt.Printf("Currently : %s - %s\n\n", result.Temp, result.Weather)
	}
}

func printSimpleForcast(result Result) {
	fmt.Print("             L   /  H\n")
	for i := 0; i < len(result.SimpleCast); i++ {
		var s = result.SimpleCast[i]
		if i == 0 {
			fmt.Printf("%9s", "Today")
		} else if i == 1 {
			fmt.Printf("%9s", "Tomorrow")
		} else {
			fmt.Printf("%9s", s.Weekday)
			// fmt.Printf("%d-%02d-%02d", s.Year, s.Month, s.Day)
		}
		fmt.Printf(" - %3d  / %3d - %s\n", s.Low, s.High, s.Conditions)
	}
}

func fatal(s string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "netfwd: %s\n", fmt.Sprintf(s, a))
	os.Exit(2)
}
