package kotsd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type KotsdConfig struct {
	Configs []Instance `yaml:"instances"`
}
type Instance struct {
	Name        string `yaml:"name"`
	Endpoint    string `yaml:"endpoint"`
	Password    string `yaml:"password"`
	KotsVersion string `yaml:"-"`
}

func ReadConfig(cfgFile string) ([]byte, error) {
	f, err := os.ReadFile(cfgFile)

	if err != nil {
		log.Fatal(err)
	}
	return f, nil

}

func ParseConfig(d []byte) (KotsdConfig, error) {
	var rc KotsdConfig
	err := yaml.Unmarshal(d, &rc)

	if err != nil {
		log.Fatal(err)
	}

	return rc, nil
}

func WriteConfig(config *KotsdConfig, cfgFile string) {
	data, err := yaml.Marshal(&config)

	if err != nil {
		log.Fatal(err)
	}

	err2 := os.WriteFile(cfgFile, data, 0)

	if err2 != nil {
		log.Fatal(err2)
	}
}

func (kc *KotsdConfig) AddInstance(name string, endpoint string, password string) {
	instance := Instance{Name: name, Endpoint: endpoint, Password: base64.StdEncoding.EncodeToString([]byte(password))}
	instances := kc.Configs
	instances = append(instances, instance)
	kc.Configs = instances
}

type HealthzResponse struct {
	Version string `json:"version"`
	GitSHA  string `json:"gitSha"`
}

func (instance Instance) GetKotsHealthz() (HealthzResponse, error) {
	url := fmt.Sprintf("%s/healthz", instance.Endpoint)
	req, err := http.NewRequest(
		"GET",
		url,
		nil,
	)

	if err != nil {
		return HealthzResponse{}, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return HealthzResponse{}, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	if resp.StatusCode != 200 {
		return HealthzResponse{}, fmt.Errorf("GET /healthz %d: %s", resp.StatusCode, body)
	}

	var healthz HealthzResponse
	err = json.Unmarshal([]byte(body), &healthz)
	if err != nil {
		return HealthzResponse{}, err
	}

	return healthz, nil

}
