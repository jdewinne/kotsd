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
						Name:     "kots1",
						Endpoint: "http://something",
						Password: base64.StdEncoding.EncodeToString([]byte("1234abcd")),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ParseConfig(input)
			actual.AddInstance("kots1", "http://something", "1234abcd")
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
						Name:     "kots2",
						Endpoint: "http://another",
						Password: base64.StdEncoding.EncodeToString([]byte("1234abcd")),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ParseConfig(input)
			actual.AddInstance("kots2", "http://another", "1234abcd")
			require.NoError(t, err)
			require.NotNil(t, actual)
			require.Len(t, actual.Configs, 2)
			assert.Equal(t, tt.wantOut.Configs[1], actual.Configs[1])
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
	require.Error(t, err, "Login 403: {\"error\":\"Invalid password. Please try again.\"}")
}
