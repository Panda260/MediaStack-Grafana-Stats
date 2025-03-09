package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	influxdb "github.com/influxdata/influxdb-client-go/v2"
)

type Movie struct {
	Genres    []string  `json:"genres"`
	Name      string    `json:"title"`
	Year      uint16    `json:"year"`
	MovieFile MovieFile `json:"movieFile"`
}

type MovieFile struct {
	Quality QualityTop `json:"quality"`
}

type QualityTop struct {
	Quality Quality `json:"quality"`
}

type Quality struct {
	Name string `json:"name"`
}

func fetchRadarr() []Movie {

	req, err := http.NewRequest("GET", "https://radarr.tjaardflix.de/api/v3/movie", nil)
	req.Header.Add("X-Api-Key", "2540104d2aa04c90bb520ae829f4daad")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: received non-200 response status:", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
	}

	var data []Movie
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}

	return data

}

func sendRadarr(movies []Movie) {
	client := influxdb.NewClient(influxDBURL, token)
	defer client.Close()

	api := client.WriteAPIBlocking(org, bucket)

	for _, movie := range movies {

		p := influxdb.NewPointWithMeasurement("movies").
			AddTag("name", movie.Name).
			AddField("genres", movie.Genres).
			AddField("year", movie.Year).
			AddField("quality", movie.MovieFile.Quality.Quality.Name).
			SetTime(time.Now())

		if err := api.WritePoint(context.Background(), p); err != nil {
			fmt.Println("Error writing point to InfluxDB:", err)
		}
	}
}
