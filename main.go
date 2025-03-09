package main

import (
	"net/http"
)

var influxDBURL = "http://server.tjaardflix.de:8086"
var token = "9bf4w79tgfd2397ghfdsamjd9o75ktdwirf5"
var bucket = "home"
var org = "my-org"

var client = &http.Client{}

func main() {
	// data := fetchRadarr()
	// sendRadarr(data)
	sonarr := fetchSonarr()
	sendSonarr(sonarr)
}
