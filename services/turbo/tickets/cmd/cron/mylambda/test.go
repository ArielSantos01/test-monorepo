package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Cron struct {
	Minute  string `json:"minute"`
	Hour    string `json:"hour"`
	Day     string `json:"day"`
	Month   string `json:"month"`
	WeekDay string `json:"weekDay"`
	Year    string `json:"year"`
}
type Schedule struct {
	Cron       Cron   `json:"cron"`
	Expression string `json:"expression"`
}

func readConfig(path string) (Schedule, error) {
	configPath := path + "/config.json"
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return Schedule{}, err
	}

	var config Schedule

	if err = json.Unmarshal(configData, &config); err != nil {
		return Schedule{}, err
	}

	return config, nil
}

func normalizeField(field string) string {
	if field == "" {
		return "*"
	}
	return field
}

func main() {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Primer error")
	}
	var config Schedule
	config, err = readConfig(path)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(normalizeField(config.Cron.Minute))
	fmt.Println(normalizeField(config.Cron.Hour))
	fmt.Println(normalizeField(config.Cron.Day))
	fmt.Println(normalizeField(config.Cron.Month))
	fmt.Println(normalizeField(config.Cron.WeekDay))
	fmt.Println(normalizeField(config.Cron.Year))

}
