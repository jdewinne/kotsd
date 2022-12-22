package kotsd

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListInstances(t *testing.T) {
	input := []byte(`instances:
- name: kots1
  endpoint: http://something
- name: kots2
  endpoint: http://another
- name: local
  endpoint: http://localhost:8800`)

	tests := []struct {
		name    string
		wantOut KotsdConfig
	}{
		{
			name: "test",

			wantOut: KotsdConfig{
				Configs: []Instance{
					{
						Name:     "kots1",
						Endpoint: "http://something",
					},
					{
						Name:     "kots2",
						Endpoint: "http://another",
					},
					{
						Name:     "local",
						Endpoint: "http://localhost:8800",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ParseConfig(input)
			require.NoError(t, err)
			require.NotNil(t, actual)
			require.Len(t, actual.Configs, 3)
			assert.Equal(t, tt.wantOut.Configs[1], actual.Configs[1])
		})
	}

}

func TestAddInstanceEmpty(t *testing.T) {
	input := []byte(``)

	tests := []struct {
		name    string
		wantOut KotsdConfig
	}{
		{
			name: "test",

			wantOut: KotsdConfig{
				Configs: []Instance{
					{
						Name:               "kots1",
						Endpoint:           "http://something",
						Password:           base64.StdEncoding.EncodeToString([]byte("1234abcd")),
						InsecureSkipVerify: false,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ParseConfig(input)
			actual.AddInstance("kots1", "http://something", "1234abcd", true)
			require.NoError(t, err)
			require.NotNil(t, actual)
			require.Len(t, actual.Configs, 1)
			assert.Equal(t, tt.wantOut.Configs[0], actual.Configs[0])
		})
	}
}

func TestAddInstanceExisting(t *testing.T) {
	input := []byte(`instances:
- name: kots1
  endpoint: http://something`)

	tests := []struct {
		name    string
		wantOut KotsdConfig
	}{
		{
			name: "test",

			wantOut: KotsdConfig{
				Configs: []Instance{
					{
						Name:     "kots1",
						Endpoint: "http://something",
					},
					{
						Name:               "kots2",
						Endpoint:           "http://another",
						Password:           base64.StdEncoding.EncodeToString([]byte("1234abcd")),
						InsecureSkipVerify: false,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ParseConfig(input)
			actual.AddInstance("kots2", "http://another", "1234abcd", true)
			require.NoError(t, err)
			require.NotNil(t, actual)
			require.Len(t, actual.Configs, 2)
			assert.Equal(t, tt.wantOut.Configs[1], actual.Configs[1])
		})
	}
}

func TestDeleteInstance(t *testing.T) {
	input := []byte(`instances:
- name: kots1
  endpoint: http://something
- name: kots2
  endpoint: http://another`)

	tests := []struct {
		name    string
		wantOut KotsdConfig
	}{
		{
			name: "test",

			wantOut: KotsdConfig{
				Configs: []Instance{
					{
						Name:     "kots1",
						Endpoint: "http://something",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ParseConfig(input)
			actual.DeleteInstance("kots2")
			require.NoError(t, err)
			require.NotNil(t, actual)
			require.Len(t, actual.Configs, 1)
			assert.Equal(t, tt.wantOut.Configs[0], actual.Configs[0])
		})
	}
}

func TestGetLoginToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/login" {
			t.Errorf("Expected to request '/api/v1/login', got: %s", r.URL.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"token":"abcdefgh"}`))
	}))
	defer server.Close()
	i := Instance{Name: "t1", Endpoint: server.URL, Password: base64.StdEncoding.EncodeToString([]byte("1234abcd"))}

	value, err := i.getLoginToken()
	require.NoError(t, err)
	assert.Equal(t, "abcdefgh", *value)
}

func TestGetInvalidPassword(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/login" {
			t.Errorf("Expected to request '/api/v1/login', got: %s", r.URL.Path)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"Invalid password. Please try again."}`))
	}))
	defer server.Close()
	i := Instance{Name: "t1", Endpoint: server.URL, Password: base64.StdEncoding.EncodeToString([]byte("1234abcd"))}

	_, err := i.getLoginToken()
	require.Error(t, err, `Login 403: {"error":"Invalid password. Please try again."}`)
}

func TestGetApps(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}
		if r.URL.Path == "/api/v1/apps" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"apps":[{"name":"DemoApp","downstream":{"currentVersion":{"versionLabel":"1.0.2"},"pendingVersions":[]}}]}`))
		} else if r.URL.Path == "/api/v1/login" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"token":"abcdefgh"}`))
		} else {
			t.Errorf("Expected to request '/api/v1/apps' or '/api/v1/login', got: %s", r.URL.Path)
		}

	}))
	defer server.Close()
	i := Instance{Name: "t1", Endpoint: server.URL, Password: base64.StdEncoding.EncodeToString([]byte("1234abcd"))}

	value, err := i.GetApps()
	require.NoError(t, err)
	rd := ListAppsResponse{Apps: []ResponseApp{{Name: "DemoApp", Downstream: ResponseDownstream{CurrentVersion: &DownstreamVersion{VersionLabel: "1.0.2"}}}}}
	assert.Equal(t, rd.Apps[0].Name, value.Apps[0].Name)
}

func TestUpdateApps(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected Accept: application/json header, got: %s", r.Header.Get("Accept"))
		}
		if r.URL.Path == "/api/v1/apps" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"apps":[{"name":"DemoApp","downstream":{"currentVersion":{"versionLabel":"1.0.2"},"pendingVersions":[]}}]}`))
		} else if r.URL.Path == "/api/v1/login" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"token":"abcdefgh"}`))
		} else {
			t.Errorf("Expected to request '/api/v1/apps' or '/api/v1/login', got: %s", r.URL.Path)
		}

	}))
	defer server.Close()
	i := Instance{Name: "t1", Endpoint: server.URL, Password: base64.StdEncoding.EncodeToString([]byte("1234abcd"))}

	err := i.UpdateApps()
	require.NoError(t, err)
}
