package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func Fetch(url string, accountID string) {
	response, err := http.Get(url + "/" + accountID)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(responseData))

	var data ResponseData
	err = json.Unmarshal(responseData, &data)
	if err != nil {
		panic(err)
	}
	fmt.Println(data.Data)
}

func Create(url string) {

	country := "GB"
	name := [4]string{
		"Dong",
		"Wang",
		"abc",
		"def",
	}

	attributes := AccountAttributes{
		Name:         name[:],
		Country:      &country,
		BaseCurrency: "GBP",
		BankID:       "400300",
		BankIDCode:   "GBDSC",
		Bic:          "NWBKGB22",
	}

	accountData := AccountData{
		Type:           "accounts",
		ID:             "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
		OrganisationID: "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
		Attributes:     &attributes,
	}

	data := RequestData{Data: &accountData}

	dataJson, err := json.Marshal(data)

	response, err := http.Post(url, "application/vnd.api+json", bytes.NewBuffer(dataJson))
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(responseData))
}

func main() {
	url := "http://localhost:8080/v1/organisation/accounts"
	// Create(url)
	Fetch(url, "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc")
}
