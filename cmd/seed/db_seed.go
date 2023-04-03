package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	shopDTO "kedai/backend/be-kedai/internal/domain/shop/dto"
	userDTO "kedai/backend/be-kedai/internal/domain/user/dto"
	"net/http"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
)

const backendURL = "https://dev-kedai-y3gq8.ondigitalocean.app/v1"

type AddressRequest struct {
	Name          string `json:"name"`
	PhoneNumber   string `json:"phoneNumber"`
	Street        string `json:"street"`
	Details       string `json:"details"`
	SubdistrictID int    `json:"subdistrictId"`
	IsDefault     *bool  `json:"isDefault"`
	IsPickup      *bool  `json:"isPickup"`
}

func main() {
	// name := "Dummy"
	password := "Secret123"
	// length := 10

	// seedUserAndAddressAndShop(length, name, password)

	test := loginUser("Dummy.1@mail.com", password)

	registerShop(test, 1)
}

func registerShop(accessToken string, addressId int) {

	shop := shopDTO.CreateShopRequest{
		Name:       gofakeit.Company(),
		AddressID:  addressId,
		CourierIDs: []int{1, 2, 3},
	}

	jsonData, err := json.Marshal(shop)

	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
	}

	fmt.Println(string(jsonData))

	req, _ := http.NewRequest("POST", backendURL+"/sellers/register", strings.NewReader(string(jsonData)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending HTTP request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		return
	}
}

func loginUser(email, password string) string {
	user := userDTO.UserLogin{
		Email:    email,
		Password: password,
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return ""
	}

	resp, err := http.Post(backendURL+"/users/login", "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		fmt.Printf("Error sending HTTP request: %v\n", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		return ""
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	var data map[string]interface{}
	_ = json.Unmarshal([]byte(respBody), &data)

	// extract the accessToken and refreshToken fields from the data map
	accessToken := data["data"].(map[string]interface{})["accessToken"].(string)

	return string(accessToken)
}

func addAddress(accessToken string) {
	falseValue := false
	address := AddressRequest{
		Name:          gofakeit.FirstName(),
		PhoneNumber:   gofakeit.Phone(),
		Street:        gofakeit.Street(),
		Details:       gofakeit.Street(),
		SubdistrictID: gofakeit.IntRange(23, 103592),
		IsDefault:     &falseValue,
		IsPickup:      &falseValue,
	}

	jsonData, _ := json.Marshal(address)

	req, _ := http.NewRequest("POST", backendURL+"/users/addresses", strings.NewReader(string(jsonData)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending HTTP request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		return
	}

	fmt.Println("Success")
}

func seedUserAndAddressAndShop(length int, name string, password string) {
	users := make([]userDTO.UserRegistrationRequest, length)

	for i := range users {
		users[i] = userDTO.UserRegistrationRequest{
			Email:    fmt.Sprintf("%s.%d@mail.com", name, i+1),
			Password: password,
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

		if i%1000 == 0 {
			fmt.Printf("Seeded %d users\n", i)
		}

		if i == length-1 {
			fmt.Printf("Finished seeding %d users\n", i+1)
		}

		addAddress(loginUser(users[i].Email, users[i].Password))
	}

	fmt.Println("Success")
}
