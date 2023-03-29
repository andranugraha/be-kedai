package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	// High traffic:
	// 	This refers to a situation where a large number of users are simultaneously accessing a web server,
	// 	creating a high level of traffic. Hey can be used to simulate high traffic by sending a large number
	// 	of requests to the server at once. For example, the following command sends 3000 requests to the server
	// 	at a concurrency level of 100:
	cmd := exec.Command("./library/hey",
		"-n", "3000",
		"-c", "100",
		"https://dev-kedai-y3gq8.ondigitalocean.app/v1/locations/provinces")
	output, err := cmd.Output()
	if err != nil {
		log.Panicln(err)
		return
	}
	fmt.Println(string(output))

	// Bursty traffic:
	// 	This refers to a situation where a sudden surge of traffic occurs, followed by a period of lower traffic.
	// 	Hey can be used to simulate bursty traffic by specifying a high rate of requests for a short period of time,
	// 	followed by a lower rate of requests. For example, the following command sends 1000 requests at a rate of 100
	// 	requests per second for 10 seconds, followed by a 5 second delay, and then another burst of 1000 requests:
	cmd = exec.Command("./library/hey",
		"-n", "1000",
		"-c", "10",
		"-q", "100",
		"https://dev-kedai-y3gq8.ondigitalocean.app/v1/locations/provinces")
	output, err = cmd.Output()
	if err != nil {
		log.Panicln(err)
		return
	}
	fmt.Println(string(output))

	// Sustained traffic: This refers to a situation where a consistent level of traffic is maintained over a longer period
	// 	of time. Hey can be used to simulate sustained traffic by sending a steady stream of requests at a constant rate.
	// 	For example, the following command sends requests at a rate of 100 requests per second for 30 seconds:
	cmd = exec.Command("./library/hey",
		"-q", "100",
		"-z", "30s",
		"https://dev-kedai-y3gq8.ondigitalocean.app/v1/locations/provinces")
	output, err = cmd.Output()
	if err != nil {
		log.Panicln(err)
		return
	}
	fmt.Println(string(output))
}
