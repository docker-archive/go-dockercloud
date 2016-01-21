package dockercloud

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var customUserAgent = "go-dockercloud/" + version

func SetUserAgent(name string) string {
	customUserAgent = ""
	customUserAgent = name + " go-dockercloud/" + version
	return customUserAgent
}

func SetBaseUrl() string {
	if os.Getenv("DOCKERCLOUD_REST_HOST") != "" {
		BaseUrl = os.Getenv("DOCKERCLOUD_REST_HOST")
		BaseUrl = BaseUrl + "/api/"
	} else if os.Getenv("DOCKERCLOUD_BASE_URL") != "" {
		BaseUrl = os.Getenv("DOCKERCLOUD_BASE_URL")
	}
	return BaseUrl
}

func DockerCloudCall(url string, requestType string, requestBody []byte) ([]byte, error) {

	LoadAuth()

	if !IsAuthenticated() {
		return nil, fmt.Errorf("Couldn't find any DockerCloud credentials in ~/.docker/config.json or environment variables DOCKERCLOUD_USER and DOCKERCLOUD_APIKEY")
	}

	BaseUrl = SetBaseUrl()

	client := &http.Client{}
	req, err := http.NewRequest(requestType, BaseUrl+url, bytes.NewBuffer(requestBody))

	req.Header.Add("Authorization", AuthHeader)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("User-Agent", customUserAgent)

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode > 300 {
		return nil, fmt.Errorf("Failed API call: %s ", response.Status)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
