package form3client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type Form3Client struct {
	client *http.Client
}

func (c *Form3Client) New() *http.Client {
	c.client = &http.Client{}
	return c.client
}

func (c *Form3Client) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error to create new request: %s", err.Error())
	}

	return request, nil
}

func (c *Form3Client) Request(request *http.Request) (*http.Response, error) {
	response, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error to request data: %s", err.Error())
	}

	return response, nil
}

func (c *Form3Client) SetQueryParams(request *http.Request, queryParams map[string]string) {
	query := request.URL.Query()
	for param, value := range queryParams {
		query.Add(param, value)
	}
	request.URL.RawQuery = query.Encode()
}

func (c *Form3Client) Response(response *http.Response) ([]byte, error) {
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error to read response body: %s", err.Error())
	}

	return responseData, nil
}

func (c *Form3Client) Fetch(url string, accountID string) (*AccountData, error) {
	request, err := c.NewRequest("GET", fmt.Sprintf("%s/%s", url, accountID), nil)
	if err != nil {
		return nil, fmt.Errorf("fetch error: %s", err.Error())
	}

	response, err := c.Request(request)
	if err != nil {
		return nil, fmt.Errorf("fetch error: %s", err.Error())
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch error: status == %d, text == %s",
			response.StatusCode, http.StatusText(response.StatusCode))
	}

	responseData, err := c.Response(response)
	if err != nil {
		return nil, fmt.Errorf("fetch error: %s", err.Error())
	}

	var data ResponseData
	err = json.Unmarshal(responseData, &data)
	if err != nil {
		return nil, fmt.Errorf("fetch error: read response body: %s", err.Error())
	}

	return data.Data, nil
}

func (c *Form3Client) Delete(url string, accountID string, queryParams map[string]string) (bool, error) {
	request, err := c.NewRequest("DELETE", fmt.Sprintf("%s/%s", url, accountID), nil)
	if err != nil {
		return false, fmt.Errorf("delete error: %s", err.Error())
	}

	c.SetQueryParams(request, queryParams)

	response, err := c.Request(request)
	if err != nil {
		return false, fmt.Errorf("delete error: %s", err.Error())
	}

	if response.StatusCode != http.StatusNoContent {
		return false, fmt.Errorf("delete error: status == %d, text == %s",
			response.StatusCode, http.StatusText(response.StatusCode))
	}

	return response.StatusCode == http.StatusNoContent, nil
}

func (c *Form3Client) Create(url string, accountData AccountData) (*AccountData, error) {
	data := RequestData{Data: &accountData}
	dataJson, err := json.Marshal(data)
	request, err := c.NewRequest("POST", url, bytes.NewBuffer(dataJson))
	request.Header.Add("Content-Type", "application/vnd.api+json")

	if err != nil {
		return nil, fmt.Errorf("create error: %s", err.Error())
	}

	response, err := c.Request(request)
	if err != nil {
		return nil, fmt.Errorf("create error: %s", err.Error())
	}

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("create error: status == %d, text == %s",
			response.StatusCode, http.StatusText(response.StatusCode))
	}

	responseData, err := c.Response(response)
	if err != nil {
		return nil, fmt.Errorf("create error: %s", err.Error())
	}

	var newAccountData ResponseData
	err = json.Unmarshal(responseData, &newAccountData)
	if err != nil {
		return nil, fmt.Errorf("create error: %s", err.Error())
	}

	return newAccountData.Data, nil
}
