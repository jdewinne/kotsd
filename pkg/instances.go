package kotsd

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type KotsdConfig struct {
	Configs []Instance `yaml:"instances"`
}
type Instance struct {
	Name               string        `yaml:"name"`
	Endpoint           string        `yaml:"endpoint"`
	Password           string        `yaml:"password"`
	InsecureSkipVerify bool          `yaml:"insecureSkipVerify"`
	KotsVersion        string        `yaml:"-"`
	Apps               []Application `yaml:"-"`
	Error              string        `yaml:"-"`
}

type Application struct {
	Name            string
	Version         string
	PendingVersions []PendingVersion
}

type PendingVersion struct {
	VersionLabel string
	Sequence     int
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

func (kc *KotsdConfig) AddInstance(name string, endpoint string, password string, tlsVerify bool) error {
	in, _ := kc.GetInstance(name)
	if in != nil {
		return errors.New("instance name must be unique")
	}
	instance := Instance{Name: name, Endpoint: endpoint, Password: base64.StdEncoding.EncodeToString([]byte(password)), InsecureSkipVerify: !tlsVerify}
	instances := kc.Configs
	instances = append(instances, instance)
	kc.Configs = instances
	return nil
}

func (kc *KotsdConfig) DeleteInstance(name string) {
	instances := []Instance{}
	for _, instance := range kc.Configs {
		if instance.Name == name {
			continue
		}
		instances = append(instances, instance)
	}
	kc.Configs = instances
}

func (kc *KotsdConfig) GetInstance(name string) (*Instance, error) {
	for _, instance := range kc.Configs {
		if instance.Name == name {
			return &instance, nil
		}
	}
	return nil, errors.New("get instance")
}

func (kc *KotsdConfig) UpdateInstance(name string, endpoint string, password string, tlsVerify bool) error {
	for i, instance := range kc.Configs {
		if instance.Name == name {
			kc.Configs[i].Endpoint = endpoint
			kc.Configs[i].Password = base64.StdEncoding.EncodeToString([]byte(password))
			kc.Configs[i].InsecureSkipVerify = !tlsVerify
			return nil
		}
	}
	return errors.New("update instance")
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
	client := http.DefaultClient
	if instance.InsecureSkipVerify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}

	resp, err := client.Do(req)

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

func (instance Instance) GetApps() (*ListAppsResponse, error) {
	token, err := instance.getLoginToken()
	if err != nil {
		errors.Wrap(err, "get login token")
	}

	if token == nil {
		return nil, fmt.Errorf("could not connect")
	}
	appsReq, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/apps", instance.Endpoint), nil)
	if err != nil {
		return nil, errors.Wrap(err, "build apps request")
	}
	appsReq.Header.Set("Accept", "application/json")
	appsReq.Header.Set("Authorization", *token)
	client := http.DefaultClient
	if instance.InsecureSkipVerify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}
	appsResp, err := client.Do(appsReq)
	if err != nil {
		return nil, errors.Wrap(err, "send apps request")
	}

	defer appsResp.Body.Close()
	if appsResp.StatusCode != 200 {
		body, _ := io.ReadAll(appsResp.Body)
		return nil, fmt.Errorf("apps %d: %s", appsResp.StatusCode, body)
	}
	bodyBytes, err := io.ReadAll(appsResp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read body")
	}
	var body ListAppsResponse
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&body); err != nil {
		return nil, errors.Wrap(err, "decode body")
	}
	return &body, nil

}

type ListAppsResponse struct {
	Apps []ResponseApp `json:"apps"`
}

type ResponseApp struct {
	Name       string             `json:"name"`
	Slug       string             `json:"slug"`
	Downstream ResponseDownstream `json:"downstream"`
}

type ResponseDownstream struct {
	CurrentVersion  *DownstreamVersion   `json:"currentVersion"`
	PendingVersions []*DownstreamVersion `json:"pendingVersions"`
}

type DownstreamVersion struct {
	VersionLabel string `json:"versionLabel"`
	Sequence     int    `json:"sequence"`
}

func (instance Instance) UpdateApps() error {
	lappr, err := instance.GetApps()
	if err != nil {
		return errors.Wrap(err, "get apps")
	}
	for _, rapp := range lappr.Apps {
		if len(rapp.Downstream.PendingVersions) > 0 {
			sequence := 0
			for _, pversion := range rapp.Downstream.PendingVersions {
				if pversion.Sequence > sequence {
					sequence = pversion.Sequence
				}
			}
			err = instance.updateApp(rapp.Slug, sequence)
			if err != nil {
				return errors.Wrap(err, "update app")
			}
		}
	}
	return nil
}

func (instance Instance) updateApp(slug string, seq int) error {
	token, err := instance.getLoginToken()
	if err != nil {
		errors.Wrap(err, "get login token")
	}

	if token == nil {
		return fmt.Errorf("could not connect")
	}

	deployParams := map[string]bool{
		"isSkipPreflights":             false,
		"continueWithFailedPreflights": true,
		"isCLI":                        false,
	}

	deployBody, err := json.Marshal(deployParams)
	if err != nil {
		return errors.Wrap(err, "marshal deploy params")
	}

	updateReq, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/app/%s/sequence/%d/deploy", instance.Endpoint, slug, seq), bytes.NewBuffer(deployBody))
	if err != nil {
		return errors.Wrap(err, "build update request")
	}
	updateReq.Header.Set("Accept", "application/json")
	updateReq.Header.Set("Authorization", *token)
	client := http.DefaultClient
	if instance.InsecureSkipVerify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}
	appsResp, err := client.Do(updateReq)
	if err != nil {
		return errors.Wrap(err, "send update request")
	}

	defer appsResp.Body.Close()
	if appsResp.StatusCode != 200 {
		body, _ := io.ReadAll(appsResp.Body)
		return fmt.Errorf("update %d: %s", appsResp.StatusCode, body)
	}
	_, err = io.ReadAll(appsResp.Body)
	if err != nil {
		return errors.Wrap(err, "read body")
	}
	return nil
}

func (instance Instance) getLoginToken() (*string, error) {
	password, err := base64.StdEncoding.DecodeString(instance.Password)
	if err != nil {
		return nil, errors.Wrap(err, "decode password")
	}

	loginParams := map[string]string{
		"password": string(password),
	}

	loginBody, err := json.Marshal(loginParams)
	if err != nil {
		return nil, errors.Wrap(err, "marshal login params")
	}

	loginReq, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/login", instance.Endpoint), bytes.NewBuffer(loginBody))
	if err != nil {
		return nil, errors.Wrap(err, "build login request")
	}
	loginReq.Header.Set("Accept", "application/json")
	loginReq.Header.Set("Content-Type", "application/json")
	client := http.DefaultClient
	if instance.InsecureSkipVerify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}
	loginResp, err := client.Do(loginReq)
	if err != nil {
		return nil, errors.Wrap(err, "send login request")
	}

	defer loginResp.Body.Close()
	if loginResp.StatusCode != 200 {
		body, _ := io.ReadAll(loginResp.Body)
		return nil, fmt.Errorf("login %d: %s", loginResp.StatusCode, body)
	}
	bodyBytes, err := io.ReadAll(loginResp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read body")
	}
	var body SessionResponse
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&body); err != nil {
		return nil, errors.Wrap(err, "decode body")
	}

	return &body.Token, nil
}

type SessionResponse struct {
	Token string `json:"token"`
}
