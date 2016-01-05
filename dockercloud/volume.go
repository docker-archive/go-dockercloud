package dockercloud

import "encoding/json"

func ListVolumes() (VolumeListResponse, error) {
	url := "infra/" + infraSubsytemVersion + "/volume/"
	request := "GET"
	//Empty Body Request
	body := []byte(`{}`)
	var response VolumeListResponse
	var finalResponse VolumeListResponse

	data, err := DockerCloudCall(url, request, body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(data, &response)
	if err != nil {
		return response, err
	}

	finalResponse = response

Loop:
	for {
		if response.Meta.Next != "" {
			var nextResponse VolumeListResponse
			data, err := DockerCloudCall(response.Meta.Next[8:], request, body)
			if err != nil {
				return nextResponse, err
			}
			err = json.Unmarshal(data, &nextResponse)
			if err != nil {
				return nextResponse, err
			}
			finalResponse.Objects = append(finalResponse.Objects, nextResponse.Objects...)
			response = nextResponse

		} else {
			break Loop
		}
	}

	return finalResponse, nil
}

func GetVolume(uuid string) (Volume, error) {

	url := ""
	if string(uuid[0]) == "/" {
		url = uuid[8:]
	} else {
		url = "infra/" + infraSubsytemVersion + "/volume/" + uuid + "/"
	}

	request := "GET"
	//Empty Body Request
	body := []byte(`{}`)
	var response Volume

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
