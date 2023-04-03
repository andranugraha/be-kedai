package main

import (
	"encoding/json"
	"fmt"
	userDTO "kedai/backend/be-kedai/internal/domain/user/dto"
	"net/http"
	"os"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
)

const backendURL = "https://dev-kedai-y3gq8.ondigitalocean.app/v1"

func seedUser(length int) {
	users := make([]userDTO.UserRegistrationRequest, length)
	recordDir := "./cmd/seed/record"

	if _, err := os.Stat(recordDir); os.IsNotExist(err) {
		err := os.Mkdir(recordDir, 0755)
		if err != nil {
			fmt.Printf("Error creating directory: %v\n", err)
			return
		}
	}

	filename := recordDir + "/user_data.txt"

	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	for i := range users {
		users[i] = userDTO.UserRegistrationRequest{
			Email:    fmt.Sprintf("%s.%s%d@mail.com", strings.ToLower(gofakeit.FirstName()), strings.ToLower(gofakeit.LastName()), i+1),
			Password: gofakeit.Password(true, true, true, false, false, 10),
		}

		jsonData, err := json.Marshal(users[i])
		if err != nil {
			fmt.Printf("Error marshaling JSON: %v\n", err)
			continue
		}

		resp, err := http.Post(backendURL+"/users/register", "application/json", strings.NewReader(string(jsonData)))
		if err != nil {

			fmt.Printf("Error sending HTTP request: %v\n", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
			continue
		}

		line := fmt.Sprintf("%s %s\n", users[i].Email, users[i].Password)
		_, err = file.WriteString(line)
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", err)
			return
		}

		if i%1000 == 0 {
			fmt.Printf("Seeded %d users\n", i)
		}

		if i == length-1 {
			fmt.Printf("Finished seeding %d users\n", i+1)
		}
	}

	fmt.Println("Success")

}

func main() {
	seedUser(20000)
}
