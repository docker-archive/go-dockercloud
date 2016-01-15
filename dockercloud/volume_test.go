package dockercloud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_ListVolumes(t *testing.T) {
	User = "test"
	ApiKey = "test"

	fake_response, err := MockupResponse("listvolumes.json")
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, fake_response)
	}))

	defer server.Close()
	url := server.URL + "/api/infra/" + infraSubsytemVersion + "/volume/"

	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var response VolumeListResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatal(err)
	}

	BaseUrl = server.URL + "/api/"

	test_response, err := ListVolumes()
	if err != nil {
		t.Fatal(err)
	}

	if reflect.DeepEqual(test_response, response) != true {
		t.Fatal("Invalid output")
	}
}

func Test_GetVolume(t *testing.T) {
	User = "test"
	ApiKey = "test"

	fake_response, err := MockupResponse("volume.json")
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, fake_response)
	}))

	defer server.Close()
	url := server.URL + "/api/infra/" + infraSubsytemVersion + "/volume/" + fake_uuid_volume

	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var response Volume
	err = json.Unmarshal(body, &response)
	if err != nil {
		t.Fatal(err)
	}

	BaseUrl = server.URL + "/api/"
	test_response, err := GetVolume(fake_uuid_volume)
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(test_response, response) != true {
		t.Fatal("Invalid output")
	}
}
