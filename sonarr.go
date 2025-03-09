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

type Series struct {
	Genres     []string   `json:"genres"`
	Name       string     `json:"title"`
	Year       uint16     `json:"year"`
	Statistics Statistics `json:"statistics"`
}

type Statistics struct {
	SeasonCount       uint8  `json:"seasonCount"`
	TotalEpisodeCount uint16 `json:"totalEpisodeCount"`
}

func fetchSonarr() []Series {

	req, err := http.NewRequest("GET", "https://sonarr.tjaardflix.de/api/v3/series", nil)
	req.Header.Add("X-Api-Key", "e165e362579c4617933b0e3426b2c0fa")

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

	var data []Series
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	}

	return data

}

func sendSonarr(series []Series) {
	client := influxdb.NewClient(influxDBURL, token)
	defer client.Close()

	api := client.WriteAPIBlocking(org, bucket)

	for _, serie := range series {

		p := influxdb.NewPointWithMeasurement("series").
			AddTag("name", serie.Name).
			AddField("genres", serie.Genres).
			AddField("year", serie.Year).
			AddField("episodeCount", serie.Statistics.TotalEpisodeCount).
			AddField("seasonCount", serie.Statistics.SeasonCount).
			SetTime(time.Now())

		if err := api.WritePoint(context.Background(), p); err != nil {
			fmt.Println("Error writing point to InfluxDB:", err)
		}
	}
}
