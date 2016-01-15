package dockercloud

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
)

const (
	appSubsystemVersion   = "v1"
	infraSubsytemVersion  = "v1"
	auditSubsystemVersion = "v1"
	repoSubsystemVersion  = "v1"
	buildSubsystemVersion = "v1"
)

var (
	User       string
	ApiKey     string
	BasicAuth  string
	AuthHeader string
	BaseUrl    = "https://cloud.docker.com/"
	StreamUrl  = "wss://cloud.docker.com/"
	version    = "0.21.0"
)

/*type config map[string]Auth

type Auth struct {
	User       string
	Apikey     string
	Basic_auth string
}*/

type AuthConfig struct {
	Auths Config `json:"auths"`
}

type Config map[string]AuthObject

type AuthObject struct {
	Email string `json:"email"`
	Auth  string `json:"auth"`
}

func LoadAuth() error {
	if User != "" && ApiKey != "" {
		sEnc := base64.StdEncoding.EncodeToString([]byte(User + ":" + ApiKey))
		AuthHeader = fmt.Sprintf("Basic %s", sEnc)
		return nil
	} else {
		if os.Getenv("TUTUM_AUTH") != "" {
			AuthHeader = os.Getenv("TUTUM_AUTH")
			return nil
		}
		// Load environment variables as an alternative option
		if os.Getenv("DOCKERCLOUD_USER") != "" && os.Getenv("DOCKERCLOUD_APIKEY") != "" {
			User = os.Getenv("DOCKERCLOUD_USER")
			ApiKey = os.Getenv("DOCKERCLOUD_APIKEY")
			sEnc := base64.StdEncoding.EncodeToString([]byte(User + ":" + ApiKey))
			AuthHeader = fmt.Sprintf("Basic %s", sEnc)
			return nil
		}
		if usr, err := user.Current(); err == nil {
			var conf AuthConfig
			confFilePath := usr.HomeDir + "/.docker/"
			if _, err := os.Stat(confFilePath + "config.json"); err == nil {
				file, e := ioutil.ReadFile(confFilePath + "config.json")
				if e != nil {
					log.Println(e)
				}
				err = json.Unmarshal(file, &conf)
				if err != nil {
					return err
				}
				sEnc := base64.StdEncoding.EncodeToString([]byte(conf.Auths["https://index.docker.io/v1/"].Email + ":" + conf.Auths["https://index.docker.io/v1/"].Auth))
				AuthHeader = fmt.Sprintf("Basic %s", sEnc)
				return nil
			}
		}
	}
	return fmt.Errorf("Couldn't find any DockerCloud credentials in ~/.tutum or environment variables DOCKERCLOUD_USER and DOCKERCLOUD_APIKEY")
}

func IsAuthenticated() bool {
	return (AuthHeader != "")
}
