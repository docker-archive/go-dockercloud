package dockercloud

import (
	"encoding/json"
	"log"
)

func ListSwarms() (SwarmListResponse, error) {

	url := ""
	if Namespace != "" {
		url = "infra/" + infraSubsytemVersion + "/" + Namespace + "/swarm/"
	} else {
		url = "infra/" + infraSubsytemVersion + "/swarm/"
	}

	request := "GET"
	//Empty Body Request
	body := []byte(`{}`)
	var response SwarmListResponse
	var finalResponse SwarmListResponse

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
			var nextResponse SwarmListResponse
			data, err := DockerCloudCall(response.Meta.Next[5:], request, body)
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

func GetSwarm(id string) (Swarm, error) {

	url := ""
	if string(id[0]) == "/" {
		url = id[5:]
	} else if Namespace != "" {
		url = "infra/" + infraSubsytemVersion + "/" + Namespace + "/swarm/" + id + "/"
	} else {
		url = "infra/" + infraSubsytemVersion + "/swarm/" + id + "/"
	}

	request := "GET"
	//Empty Body Request
	body := []byte(`{}`)
	var response Swarm

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

func CreateSwarm(createRequest SwarmCreateRequest, provider SwarmProviderOptions) (Swarm, error) {
	url := ""
	if Namespace != "" {
		url = "infra/" + infraSubsytemVersion + "/" + Namespace + "/swarm/" + provider.Name + "/" + provider.Region + "/"
	} else {
		url = "infra/" + infraSubsytemVersion + "/swarm/" + provider.Name + "/" + provider.Region + "/"
	}

	request := "POST"
	var response Swarm

	newSwarm, err := json.Marshal(createRequest)
	if err != nil {
		return response, err
	}
	log.Println(string(newSwarm))

	data, err := DockerCloudCall(url, request, newSwarm)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(data, &response)
	if err != nil {
		return response, err
	}

	return response, nil

}
