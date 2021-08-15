package form3client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

type test struct {
	status int
}

func createClient() Form3Client {
	client := Form3Client{}
	client.New()
	return client
}

func readAccountSampleData() ([]byte, error) {
	filename := "testdata/account_sample.json"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func getAccount(sampleData []byte) *AccountData {
	data := ResponseData{}
	json.Unmarshal(sampleData, &data)
	return data.Data
}

func TestClient_Fetch(t *testing.T) {
	expectedAccountData, err := readAccountSampleData()
	if err != nil {
		t.Errorf("Error to read account sample data: %s", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(expectedAccountData)
		if r.Method != "GET" {
			t.Errorf("Expected 'GET' request, got '%s'", r.Method)
		}
	}))
	defer server.Close()

	accountId := getAccount(expectedAccountData).ID
	client := createClient()
	accountData, err := client.Fetch(server.URL, accountId)
	if err != nil {
		t.Errorf("Expected err == nil, got '%s'", err)
	}

	if accountData.ID != accountId {
		t.Errorf("Expected account id == %s, got %s", accountId, accountData.ID)
	}
}

func TestClient_Fetch_Fail(t *testing.T) {
	tests := []test{
		{http.StatusBadRequest},
		{http.StatusUnauthorized},
		{http.StatusForbidden},
		{http.StatusNotFound},
		{http.StatusInternalServerError},
		{http.StatusServiceUnavailable},
		{http.StatusGatewayTimeout},
	}

	expectedAccountData, err := readAccountSampleData()
	if err != nil {
		t.Errorf("Error to read account sample data: %s", err)
	}

	for _, tc := range tests {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tc.status)
		}))
		defer server.Close()

		accountId := getAccount(expectedAccountData).ID
		client := createClient()
		accountData, err := client.Fetch(server.URL, accountId)
		if err == nil {
			t.Errorf("Expected err == nil, got '%s'", err)
		}

		if accountData != nil {
			t.Errorf("Expected account == nil")
		}
	}
}

func TestClient_Delete(t *testing.T) {
	expectedAccountData, err := readAccountSampleData()
	if err != nil {
		t.Errorf("Error to read account sample data: %s", err)
	}
	account := getAccount(expectedAccountData)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		if r.Method != "DELETE" {
			t.Errorf("Expected 'DELETE' request, got '%s'", r.Method)
		}
	}))
	defer server.Close()

	client := createClient()
	response, err := client.Delete(server.URL, account.ID,
		map[string]string{"Version": strconv.FormatInt(*account.Version, 10)})
	if err != nil {
		t.Errorf("Expected err == nil, got '%s'", err)
	}

	if response == false {
		t.Errorf("Expected response == true, got %t", response)
	}
}

func TestClient_Delete_Fail(t *testing.T) {
	tests := []test{
		{http.StatusBadRequest},
		{http.StatusUnauthorized},
		{http.StatusForbidden},
		{http.StatusNotFound},
		{http.StatusInternalServerError},
		{http.StatusServiceUnavailable},
		{http.StatusGatewayTimeout},
	}

	expectedAccountData, err := readAccountSampleData()
	if err != nil {
		t.Errorf("Error to read account sample data: %s", err)
	}
	account := getAccount(expectedAccountData)

	for _, tc := range tests {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tc.status)
		}))
		defer server.Close()

		client := createClient()
		response, err := client.Delete(server.URL, account.ID,
			map[string]string{"Version": strconv.FormatInt(*account.Version, 10)})
		if err == nil {
			t.Errorf("Expected err == nil, got '%s'", err)
		}

		if response != false {
			t.Errorf("Expected account == nil")
		}
	}
}

func TestClient_Create(t *testing.T) {
	expectedAccountData, err := readAccountSampleData()
	if err != nil {
		t.Errorf("Error to read account sample data: %s", err)
	}
	account := getAccount(expectedAccountData)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write(expectedAccountData)
		if r.Method != "POST" {
			t.Errorf("Expected 'POST' request, got '%s'", r.Method)
		}
	}))
	defer server.Close()

	client := createClient()
	response, err := client.Create(server.URL, *account)
	if err != nil {
		t.Errorf("Expected err == nil, got '%s'", err)
	}
	if response.ID != account.ID {
		t.Errorf("Expected account == %s, got %s", response.ID, account.ID)
	}
}

func TestClient_Create_Fail(t *testing.T) {
	tests := []test{
		{http.StatusBadRequest},
		{http.StatusUnauthorized},
		{http.StatusForbidden},
		{http.StatusNotFound},
		{http.StatusRequestEntityTooLarge},
		{http.StatusTooManyRequests},
		{http.StatusInternalServerError},
		{http.StatusServiceUnavailable},
		{http.StatusGatewayTimeout},
	}

	expectedAccountData, err := readAccountSampleData()
	if err != nil {
		t.Errorf("Error to read account sample data: %s", err)
	}
	account := getAccount(expectedAccountData)

	for _, tc := range tests {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tc.status)
		}))
		defer server.Close()

		client := createClient()
		accountData, err := client.Create(server.URL, *account)
		if err == nil {
			t.Errorf("Expected err == nil, got '%s'", err)
		}

		if accountData != nil {
			t.Errorf("Expected account == nil")
		}
	}
}
