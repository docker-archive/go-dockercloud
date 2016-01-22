package dockercloud

import "encoding/json"

func GetImageTag(name string, tag string) (ImageTags, error) {
	url := ""
	if string(name[0]) == "/" {
		url = name[5:]
	} else {
		url = "repo/" + repoSubsystemVersion + "/repository/" + name + "/tag/" + tag + "/"
	}

	request := "GET"
	//Empty Body Request
	body := []byte(`{}`)
	var response ImageTags

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

func GetImageBuildSource(uuid string) (BuildSource, error) {
	url := "build/" + buildSubsystemVersion + "/source/" + uuid + "/"
	request := "GET"
	body := []byte(`{}`)
	var response BuildSource

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

func GetImageBuildSetting(uuid string) (BuildSettings, error) {

	url := "build/" + buildSubsystemVersion + "/setting/" + uuid + "/"

	request := "GET"
	//Empty Body Request
	body := []byte(`{}`)
	var response BuildSettings

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

/*func (self *BuildSettings) Build() (BuildSettings, error) {
	url := "build/" + buildSubsystemVersion + "/setting/" + self.Uuid + "/call/"
	request := "POST"

	body := []byte(`{}`)
	var response BuildSettings

	data, err := DockerCloudCall(url, request, body)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(data, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}*/
