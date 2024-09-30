package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Route struct {
	Method string `json:"method"`
	Route  string `json:"route"`
}
type Routes struct {
	Routes []Route `json:"routes"`
	ApiId  string  `json:"apiId"`
}

func readConfig(path string) (Routes, error) {
	configPath := path + "/config.json"
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return Routes{}, err
	}
	var routes Routes
	if err = json.Unmarshal(configData, &routes); err != nil {
		return Routes{}, err
	}
	return routes, nil
}

func main() {
	path, err := os.Getwd()
	fmt.Println(path)
	if err != nil {
		fmt.Println("Primer error")
	}
	var routes Routes
	routes, err = readConfig(path)
	if err != nil {
		fmt.Println("Segundo error")
	}
	fmt.Println("ApiId: ", routes.ApiId)
	for _, route := range routes.Routes {
		fmt.Println("Route: ", route)
		fmt.Println("Method: ", route.Method)
		fmt.Println("Route: ", route.Route)
	}
}
