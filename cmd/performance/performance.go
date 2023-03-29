package main

import (
	"fmt"
	"kedai/backend/be-kedai/config"
	"kedai/backend/be-kedai/internal/server"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func createReport(content string, filename string) {
	// Create the directory recursively
	err := os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	if err != nil {
		panic(err)
	}

	// Create the file if it doesn't exist
	file, err := os.Create(filename + ".txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func generateEndpointPath(r *gin.Engine, filename string) gin.RoutesInfo {
	routes := r.Routes()
	// Generate output text
	var output strings.Builder
	for _, route := range routes {
		if route.Method == "GET" || route.Method == "POST" || route.Method == "PUT" || route.Method == "DELETE" {
			output.WriteString(route.Method + "\t" + route.Path + "\n")
		}
	}
	// Generate a text file
	createReport(output.String(), filename)

	return routes
}

func main() {
	const REPORT_DIR = "./cmd/performance/report"
	var BASE_URL = config.GetEnv("BACKEND_URL", "https://dev-kedai-y3gq8.ondigitalocean.app")

	// Clear previous reports
	err := os.RemoveAll(REPORT_DIR)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Report reset succesfully!")
	}

	routes := generateEndpointPath(server.NewRouter(&server.RouterConfig{}), REPORT_DIR+"/endpoint-paths")

	for i, route := range routes {
		fmt.Print("Testing on endpoint " + strconv.Itoa(i) + " of " + strconv.Itoa(len(routes)) + ": " + route.Method + " " + route.Path)
		// High traffic:
		// 	This refers to a situation where a large number of users are simultaneously accessing a web server,
		// 	creating a high level of traffic. Hey can be used to simulate high traffic by sending a large number
		// 	of requests to the server at once. For example, the following command sends 100 requests to the server
		// 	at a concurrency level of 10:
		cmd := exec.Command("./library/hey",
			"-n", "100",
			"-c", "10",
			"-m", route.Method,
			BASE_URL+route.Path)
		output, err := cmd.Output()
		if err != nil {
			log.Panicln(err)
			return
		}
		createReport(string(output), REPORT_DIR+"/high-traffic/"+route.Method+route.Path)

		// Bursty traffic:
		//
		//	This refers to a situation where a sudden surge of traffic occurs, followed by a period of lower traffic.
		//	Hey can be used to simulate bursty traffic by specifying a high rate of requests for a short period of time,
		//	followed by a lower rate of requests. For example, the following command sends 100 requests at a rate of 10
		//	requests per second for 10 seconds, and then another burst of 100 requests:
		cmd = exec.Command("./library/hey",
			"-n", "100",
			"-c", "10",
			"-q", "10",
			"-m", route.Method,
			BASE_URL+route.Path)
		output, err = cmd.Output()
		if err != nil {
			log.Panicln(err)
			return
		}
		createReport(string(output), REPORT_DIR+"/bursty-traffic/"+route.Method+route.Path)

		// Sustained traffic: This refers to a situation where a consistent level of traffic is maintained over a longer period
		//
		//	of time. Hey can be used to simulate sustained traffic by sending a steady stream of requests at a constant rate.
		//	For example, the following command sends requests at a rate of 5 requests per second for 10 seconds:
		cmd = exec.Command("./library/hey",
			"-q", "100",
			"-z", "5s",
			"-m", route.Method,
			BASE_URL+route.Path)
		output, err = cmd.Output()
		if err != nil {
			log.Panicln(err)
			return
		}
		createReport(string(output), REPORT_DIR+"/sustained-traffic/"+route.Method+route.Path)

		fmt.Println(": DONE")
	}
}
