package dockercloud

import "encoding/json"

func CreateToken() (Token, error) {
	url := "token/"
	request := "POST"
	body := []byte(`{}`)
	var response Token

	data, err := DockerCloudCall(url, request, body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(data, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}
