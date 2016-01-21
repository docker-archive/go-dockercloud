package dockercloud

import (
	"os"
	"testing"
)

func Test_SetBaseUrl(t *testing.T) {

	url := ""

	os.Setenv("DOCKERCLOUD_REST_HOST", "https://cloud.docker.com")
	url = SetBaseUrl()
	if url != "https://cloud.docker.com/api/" {
		t.Fatal("Wrong url set")
	}
	os.Setenv("DOCKERCLOUD_REST_HOST", "")
	os.Setenv("DOCKERCLOUD_BASE_URL", "https://cloud.docker.com/api/")
	url = SetBaseUrl()
	if url != "https://cloud.docker.com/api/" {
		t.Fatal("Wrong url set")
	}
	os.Setenv("DOCKERCLOUD_BASE_URL", "")
}
